package sql

import (
	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/error"
)

func (db *db) Insert(query dbcore.QueryModel) (result interface{}, err error.ErrorModel) {
	result = query.GetBody()
	if result == nil {
		err = errorHandler(error.New().
			WithType(error.TypeValidation).
			WithId("InsertGetBodyNil").
			WithDetail("nil body in insertion"))
		return
	}
	dbc := db.GetGormDB(query).Create(result)
	if dbc.Error != nil {
		err = errorHandler(error.Internal(dbc.Error))
		return
	}
	return
}
