package models

import (
	"time"
)

type ErrorLog struct {
	BaseModel
	ID        int64     `orm:"auto;column(id)" json:"id"`
	Title     string    `orm:"size(100)" json:"title"`
	Context   string    `orm:"size(150)" json:"context"`
	Error     string    `orm:"type(longtext)" json:"error"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)" json:"created_at"`
}

func (e *ErrorLog) TableName() string { return "error_log" }
