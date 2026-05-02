package sql

import (
	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/errors"
)

func (db *repo) Delete(query dbcore.QueryModel) (err errors.ErrorModel) {
	model := query.GetModel()
	q := db.GetGormDB(query).Model(model)
	q, filtered := db.Filter(q, query)
	q, dFiltered := db.dFilter(q, query)
	if !filtered && !dFiltered {
		err = errorHandler(errors.New().
			WithType(errors.TypeValidation).
			WithId("NoFilterSelected").
			WithMessage("you can not delete all items"))
		return
	}
	dbc := q.Delete(model)
	if dbc.Error != nil {
		err = errorHandler(errors.Internal(dbc.Error))
		return
	}
	if dbc.RowsAffected == 0 {
		err = errorHandler(errors.NotFound(nil).WithDetail("NoRowsAffected"))
		return
	}
	return
}
