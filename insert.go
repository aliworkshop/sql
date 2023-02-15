package sql

import (
	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/errorslib"
)

func (db *db) Insert(query dbcore.QueryModel) (result interface{}, err errorslib.ErrorModel) {
	result = query.GetBody()
	if result == nil {
		err = errorHandler(errorslib.New().
			WithType(errorslib.TypeValidation).
			WithId("InsertGetBodyNil").
			WithDetail("nil body in insertion"))
		return
	}
	dbc := db.GetGormDB(query).Create(result)
	if dbc.Error != nil {
		err = errorHandler(errorslib.Internal(dbc.Error))
		return
	}
	return
}
