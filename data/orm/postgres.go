package orm

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type postgres struct{ db *pg.DB }

func (p *postgres) Search(ctx context.Context, cond *Condition, res interface{}) (n int, e error) {
	q := p.db.ModelContext(ctx, res)
	if cond == nil {
		return 0, p.skipNoRows(q.Column("*").Select())
	}
	if cond.WithDeleted {
		q.AllWithDeleted()
	}
	if cond.Paging.onlyCount() {
		n, e = cond.applyPostgresQuery(q).Count()
	} else if cond.Paging.onlySelect() {
		n, e = 0, cond.applyPostgresQuery(q).Select()
	} else {
		n, e = cond.applyPostgresQuery(q).SelectAndCount()
	}
	return n, p.skipNoRows(e)
}

func (p *postgres) Insert(ctx context.Context, m interface{}) (int, error) {
	res, err := p.db.ModelContext(ctx, m).Returning("NULL").Insert()
	return p.rowsAffected(res), err
}

func (p *postgres) InsertUpdate(ctx context.Context, model interface{}, cols ...string) (int, error) {
	q := p.db.ModelContext(ctx, model)
	for _, col := range cols {
		q.Set("?0 = EXCLUDED.?0", pg.Ident(col))
	}
	var holders = make([]string, 0)
	var columns = make([]interface{}, 0)
	for _, pk := range q.TableModel().Table().PKs {
		holders = append(holders, "?")
		columns = append(columns, pg.Ident(pk.SQLName))
	}
	sql := fmt.Sprintf("(%s) DO UPDATE", strings.Join(holders, ","))
	res, err := q.OnConflict(sql, columns...).Insert()
	return p.rowsAffected(res), err
}

func (p *postgres) Update(ctx context.Context, m interface{}, cols ...string) (int, error) {
	res, err := p.db.ModelContext(ctx, m).Column(cols...).WherePK().Update()
	return p.rowsAffected(res), err
}

func (p *postgres) Delete(ctx context.Context, m interface{}) (int, error) {
	q := p.db.ModelContext(ctx, m)
	res, err := deleteBugFixed(q).Delete()
	return p.rowsAffected(res), err
}

func (p *postgres) DeleteForce(ctx context.Context, m interface{}) (int, error) {
	q := p.db.ModelContext(ctx, m)
	res, err := deleteBugFixed(q).ForceDelete()
	return p.rowsAffected(res), err
}

func (p *postgres) DeleteCond(ctx context.Context, m interface{}, cond *Condition) (int, error) {
	if cond == nil {
		return 0, nil
	}
	res, err := cond.applyPostgresQuery(p.db.ModelContext(ctx, m)).Delete()
	return p.rowsAffected(res), err
}

func (p *postgres) UpdateSafe(ctx context.Context, m interface{}, l *Locker, cols ...string) (int, error) {
	res, err := l.applyPostgresQuery(p.db.ModelContext(ctx, m).WherePK()).Column(cols...).Update()
	return p.rowsAffected(res), err
}

func (p *postgres) UpdateCond(ctx context.Context, m interface{}, kvs map[string]interface{}, cond *Condition) (int, error) {
	q := p.db.ModelContext(ctx, m)
	for key, value := range kvs {
		q.Set("?0 = ?1", pg.Ident(key), value)
	}
	res, err := cond.applyPostgresQuery(q).Update()
	return p.rowsAffected(res), err
}

func (p *postgres) UpdateLock(
	ctx context.Context, w *Where, res interface{},
	fn func(res interface{}, uw UnitWorkInterface) error,
) error {
	return p.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		q := w.Condition().applyPostgresQuery(tx.Model(res))
		if err := q.For("UPDATE").Select(); err != nil {
			return err
		} else {
			return fn(res, postgresUnitWork(tx))
		}
	})
}

func (p *postgres) Increase(
	ctx context.Context, model interface{},
	kvs map[string]int64, cond *Condition,
) (int, error) {
	q := p.db.ModelContext(ctx, model)
	if cond == nil {
		q.WherePK()
	} else {
		cond.applyPostgresQuery(q)
	}
	for field, value := range kvs {
		q.Set("?0 = ?0 + ?1", pg.Ident(field), value)
	}
	res, err := q.Update()
	return p.rowsAffected(res), err
}

func (p *postgres) UnitWork(ctx context.Context) UnitWorkInterface {
	return postgresUnitWork(p.db.WithContext(ctx))
}

func (p *postgres) rowsAffected(res orm.Result) int {
	if res == nil {
		return 0
	}
	return res.RowsAffected()
}

func (p *postgres) skipNoRows(err error) error {
	if err == pg.ErrNoRows {
		return nil
	}
	return err
}

func Postgres(db *pg.DB) Interface { return &postgres{db: db} }
