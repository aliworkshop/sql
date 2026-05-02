package sql

import (
	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/errors"
	"gorm.io/gorm"
	"reflect"
)

func (db *repo) count(gq *gorm.DB, query dbcore.QueryModel) (uint64, errors.ErrorModel) {
	gq, _ = db.Filter(gq, query)
	gq = db.Join(gq, query)
	if table, args := query.GetTable(); table != "" {
		gq.Table(table, db.handleArgs(args)...)
	}
	for _, field := range query.GetGroupBy() {
		gq = gq.Group(field)
	}

	if query.IsUnscoped() {
		gq.Unscoped()
	}

	var c int64
	r := gq.Count(&c)
	if r.Error != nil {
		return 0, errorHandler(errors.Internal(r.Error))
	}
	return uint64(c), nil
}

func (db *repo) Count(query dbcore.QueryModel) (count uint64, err errors.ErrorModel) {
	dbQuery := db.GetGormDB(query)
	return db.count(dbQuery, query)
}

func (db *repo) CountWithDFilter(query dbcore.QueryModel) (count uint64, err errors.ErrorModel) {
	gq := db.GetGormDB(query)
	gq, _ = db.dFilter(gq, query)
	return db.count(gq, query)
}

func (db *repo) list(gq *gorm.DB, query dbcore.QueryModel) (pr interface{}, err errors.ErrorModel) {
	q, _ := db.Filter(gq, query)
	q, err = db.sort(q, query)
	if err != nil {
		return nil, err
	}
	offset := (query.GetPage() - 1) * query.GetPageSize()
	q = db.Join(q, query)
	for _, sel := range query.GetSelects() {
		q = q.Select(sel.Columns, db.handleArgs(sel.Args)...)
	}

	typ := reflect.TypeOf(query.GetModel())
	result := reflect.New(reflect.SliceOf(typ)).Elem().Interface()
	if table, args := query.GetTable(); table != "" {
		q.Table(table, db.handleArgs(args)...)
	}
	for _, field := range query.GetGroupBy() {
		q = q.Group(field)
	}
	for _, relation := range query.GetPreloads() {
		q = q.Preload(relation.Preload, relation.Args...)
	}
	if query.GetPageSize() != -1 {
		q = q.Offset(offset).Limit(query.GetPageSize())
	}
	if query.IsUnscoped() {
		q.Unscoped()
	}
	dbc := q.Find(&result)
	if dbc.Error != nil {
		return nil, errorHandler(errors.Internal(dbc.Error))
	}
	return result, nil
}

func (db *repo) List(query dbcore.QueryModel) (pr interface{}, err errors.ErrorModel) {
	gq := db.GetGormDB(query)
	return db.list(gq, query)
}

func (db *repo) ListWithDFilter(query dbcore.QueryModel) (items interface{}, err errors.ErrorModel) {
	gq := db.GetGormDB(query)
	gq, _ = db.dFilter(gq, query)
	return db.list(gq, query)
}

func (db *repo) Get(query dbcore.QueryModel) (item interface{}, err errors.ErrorModel) {
	q := db.GetGormDB(query)
	q, filtered := db.Filter(q, query)
	if !filtered {
		err = errorHandler(errors.New().
			WithType(errors.TypeValidation).
			WithDetail("query must be set..no query is set as filter"))
		return
	}
	q, err = db.sort(q, query)
	if err != nil {
		return nil, err
	}
	q = db.Join(q, query)
	q = q.Limit(1)
	for _, sel := range query.GetSelects() {
		q = q.Select(sel.Columns, db.handleArgs(sel.Args)...)
	}
	result := query.GetModel()
	if table, args := query.GetTable(); table != "" {
		q.Table(table, db.handleArgs(args)...)
	}
	for _, relation := range query.GetPreloads() {
		q = q.Preload(relation.Preload, relation.Args...)
	}
	dbc := q.Find(result)
	if dbc.Error != nil {
		err = errors.Internal(dbc.Error)
		return
	}
	if dbc.RowsAffected == 0 {
		err = dbcore.NotFoundErr
		return
	}
	item = result
	return
}

func (db *repo) Exist(query dbcore.QueryModel) (exists bool, err errors.ErrorModel) {
	count, e := db.Count(query)
	if e != nil {
		return false, e
	}

	return count > 0, nil
}

func (db *repo) handleArgs(args []interface{}) []interface{} {
	for i, arg := range args {
		if dbcore.IsQueryModel(arg) {
			query := dbcore.GetQueryModel(arg)
			q := db.GetGormDB(query)
			q, _ = db.Filter(q, query)
			q = db.Join(q, query)
			for _, sel := range query.GetSelects() {
				q = q.Select(sel.Columns, db.handleArgs(sel.Args)...)
			}
			if table, tableArgs := query.GetTable(); table != "" {
				q.Table(table, db.handleArgs(tableArgs)...)
			}
			if qs, qArgs := query.GetQuery(); qs != "" {
				q = q.Raw(qs, qArgs...)
			}
			for _, field := range query.GetGroupBy() {
				q = q.Group(field)
			}
			args[i] = q
		}
	}
	return args
}
