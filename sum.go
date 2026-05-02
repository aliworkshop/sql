package sql

import (
	"fmt"
	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func (db *repo) sum(g *gorm.DB, query dbcore.QueryModel, key string) (decimal.Decimal, errors.ErrorModel) {
	q, _ := db.Filter(g, query)
	var sum = new(decimal.Decimal)
	q = db.Join(q, query)
	for _, field := range query.GetGroupBy() {
		q = q.Group(field)
	}
	r := q.Select(fmt.Sprintf("COALESCE(sum(%v), 0)", key))
	if r.Error != nil {
		return decimal.Zero, errorHandler(errors.Internal(r.Error))
	}
	if err := r.Row().Scan(sum); err != nil {
		return decimal.Zero, errorHandler(errors.Internal(err))
	}
	return *sum, nil
}

func (db *repo) Sum(query dbcore.QueryModel, key string) (decimal.Decimal, errors.ErrorModel) {
	model := query.GetModel()
	dbQuery := db.GetGormDB(query).Model(model)
	return db.sum(dbQuery, query, key)
}

func (db *repo) SumWithDFilters(query dbcore.QueryModel, key string) (decimal.Decimal, errors.ErrorModel) {
	gq := db.GetGormDB(query)
	gq, _ = db.dFilter(gq, query)
	return db.sum(gq, query, key)
}
