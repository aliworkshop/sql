package sql

import (
	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/errors"
	"gorm.io/gorm/clause"
)

func (db *repo) Upsert(query dbcore.QueryModel) (err errors.ErrorModel) {
	entity := query.GetBody()
	if entity == nil {
		err = errorHandler(errors.Validation())
		return
	}
	q := db.GetGormDB(query)
	q, _ = db.Filter(q, query)
	dbc := q.
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(entity)
	if dbc.Error != nil {
		err = errorHandler(errors.Internal(dbc.Error))
		return
	}
	return
}
