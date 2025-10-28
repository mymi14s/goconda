package models

import (
	"time"
)

type NHSJob struct {
	Id          string `orm:"pk" json:"id"`
	Reference   string `orm:"size(100)" json:"reference"`
	Title       string `orm:"size(255)" json:"title"`
	Description string `orm:"type(text)" json:"description"`
	// URL         string    `orm:"size(500);column(url)" json:"url"`
	Employer  string    `orm:"size(255)" json:"employer"`
	Type      string    `orm:"size(100)" json:"type"`
	Location  string    `orm:"size(255)" json:"location"`
	Salary    string    `orm:"size(100)" json:"salary"`
	PostDate  time.Time `orm:"type(datetime)" json:"postDate"`
	CloseDate time.Time `orm:"type(datetime)" json:"closeDate"`
}

func (u *NHSJob) TableName() string { return "nhs_job" }
