package controllers

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/beego/beego/v2/client/orm"
	base_controller "github.com/mymi14s/goconda/controllers"
	"github.com/mymi14s/goconda/models"
	"github.com/mymi14s/goconda/utils"
)

const LAYOUT = "2006-01-02T15:04:05.999999999"
const ERROR_TITLE = "NHS Crawler"

type NHSController struct {
	base_controller.BaseController
}

type NHSJobSearch struct {
	Title, Description, Employer, Type, Location, PostDate, CloseDate string
	PageNumber, PageSize                                              int
}

type NHSJobData struct {
	Id          string `xml:"id"`
	Reference   string `xml:"reference"`
	Title       string `xml:"title"`
	Description string `xml:"description"`
	URL         string `xml:"url"`
	Employer    string `xml:"employer"`
	Type        string `xml:"type"`
	Location    string `xml:"locations>location"`
	Salary      string `xml:"salary"`
	PostDate    string `xml:"postDate"`
	CloseDate   string `xml:"closeDate"`
}

type NHSJobsXML struct {
	XMLName        xml.Name     `xml:"nhsJobs"`
	VacancyDetails []NHSJobData `xml:"vacancyDetails"`
	TotalPages     int          `xml:"totalPages"`
	TotalResults   int          `xml:"totalResults"`
}

func (c *NHSController) Index() {
	c.TplName = "frontend/nhs.html"

	if err := c.Render(); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.WriteString("template render error")
	}

}

func (c *NHSController) GetData() {
	var form NHSJobSearch
	if err := c.ParseJSON(&form); err != nil {
		fmt.Println((err))
		c.JSONError(400, err.Error())
		return
	}

	form.PageSize = 10

	jobs, total, err := GetNHSJobs(form)
	if err != nil {
		c.JSONError(400, err.Error())
	}

	c.JSONOK(map[string]any{"data": jobs, "total": total})
}

func GetNHSJobs(search NHSJobSearch) ([]models.NHSJob, int64, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.NHSJob))

	// Filters
	if search.Title != "" {
		qs = qs.Filter("Title__icontains", search.Title)
	}
	if search.Description != "" {
		qs = qs.Filter("Description__icontains", search.Description)
	}
	if search.Employer != "" {
		qs = qs.Filter("Employer__icontains", search.Employer)
	}
	if search.Type != "" {
		qs = qs.Filter("Type__icontains", search.Type)
	}
	if search.Location != "" {
		qs = qs.Filter("Location__icontains", search.Location)
	}

	today := time.Now().Format("2006-01-02")

	if search.PostDate != "" {
		postDate, err := time.Parse("2006-01-02", search.PostDate)
		if err == nil {
			todayDate, _ := time.Parse("2006-01-02", today)

			if postDate.Before(todayDate) {
				qs = qs.Filter("PostDate__gte", postDate).
					Filter("PostDate__lte", todayDate)
			} else {
				qs = qs.Filter("PostDate__gte", postDate)
			}
		}
	}
	if search.CloseDate != "" {
		qs = qs.Filter("CloseDate__lte", search.CloseDate)
	}

	qs = qs.OrderBy("-PostDate")

	var jobs []models.NHSJob
	var total int64
	var err error

	if search.Title == "" && search.Description == "" && search.Employer == "" &&
		search.Type == "" && search.Location == "" && search.PostDate == "" && search.CloseDate == "" {
		// No filters: get last 10
		_, err = qs.Limit(10).All(&jobs)
		total = int64(len(jobs))
	} else {
		// Pagination
		if search.PageNumber < 1 {
			search.PageNumber = 1
		}
		if search.PageSize < 1 {
			search.PageSize = 10
		}
		offset := (search.PageNumber - 1) * search.PageSize
		total, err = qs.Count()
		if err != nil {
			return nil, 0, err
		}
		_, err = qs.Limit(search.PageSize, offset).All(&jobs)
	}

	return jobs, total, err
}

func NHSCrawler() error {
	loc, err := time.LoadLocation("Europe/London")
	if err != nil {
		panic(err)
	}
	now := time.Now().In(loc)
	oneHourAgo := now.Add(-1 * time.Hour)
	pageNumber := 1
	for pageNumber > 0 {
		requestURL := fmt.Sprintf("https://www.jobs.nhs.uk/api/v1/search_xml?sort=publicationDateDesc&page=%d", pageNumber)
		req, err := http.NewRequest(http.MethodGet, requestURL, nil)
		if err != nil {
			utils.LogError(map[string]any{
				"title":  ERROR_TITLE,
				"source": "NHS",
				"error":  err,
			})
			os.Exit(1)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			utils.LogError(map[string]any{
				"title":  ERROR_TITLE,
				"source": "NHS",
				"error":  err,
			})
			os.Exit(1)
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			utils.LogError(map[string]any{
				"title":  ERROR_TITLE,
				"source": "NHS",
				"error":  err,
			})
			os.Exit(1)
		}
		var jobs NHSJobsXML

		if err := xml.Unmarshal(resBody, &jobs); err != nil {
			utils.LogError(map[string]any{
				"title":  ERROR_TITLE,
				"source": "NHS",
				"error":  err,
			})
		}
		for _, v := range jobs.VacancyDetails {
			postDate, err := time.ParseInLocation(LAYOUT, v.PostDate, loc)

			if err != nil {
				utils.LogError(map[string]any{
					"title":  ERROR_TITLE,
					"source": "NHS",
					"error":  err,
				})
				return nil
			}
			if postDate.After(oneHourAgo) && postDate.Before(now) {

				data, _ := ConvertNHSJobDataToModel(v)
				o := orm.NewOrm()
				if _, err := o.Insert(&data); err != nil {
					fmt.Println(err)
				}

			} else {
				pageNumber = 0
				break
			}
		}
		if pageNumber > 0 {
			pageNumber++
		}
	}

	return nil
}

func DeleteExpiredNHSJobs() {

	o := orm.NewOrm()
	today := time.Now()

	_, err := o.QueryTable(new(models.NHSJob)).
		Filter("CloseDate__lt", today).
		Delete()
	if err != nil {
		utils.LogError(map[string]any{
			"title":  ERROR_TITLE,
			"source": "NHS",
			"error":  err,
		})
	}
}

// Convert NHSJobData → NHSJob
func ConvertNHSJobDataToModel(data NHSJobData) (models.NHSJob, error) {
	var job models.NHSJob
	var err error

	// Parse PostDate
	if data.PostDate != "" {
		job.PostDate, err = time.ParseInLocation(LAYOUT, data.PostDate, time.Local)
		if err != nil {
			return job, err
		}
	}

	// Parse CloseDate
	if data.CloseDate != "" {
		job.CloseDate, err = time.ParseInLocation("2006-01-02", data.CloseDate, time.Local)
		if err != nil {
			return job, err
		}
	}

	// Map other fields
	job.Reference = data.Reference
	job.Title = data.Title
	job.Description = data.Description
	// job.URL = data.URL
	job.Employer = data.Employer
	job.Type = data.Type
	job.Location = data.Location
	job.Salary = data.Salary
	job.Id = data.Id

	return job, nil
}
