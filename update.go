package sql

import (
	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/errors"
)

func (db *repo) Update(query dbcore.QueryModel) (err errors.ErrorModel) {
	entity := query.GetBody()
	if entity == nil {
		err = errorHandler(errors.Validation())
		return
	}
	q := db.GetGormDB(query)
	q, filtered := db.Filter(q, query)
	if !filtered {
		err = errorHandler(errors.New().
			WithType(errors.TypeValidation).
			WithDetail("query must be set..no query is set as filter"))
		return
	}
	dbc := q.Updates(entity)
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
