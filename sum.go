package sql

import (
	"fmt"
	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/error"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func (db *db) sum(g *gorm.DB, query dbcore.QueryModel, key string) (decimal.Decimal, error.ErrorModel) {
	q, _ := db.Filter(g, query)
	q = db.handleQueryActions(q, query)
	var sum = new(decimal.Decimal)
	q = db.Join(q, query)
	for _, field := range query.GetGroupBy() {
		q = q.Group(field)
	}
	r := q.Select(fmt.Sprintf("ifnull(sum(%v), 0)", key))
	if r.Error != nil {
		return decimal.Zero, errorHandler(error.Internal(r.Error))
	}
	if err := r.Row().Scan(sum); err != nil {
		return decimal.Zero, errorHandler(error.Internal(err))
	}
	return *sum, nil
}

func (db *db) Sum(query dbcore.QueryModel, key string) (decimal.Decimal, error.ErrorModel) {
	model := query.GetModel()
	dbQuery := db.GetGormDB(query).Model(model)
	return db.sum(dbQuery, query, key)
}

func (db *db) SumWithDFilters(query dbcore.QueryModel, key string) (decimal.Decimal, error.ErrorModel) {
	gq := db.GetGormDB(query)
	gq, _ = db.dFilter(gq, query)
	return db.sum(gq, query, key)
}
