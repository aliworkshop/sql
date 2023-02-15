package sql

import (
	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/errorslib"
	"gorm.io/gorm"
)

func (db *db) sort(con *gorm.DB, query dbcore.QueryModel) (q *gorm.DB, err errorslib.ErrorModel) {
	q = con
	sort := query.GetSort()
	if sort != nil {
		for _, s := range sort {
			sort := s.Field
			if s.Order != "" {
				if s.Order.IsDescending() {
					sort += " DESC"
				} else if !s.Order.IsAscending() {
					err = errorHandler(errorslib.New().
						WithType(errorslib.TypeValidation).
						WithId("InvalidSortQuery").
						WithDetail("invalid sort order. order takes one of values of `DESC` or `ASC`"))
					return
				}
			}
			q = q.Order(sort)
		}
	}
	return
}
