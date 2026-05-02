package sql

import (
	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/errors"
	"github.com/jackc/pgx/v5/pgconn"
)

func (db *repo) Insert(query dbcore.QueryModel) (result interface{}, err errors.ErrorModel) {
	result = query.GetBody()
	if result == nil {
		err = errorHandler(errors.New().
			WithType(errors.TypeValidation).
			WithId("InsertGetBodyNil").
			WithDetail("nil body in insertion"))
		return
	}
	dbc := db.GetGormDB(query).Create(result)
	if dbc.Error != nil {
		var e *pgconn.PgError
		if errors.As(dbc.Error, &e) {
			if e.Code == "23505" {
				return nil, errors.Duplicate(dbc.Error)
			}
		}
		return nil, errorHandler(errors.Internal(dbc.Error))
	}
	return
}
