package orm

import (
	"context"
	"errors"

	"github.com/go-pg/pg/v10"
)

// updateUnit 更新单元
type updateUnit struct {
	cols   []string    // 更新字段
	model  interface{} // 业务模型
	locker *Locker     // 乐观加锁
}

// condUpdateUnit 条件更新
type condUpdateUnit struct {
	model interface{}            // 业务模型
	kvMap map[string]interface{} // 更新数据
	cond  *Condition             // 更新条件
}

// increaseUnit 递增单元
type increaseUnit struct {
	model interface{}      // 业务模型
	kvMap map[string]int64 // 更新数据
	cond  *Condition       // 更新条件
}

// condDeleteUnit 条件删除
type condDeleteUnit struct {
	model interface{} // 业务模型
	cond  *Condition  // 删除条件
}

// UnitWork 工作单元
type UnitWork struct {
	db        pg.DBI
	inserts   []interface{}
	updates   []updateUnit
	condUps   []condUpdateUnit
	deletes   []interface{}
	forceDels []interface{}
	condDels  []condDeleteUnit
	increases []increaseUnit
}

func (uw *UnitWork) Insert(ms ...interface{}) UnitWorkInterface {
	uw.inserts = append(uw.inserts, ms...)
	return uw
}

func (uw *UnitWork) Delete(ms ...interface{}) UnitWorkInterface {
	uw.deletes = append(uw.deletes, ms...)
	return uw
}

func (uw *UnitWork) DeleteForce(ms ...interface{}) UnitWorkInterface {
	uw.forceDels = append(uw.forceDels, ms...)
	return uw
}

func (uw *UnitWork) DeleteCond(model interface{}, cond *Condition) UnitWorkInterface {
	uw.condDels = append(uw.condDels, condDeleteUnit{model: model, cond: cond})
	return uw
}

func (uw *UnitWork) Update(m interface{}, cols ...string) UnitWorkInterface {
	uw.updates = append(uw.updates, updateUnit{cols: cols, model: m})
	return uw
}

func (uw *UnitWork) UpdateSafe(m interface{}, locker *Locker, cols ...string) UnitWorkInterface {
	uw.updates = append(uw.updates, updateUnit{cols: cols, model: m, locker: locker})
	return uw
}

func (uw *UnitWork) UpdateCond(m interface{}, kvs map[string]interface{}, cond *Condition) UnitWorkInterface {
	uw.condUps = append(uw.condUps, condUpdateUnit{model: m, kvMap: kvs, cond: cond})
	return uw
}

func (uw *UnitWork) Increase(model interface{}, kvs map[string]int64, cond *Condition) UnitWorkInterface {
	uw.increases = append(uw.increases, increaseUnit{model: model, kvMap: kvs, cond: cond})
	return uw
}

func (uw *UnitWork) Execute(check func(i, u, d, incr int) bool) error {
	if tx, ok := uw.db.(*pg.Tx); ok {
		return uw.executeWithTx(tx, check)
	} else {
		return uw.db.RunInTransaction(context.TODO(), func(tx *pg.Tx) error {
			return uw.executeWithTx(tx, check)
		})
	}
}

func (uw *UnitWork) ExecuteCheckAffected(affected int) error {
	return uw.Execute(func(i, u, d, incr int) bool { return i+u+d+incr == affected })
}

func (uw *UnitWork) executeWithTx(tx *pg.Tx, check func(i, u, d, incr int) bool) error {
	var i, u, d, incr int // 影响行计数器
	for _, m := range uw.deletes {
		q := deleteBugFixed(tx.Model(m))
		if r, e := q.Delete(); e != nil {
			return e
		} else {
			d += r.RowsAffected()
		}
	}
	for _, m := range uw.forceDels {
		q := deleteBugFixed(tx.Model(m))
		if r, e := q.ForceDelete(); e != nil {
			return e
		} else {
			d += r.RowsAffected()
		}
	}
	for _, m := range uw.condDels {
		q := m.cond.applyPostgresQuery(tx.Model(m.model))
		if r, e := q.Delete(); e != nil {
			return e
		} else {
			d += r.RowsAffected()
		}
	}
	for _, m := range uw.inserts {
		if r, e := tx.Model(m).Returning("NULL").Insert(); e != nil {
			return e
		} else {
			i += r.RowsAffected()
		}
	}
	for _, m := range uw.updates {
		q := tx.Model(m.model).WherePK()
		if m.locker != nil {
			m.locker.applyPostgresQuery(q)
		}
		if r, e := q.Column(m.cols...).Update(); e != nil {
			return e
		} else {
			u += r.RowsAffected()
		}
	}
	for _, m := range uw.condUps {
		q := m.cond.applyPostgresQuery(tx.Model(m.model))
		for key, value := range m.kvMap {
			q.Set("?0 = ?1", pg.Ident(key), value)
		}
		if r, e := q.Update(); e != nil {
			return e
		} else {
			u += r.RowsAffected()
		}
	}
	for _, m := range uw.increases {
		q := tx.Model(m.model)
		if m.cond == nil {
			q.WherePK()
		} else {
			m.cond.applyPostgresQuery(q)
		}
		for field, value := range m.kvMap {
			q.Set("?0 = ?0 + ?1", pg.Ident(field), value)
		}
		if r, e := q.Update(); e != nil {
			return e
		} else {
			incr += r.RowsAffected()
		}
	}
	if check == nil || check(i, u, d, incr) {
		return nil
	}
	return errors.New("unit work reject by call check function")
}

func postgresUnitWork(db pg.DBI) UnitWorkInterface {
	return &UnitWork{
		db:      db,
		inserts: make([]interface{}, 0),
		updates: make([]updateUnit, 0),
		deletes: make([]interface{}, 0),
	}
}
