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
				if strings.Contains(m.Key, "->") {
					field := strings.Split(m.Key, "->")
					s += fmt.Sprintf("%s->>'%s' %s ? %s ", field[0], field[1], m.Operator, m.Op)
					if m.Operator == dbcore.Ct {
						vals = append(vals, "%"+m.Value.(string)+"%")
						continue
					}
				} else {
					if items, ok := m.Value.([]string); ok {
						filtered = true
						if len(items) > 1 {
							s += fmt.Sprintf("\"%s\" IN (?) %s ", m.Key, m.Op)
							vals = append(vals, m.Value.([]string))
							continue
						}
						if strings.Contains(items[0], ",") {
							s += fmt.Sprintf("\"%s\" IN (?) %s", m.Key, m.Op)
							vals = append(vals, strings.Split(m.Value.([]string)[0], ","))
							continue
						}
					}
					if m.Operator == dbcore.Ct {
						s += fmt.Sprintf("%s %s ? %s ", m.Key, m.Operator, m.Op)
						vals = append(vals, "%"+m.Value.(string)+"%")
						continue
					}
					s += fmt.Sprintf("%s %s ? %s ", m.Key, m.Operator, m.Op)
				}
				vals = append(vals, m.Value)
			}

			s = s[:len(s)-len(op)-2]
			if filter.GetOperation() == dbcore.And {
				q = q.Where(s, vals...)
			} else if filter.GetOperation() == dbcore.OR {
				q = q.Or(s, vals...)
			} else if filter.GetOperation() == dbcore.Not {
				q = q.Not(s, vals...)
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
				q = q.Where(fmt.Sprintf("%s->>'%s' IN ?", field[0], field[1]), v.Value)
			case dfilter.OperatorNotIn:
				q = q.Where(fmt.Sprintf("%s->>'%s' NOT IN ?", field[0], field[1]), v.Value)
			case dfilter.OperatorEq:
				q = q.Where(fmt.Sprintf("%s->>'%s' = ?", field[0], field[1]), v.Value)
			case dfilter.OperatorNot:
				q = q.Where(fmt.Sprintf("%s->>'%s' != ?", field[0], field[1]), v.Value)
			default:
				q = q.Where(fmt.Sprintf("%s->>'%s' %s ?", field[0], field[1], v.SQLOperator), v.Value)
			}
			continue
		}
		if r, ok := query.GetReplaces()[v.Key]; ok {
			v.Key = r
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
			switch v.Operator {
			case dfilter.OperatorNot:
				q = q.Where(fmt.Sprintf("%s IS NOT NULL", v.Key))
			case dfilter.OperatorIs:
				q = q.Where(fmt.Sprintf("%s IS NULL", v.Key))
			}
			continue
		}
		if v.ValueType == dfilter.Date {
			q = q.Where(fmt.Sprintf("date(%s) %s ?", v.Key, v.SQLOperator), v.Value)
			continue
		}
		if v.Operator == dfilter.OperatorCt {
			q = q.Where(fmt.Sprintf("%s %s ?", v.Key, v.SQLOperator), "%"+v.Value.(string)+"%")
			continue
		}
		q = q.Where(fmt.Sprintf("%s %s ?", v.Key, v.SQLOperator), v.Value)
	}
	return
}
