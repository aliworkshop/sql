package sql

import (
	"github.com/aliworkshop/errorslib"
	"github.com/aliworkshop/dbcore"
	"gorm.io/gorm"
	"gorm.io/hints"
)

func (db *db) GetByFunction(query dbcore.QueryModel, function string) (resultByFunction interface{}, err errorslib.ErrorModel) {
	q := db.GetGormDB(query)
	//q, _ = db.Filter(q, query)
	switch function {
	case "count":
		var count int64
		q = q.Count(&count)
		resultByFunction = count
	}
	if q.Error != nil {
		err = errorHandler(errorslib.Internal(q.Error))
		return
	}
	return
}

func (db *db) GetDbItemsCount(gq *gorm.DB, query dbcore.QueryModel) (uint64, errorslib.ErrorModel) {
	gq, _ = db.Filter(gq, query)
	gq = db.handleQueryActions(gq, query)
	gq = db.Join(gq, query)
	var c int64
	r := gq.Count(&c)
	if r.Error != nil {
		return 0, errorHandler(errorslib.Internal(r.Error))
	}
	return uint64(c), nil
}

func (db *db) GetItemsCount(query dbcore.QueryModel) (count uint64, err errorslib.ErrorModel) {
	dbQuery := db.GetGormDB(query)
	return db.GetDbItemsCount(dbQuery, query)
}

func (db *db) GetItemsCountWithDFilters(query dbcore.QueryModel) (count uint64, err errorslib.ErrorModel) {
	gq := db.GetGormDB(query)
	gq, _ = db.dFilter(gq, query)
	return db.GetDbItemsCount(gq, query)
}

// GetDbItems handle get items using given gq.
// parameter gq is gorm query which is already initialized by caller
func (db *db) GetDbItems(gq *gorm.DB, query dbcore.QueryModel) (pr interface{}, err errorslib.ErrorModel) {
	q, _ := db.Filter(gq, query)
	q = db.handleQueryActions(q, query)
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
	result := query.GetModels()
	if table, args := query.GetTable(); table != "" {
		q.Table(table, db.handleArgs(args))
	}
	for _, field := range query.GetGroupBy() {
		q = q.Group(field)
	}
	dbc := q.Offset(offset).Limit(query.GetPageSize()).Find(result)
	if dbc.Error != nil {
		return nil, errorHandler(errorslib.Internal(dbc.Error))
	}
	return result, nil
}

func (db *db) GetItems(query dbcore.QueryModel) (pr interface{}, err errorslib.ErrorModel) {
	gq := db.GetGormDB(query)
	return db.GetDbItems(gq, query)
}

func (db *db) GetItemsWithDFilters(query dbcore.QueryModel) (items interface{}, err errorslib.ErrorModel) {
	gq := db.GetGormDB(query)
	gq, _ = db.dFilter(gq, query)
	return db.GetDbItems(gq, query)
}

func (db *db) GetItem(query dbcore.QueryModel) (item interface{}, err errorslib.ErrorModel) {
	q := db.GetGormDB(query)
	q, filtered := db.Filter(q, query)
	if !filtered {
		err = errorHandler(errorslib.New().
			WithType(errorslib.TypeValidation).
			WithDetail("query must be set..no query is set as filter"))
		return
	}
	q, err = db.sort(q, query)
	if err != nil {
		return nil, err
	}
	q = db.Join(q, query)
	q = q.Limit(1)
	q = db.handleQueryActions(q, query)
	for _, sel := range query.GetSelects() {
		q = q.Select(sel.Columns, db.handleArgs(sel.Args))
	}
	result := query.GetModel()
	if table, args := query.GetTable(); table != "" {
		q.Table(table, db.handleArgs(args))
	}
	dbc := q.Find(result)
	if dbc.Error != nil {
		err = errorslib.Internal(dbc.Error)
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
