package sql

import (
	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/error"
)

func (db *repo) Delete(query dbcore.QueryModel) (err error.ErrorModel) {
	model := query.GetModel()
	q := db.GetGormDB(query).Model(model)
	q, filtered := db.Filter(q, query)
	q, dFiltered := db.dFilter(q, query)
	if !filtered && !dFiltered {
		err = errorHandler(error.New().
			WithType(error.TypeValidation).
			WithId("NoFilterSelected").
			WithMessage("you can not delete all items"))
		return
	}
	dbc := q.Delete(model)
	if dbc.Error != nil {
		err = errorHandler(error.Internal(dbc.Error))
		return
	}
	if dbc.RowsAffected == 0 {
		err = errorHandler(error.NotFound(nil).WithDetail("NoRowsAffected"))
		return
	}
	return
}
