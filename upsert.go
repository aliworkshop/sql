package sql

import (
	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/error"
	"gorm.io/gorm/clause"
)

func (db *repo) Upsert(query dbcore.QueryModel) (err error.ErrorModel) {
	entity := query.GetBody()
	if entity == nil {
		err = errorHandler(error.DefaultValidationError)
		return
	}
	q := db.GetGormDB(query)
	q, _ = db.Filter(q, query)
	dbc := q.
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(entity)
	if dbc.Error != nil {
		err = errorHandler(error.Internal(dbc.Error))
		return
	}
	return
}
