package orm

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/go-pg/pg/v10/types"
)

// applyPostgresQuery 应用乐观锁
func (l *Locker) applyPostgresQuery(query *orm.Query) *orm.Query {
	if l.Next != nil {
		next := strconv.FormatInt(l.Next(), 10)
		query.Column(l.Field).Value(l.Field, next)
	}
	return query.Where("?0 = ?1", pg.Ident(l.Field), l.Value)
}

// applyPostgresQuery 应用查询条件
func (c *Condition) applyPostgresQuery(query *orm.Query) *orm.Query {
	if len(c.Fields) == 0 {
		query.Column("*")
	} else {
		query.Column(c.Fields...)
	}
	c.Where.applyPostgresQuery(query)
	for _, orderBy := range c.OrderBy {
		orderBy.applyPostgresQuery(query)
	}
	if len(c.GroupBy) > 0 {
		query.Group(c.GroupBy...)
	}
	return c.Paging.applyPostgresQuery(query)
}

func (w *Where) applyPostgresQuery(query *orm.Query) *orm.Query {
	if w == nil {
		return query
	}
	for _, filter := range w.Filters {
		filter.applyPostgresQuery(query)
	}
	switch strings.ToUpper(w.Connect) {
	case "OR":
		query.WhereOrGroup(func(q *orm.Query) (*orm.Query, error) {
			return w.SubWhere.applyPostgresQuery(q), nil
		})
	case "OR NOT":
		query.WhereOrNotGroup(func(q *orm.Query) (*orm.Query, error) {
			return w.SubWhere.applyPostgresQuery(q), nil
		})
	case "NOT":
		query.WhereNotGroup(func(q *orm.Query) (*orm.Query, error) {
			return w.SubWhere.applyPostgresQuery(q), nil
		})
	default:
		query.WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			return w.SubWhere.applyPostgresQuery(q), nil
		})
	}
	return query
}

func (f *Filter) applyPostgresQuery(query *orm.Query) *orm.Query {
	if f.Express != "" {
		return query.Where(f.Express)
	}
	if f.UseZero == false {
		if reflect.ValueOf(f.Value).IsZero() {
			return query
		}
	}
	field := pg.Ident(f.Field)
	switch strings.ToUpper(f.Symbol) {
	case "!=":
		query.Where("?0 != ?1", field, f.Value)
	case ">":
		query.Where("?0 > ?1", field, f.Value)
	case ">=":
		query.Where("?0 >= ?1", field, f.Value)
	case "<":
		query.Where("?0 < ?1", field, f.Value)
	case "<=":
		query.Where("?0 <= ?1", field, f.Value)
	case "%L":
		query.Where("?0 LIKE ?1", field, "%"+f.string())
	case "L%":
		query.Where("?0 LIKE ?1", field, f.string()+"%")
	case "%%":
		query.Where("?0 LIKE ?1", field, "%"+f.string()+"%")
	case "IN":
		query.Where("?0 IN (?1)", field, types.In(f.split(",")))
	case "NOT IN":
		query.Where("?0 NOT IN (?1)", field, types.In(f.split(",")))
	case "HAS":
		query.Where("?0 = ANY(?1)", f.Value, field)
	default:
		query.Where("?0 = ?1", field, f.Value)
	}
	return query
}

func (p *Paging) applyPostgresQuery(query *orm.Query) *orm.Query {
	if p.PageSize == 0 {
		return query
	} else if p.PageSize < 0 {
		return query.Limit(-1)
	}
	if p.PageNo > 1 {
		query.Offset((p.PageNo - 1) * p.PageSize)
	}
	return query.Limit(p.PageSize)
}

func (o *OrderBy) applyPostgresQuery(query *orm.Query) *orm.Query {
	if o.Express != "" {
		return query.OrderExpr(o.Express)
	}
	return query.Order(fmt.Sprintf("%s %s", o.Field, o.Direction))
}

func deleteBugFixed(query *orm.Query) *orm.Query {
	tk := query.TableModel().Kind()
	if tk == reflect.Struct {
		query.WherePK()
	}
	return query
}
