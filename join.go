package sql

import (
	"github.com/aliworkshop/dbcore"
	"gorm.io/gorm"
)

func (db *db) Join(gq *gorm.DB, query dbcore.QueryModel) *gorm.DB {
	joins := query.GetJoin()
	for _, join := range joins {
		gq = gq.Joins(join.Query, db.handleArgs(join.Args))
	}
	return gq
}
