package sql

import (
	"fmt"
	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/dfilter"
	"gorm.io/gorm"
	"strings"
)

func (db *repo) Filter(gormQuery *gorm.DB, query dbcore.QueryModel) (q *gorm.DB, filtered bool) {
	q = gormQuery

	if filters := query.GetFilters(); len(filters) > 0 {
		for _, filter := range filters {
			s, op, vals := "", "", []any{}
			for _, m := range filter.GetMatches() {
				op = string(m.Op)
				s += fmt.Sprintf("%s %s ? %s ", m.Key, m.Operator, m.Op)
				vals = append(vals, m.Value)
			}

			s = s[:len(s)-len(op)-2]
			if filter.GetOperation() == dbcore.And {
				q = q.Where(s, vals...)
			} else if filter.GetOperation() == dbcore.OR {
				q = q.Or(s, vals...)
			}
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

func (db *repo) dFilter(dbQuery *gorm.DB, query dbcore.QueryModel) (q *gorm.DB, filtered bool) {
	q = dbQuery
	for _, v := range query.GetDynamicFilters() {
		if v.KeyType == string(dfilter.Json) {
			field := strings.Split(v.Key, ".")
			switch v.Operator {
			case dfilter.OperatorIn:
				q = q.Where(fmt.Sprintf("JSON_CONTAINS(%s->'$.%s', '[%s]')", field[0], field[1], strings.Join(v.Value.([]string), ",")))
			case dfilter.OperatorNotIn:
				//todo: implement not in
			case dfilter.OperatorEq:
				q = q.Where(fmt.Sprintf("%s->>'$.%s' = ?", field[0], field[1]), v.Value)
			case dfilter.OperatorNot:
				q = q.Where(fmt.Sprintf("%s->>'$.%s' != ?", field[0], field[1]), v.Value)
			default:
				q = q.Where(fmt.Sprintf("%s->>'$.%s'%s?", field[0], field[1], v.SQLOperator), v.Value)
			}
			continue
		}
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
