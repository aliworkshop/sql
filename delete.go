package sql

import (
	"github.com/aliworkshop/errorslib"
	"github.com/aliworkshop/dbcore"
)

func (db *db) Delete(query dbcore.QueryModel) (err errorslib.ErrorModel) {
	model := query.GetModel()
	q := db.GetGormDB(query).Model(model)
	q, filtered := db.Filter(q, query)
	q, dFiltered := db.dFilter(q, query)
	if !filtered && !dFiltered {
		err = errorHandler(errorslib.New().
			WithType(errorslib.TypeValidation).
			WithId("NoFilterSelected").
			WithMessage("you can not delete all items"))
		return
	}
	q = db.handleQueryActions(q, query)
	dbc := q.Delete(model)
	if dbc.Error != nil {
		err = errorHandler(errorslib.Internal(dbc.Error))
		return
	}
	if dbc.RowsAffected == 0 {
		err = errorHandler(errorslib.NotFound(nil).WithDetail("NoRowsAffected"))
		return
	}
	return
}
