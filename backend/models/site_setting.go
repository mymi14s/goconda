package models

import (
	"context"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type SiteSetting struct {
	Id        int64     `orm:"pk"`
	Title     string    `orm:"size(128)"`
	SiteName  string    `orm:"size(128)"`
	BaseURL   string    `orm:"size(255);column(base_url)"`
	Email     string    `orm:"size(128)"`
	Tagline   string    `orm:"size(256);null"`
	Header    string    `orm:"size(5000);null"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime)"`
	Version   int       `orm:"version"`
	Sentinel  string    `orm:"unique;size(16)"`
}

func (s *SiteSetting) TableName() string { return "site_setting" }

func (s *SiteSetting) Get() (*SiteSetting, error) {
	o := orm.NewOrm()
	ss := SiteSetting{Id: 1, Sentinel: "singleton"}
	_, _, err := o.ReadOrCreate(&ss, "Sentinel")
	if err != nil {
		return nil, err
	}
	return &ss, nil
}

func Update(apply func(*SiteSetting) error) error {
	o := orm.NewOrm()
	return o.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		s := SiteSetting{Id: 1}
		if err := txOrm.ReadForUpdate(&s); err != nil {
			if err == orm.ErrNoRows {
				_, _, rcErr := txOrm.ReadOrCreate(&s, "Id")
				if rcErr != nil {
					return rcErr
				}
			} else {
				return err
			}
		}

		if err := apply(&s); err != nil {
			return err
		}

		if _, err := txOrm.Update(&s); err != nil {
			return err
		}
		return nil
	})
}

func init() {
	orm.RegisterModel(new(SiteSetting))
}
