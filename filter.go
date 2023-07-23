package sql

import (
	"fmt"
	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/dfilter"
	"gorm.io/gorm"
	"strings"
)

func (db *db) Filter(gormQuery *gorm.DB, query dbcore.QueryModel) (q *gorm.DB, filtered bool) {
	q = gormQuery
	filters := query.GetFilters()
	if filters != nil && len(filters) > 0 {
		for k, v := range filters {
			q = q.Where(fmt.Sprintf(`%s=?`, k), v)
		}
		filtered = true
	}
	extraFilters := query.GetExtraFilters()
	if extraFilters != nil {
		for _, filter := range extraFilters {
			q = q.Where(filter.Query, db.handleArgs(filter.Params)...)
			filtered = true
		}
	}
	if db.queryParser != nil {
		pr := db.queryParser.Parse(query)
		if pr != nil {
			prQuery := pr.GetQuery()
			prParams := pr.GetParams()
			q = q.Where(prQuery, prParams...)
			filtered = true
		}
	}
	return
}

func (db *db) dFilter(dbQuery *gorm.DB, query dbcore.QueryModel) (q *gorm.DB, filtered bool) {
	q = dbQuery
	for _, v := range query.GetDynamicFilters() {
		if t := query.GetDynamicFilterTable(); t != "" {
			v.Key = fmt.Sprintf("%s.%s", t, v.Key)
		}
		if items, ok := v.Value.([]string); ok {
			filtered = true
			if len(items) > 1 {
				q = q.Where(fmt.Sprintf("%s IN (?)", v.Key), v.Value)
				continue
			}
			if strings.Contains(items[0], ",") {
				q = q.Where(fmt.Sprintf("%s IN (?)", v.Key), strings.Split(v.Value.([]string)[0], ","))
				continue
			}
		}
		if v.ValueType == dfilter.Null {
			q = q.Where(fmt.Sprintf("%s IS NULL", v.Key))
			continue
		}
		q = q.Where(fmt.Sprintf("%s %s ?", v.Key, v.SQLOperator), v.Value)
	}
	return
}

func (db *db) DFilter(query dbcore.QueryModel) (dbQuery interface{}, filtered bool) {
	gormDb := db.GetGormDB(query)
	return db.dFilter(gormDb, query)
}
