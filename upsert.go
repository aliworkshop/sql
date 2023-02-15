package sql

import (
	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/errorslib"
	"gorm.io/gorm/clause"
)

func (db *db) Upsert(query dbcore.QueryModel) (err errorslib.ErrorModel) {
	entity := query.GetBody()
	if entity == nil {
		err = errorHandler(errorslib.DefaultValidationError)
		return
	}
	q := db.GetGormDB(query)
	q = db.handleQueryActions(q, query)
	q, _ = db.Filter(q, query)
	dbc := q.
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(entity)
	if dbc.Error != nil {
		err = errorHandler(errorslib.Internal(dbc.Error))
		return
	}
	return
}
