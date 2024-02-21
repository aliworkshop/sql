package sql

import (
	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/error"
	"gorm.io/gorm"
)

func (db *repo) sort(con *gorm.DB, query dbcore.QueryModel) (q *gorm.DB, err error.ErrorModel) {
	q = con
	sort := query.GetSort()
	if sort != nil {
		for field, s := range sort {
			if s.ReplaceWith != "" {
				field = s.ReplaceWith
			}
			if s.Order != "" {
				if s.Order.IsDescending() {
					field += " DESC"
				} else if !s.Order.IsAscending() {
					err = errorHandler(error.New().
						WithType(error.TypeValidation).
						WithId("InvalidSortQuery").
						WithDetail("invalid sort order. order takes one of values of `DESC` or `ASC`"))
					return
				}
			}
			q = q.Order(field)
		}
	}
	return
}
