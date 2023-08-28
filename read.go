package sql

import (
	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/error"
	"gorm.io/gorm"
	"gorm.io/hints"
	"reflect"
)

func (db *db) count(gq *gorm.DB, query dbcore.QueryModel) (uint64, error.ErrorModel) {
	gq, _ = db.Filter(gq, query)
	gq = db.Join(gq, query)
	var c int64
	r := gq.Count(&c)
	if r.Error != nil {
		return 0, errorHandler(error.Internal(r.Error))
	}
	return uint64(c), nil
}

func (db *db) Count(query dbcore.QueryModel) (count uint64, err error.ErrorModel) {
	dbQuery := db.GetGormDB(query)
	return db.count(dbQuery, query)
}

func (db *db) CountWithDFilter(query dbcore.QueryModel) (count uint64, err error.ErrorModel) {
	gq := db.GetGormDB(query)
	gq, _ = db.dFilter(gq, query)
	return db.count(gq, query)
}

func (db *db) list(gq *gorm.DB, query dbcore.QueryModel) (pr interface{}, err error.ErrorModel) {
	q, _ := db.Filter(gq, query)
	q, err = db.sort(q, query)
	if err != nil {
		return nil, err
	}
	offset := (query.GetPage() - 1) * query.GetPageSize()
	q = db.Join(q, query)
	q = db.handleHints(q, query)
	for _, sel := range query.GetSelects() {
		q = q.Select(sel.Columns, db.handleArgs(sel.Args))
	}

	typ := reflect.TypeOf(query.GetModel())
	result := reflect.New(reflect.SliceOf(typ)).Elem().Interface()
	if table, args := query.GetTable(); table != "" {
		q.Table(table, db.handleArgs(args))
	}
	for _, field := range query.GetGroupBy() {
		q = q.Group(field)
	}
	dbc := q.Offset(offset).Limit(query.GetPageSize()).Find(&result)
	if dbc.Error != nil {
		return nil, errorHandler(error.Internal(dbc.Error))
	}
	return result, nil
}

func (db *db) List(query dbcore.QueryModel) (pr interface{}, err error.ErrorModel) {
	gq := db.GetGormDB(query)
	return db.list(gq, query)
}

func (db *db) ListWithDFilter(query dbcore.QueryModel) (items interface{}, err error.ErrorModel) {
	gq := db.GetGormDB(query)
	gq, _ = db.dFilter(gq, query)
	return db.list(gq, query)
}

func (db *db) Get(query dbcore.QueryModel) (item interface{}, err error.ErrorModel) {
	q := db.GetGormDB(query)
	q, filtered := db.Filter(q, query)
	if !filtered {
		err = errorHandler(error.New().
			WithType(error.TypeValidation).
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
		q = q.Select(sel.Columns, db.handleArgs(sel.Args))
	}
	result := query.GetModel()
	if table, args := query.GetTable(); table != "" {
		q.Table(table, db.handleArgs(args))
	}
	dbc := q.Find(result)
	if dbc.Error != nil {
		err = error.Internal(dbc.Error)
		return
	}
	if dbc.RowsAffected == 0 {
		err = dbcore.NotFoundErr
		return
	}
	item = result
	return
}

func (db *db) handleHints(gq *gorm.DB, query dbcore.QueryModel) *gorm.DB {
	h := query.GetHint()
	if h != nil {
		switch h.Kind {
		case dbcore.HintForce:
			if h.Name != "" {
				gq.Clauses(hints.ForceIndex(h.Name))
			}
		case dbcore.HintUse:
			if h.Name != "" {
				gq.Clauses(hints.UseIndex(h.Name))
			}
		case dbcore.HintIgnore:
			if h.Name != "" {
				gq.Clauses(hints.IgnoreIndex(h.Name))
			}
		}
	}
	return gq
}

func (db *db) handleArgs(args []interface{}) []interface{} {
	for i, arg := range args {
		if dbcore.IsQueryModel(arg) {
			query := dbcore.GetQueryModel(arg)
			q := db.GetGormDB(query)
			q, _ = db.Filter(q, query)
			for _, sel := range query.GetSelects() {
				q = q.Select(sel.Columns, sel.Args)
			}
			if table, _ := query.GetTable(); table != "" {
				q.Table(query.GetTable())
			}
			for _, field := range query.GetGroupBy() {
				q = q.Group(field)
			}
			args[i] = q
		}
	}
	return args
}
