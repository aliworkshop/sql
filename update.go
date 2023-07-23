package sql

import (
	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/error"
	"reflect"
)

func (db *db) Update(query dbcore.QueryModel) (err error.ErrorModel) {
	entity := query.GetBody()
	if entity == nil {
		err = errorHandler(error.DefaultValidationError)
		return
	}
	q := db.GetGormDB(query)
	q = db.handleQueryActions(q, query)
	q, _ = db.Filter(q, query)
	dbc := q.Updates(entity)
	if dbc.Error != nil {
		err = errorHandler(error.Internal(dbc.Error))
		return
	}
	return
}

func (db *db) UpdateMap(query dbcore.QueryModel) (err error.ErrorModel) {
	entity := query.GetBody()
	if entity == nil {
		err = errorHandler(error.DefaultValidationError)
		return
	}
	q := db.GetGormDB(query)
	q = db.handleQueryActions(q, query)
	q, _ = db.Filter(q, query)
	dbc := q.Updates(structToMap(entity))
	if dbc.Error != nil {
		err = errorHandler(error.Internal(dbc.Error))
		return
	}
	return
}

func structToMap(i interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for j := 0; j < v.NumField(); j++ {
		data[v.Type().Field(j).Name] = v.Field(j).Interface()
	}
	return data
}
