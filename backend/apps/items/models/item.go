package models

import (
	"time"

	"github.com/mymi14s/goconda/models"

	"github.com/beego/beego/v2/client/orm"
)

type Item struct {
	ID          int64        `orm:"auto;pk;column(id)" json:"id"`
	Name        string       `orm:"size(200)" json:"name"`
	Description string       `orm:"type(text)" json:"description"`
	Owner       *models.User `orm:"rel(fk);column(owner_email);index" json:"owner"`
	CreatedAt   time.Time    `orm:"auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt   time.Time    `orm:"auto_now;type(datetime)" json:"updated_at"`
}

func (i *Item) TableName() string { return "item" }

func CreateItem(i *Item) error {
	_, err := orm.NewOrm().Insert(i)
	return err
}

func GetItemByID(id int64) (*Item, error) {
	o := orm.NewOrm()
	it := Item{ID: id}
	if err := o.Read(&it); err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &it, nil
}

func ListItemsByOwner(ownerEmail string, offset, limit int64) ([]*Item, int64, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Item)).Filter("Owner__Email", ownerEmail)
	total, err := qs.Count()
	if err != nil {
		return nil, 0, err
	}
	var items []*Item
	_, err = qs.OrderBy("-ID").Limit(limit, offset).All(&items)
	return items, total, err
}

func init() {
	orm.RegisterModel(new(Item))
}
