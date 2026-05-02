package sql

import (
	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/errors"
)

func (db *repo) Update(query dbcore.QueryModel) (err errors.ErrorModel) {
	entity := query.GetBody()
	if entity == nil {
		err = errorHandler(errors.DefaultValidationError)
		return
	}
	q := db.GetGormDB(query)
	q, _ = db.Filter(q, query)
	dbc := q.Updates(entity)
	if dbc.Error != nil {
		err = errorHandler(errors.Internal(dbc.Error))
		return
	}
	return
}
