package sql

import (
	"github.com/aliworkshop/dbcore"
	"gorm.io/gorm"
)

func (db *db) handleQueryActions(q *gorm.DB, query dbcore.QueryModel) *gorm.DB {
	extraActions := query.GetExtraActions()
	if extraActions != nil {
		for k, v := range extraActions {
			switch k {
			case "func":
				f := v.(func(*gorm.DB))
				f(q)
				break
			}
		}
	}
	return q
}
