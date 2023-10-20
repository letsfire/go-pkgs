package orm

import (
	"context"
	"reflect"
	"strings"
)

type Interface interface {
	Search(ctx context.Context, cond *Condition, res interface{}) (int, error)
	Insert(ctx context.Context, model interface{}) (int, error)
	InsertUpdate(ctx context.Context, model interface{}, cols ...string) (int, error)
	Delete(ctx context.Context, model interface{}) (int, error)
	DeleteCond(ctx context.Context, model interface{}, cond *Condition) (int, error)
	DeleteForce(ctx context.Context, model interface{}) (int, error)
	Update(ctx context.Context, model interface{}, cols ...string) (int, error)
	UpdateSafe(ctx context.Context, model interface{}, locker *Locker, cols ...string) (int, error)
	UpdateCond(ctx context.Context, model interface{}, kvs map[string]interface{}, cond *Condition) (int, error)
	UpdateLock(
		ctx context.Context, w *Where, res interface{},
		fn func(res interface{}, uw UnitWorkInterface) error,
	) error
	Increase(ctx context.Context, model interface{}, kvs map[string]int64, cond *Condition) (int, error)
	UnitWork(ctx context.Context) UnitWorkInterface
}

type UnitWorkInterface interface {
	Insert(models ...interface{}) UnitWorkInterface
	Delete(models ...interface{}) UnitWorkInterface
	DeleteForce(models ...interface{}) UnitWorkInterface
	DeleteCond(model interface{}, cond *Condition) UnitWorkInterface
	Update(model interface{}, cols ...string) UnitWorkInterface
	UpdateSafe(model interface{}, locker *Locker, cols ...string) UnitWorkInterface
	UpdateCond(model interface{}, kvs map[string]interface{}, cond *Condition) UnitWorkInterface
	Increase(model interface{}, kvs map[string]int64, cond *Condition) UnitWorkInterface
	Execute(check func(i, u, d, incr int) bool) error
	ExecuteCheckAffected(affected int) error
}

// Condition 查询组合
type Condition struct {
	Where       *Where
	Fields      []string
	Paging      Paging
	GroupBy     []string
	OrderBy     []OrderBy
	FieldExps   []string
	Relations   map[string]*Condition
	WithDeleted bool
}

func (c *Condition) Deleted() *Condition {
	c.WithDeleted = true
	return c
}

func (c *Condition) Limit(num int) *Condition {
	return c.SetPage(0, num)
}

func (c *Condition) LimitNoCount(num int) *Condition {
	return c.SetPageNoCount(0, num)
}

func (c *Condition) SetPage(no, size int) *Condition {
	c.Paging = Paging{PageNo: no, PageSize: size}
	return c
}

func (c *Condition) SetPageNoCount(no, size int) *Condition {
	c.Paging = Paging{PageNo: no, PageSize: size, noCount: true}
	return c
}

func (c *Condition) AddOrder(field, direction string) *Condition {
	c.OrderBy = append(c.OrderBy, OrderBy{Field: field, Direction: direction})
	return c
}

func (c *Condition) AddOrderExp(express string) *Condition {
	if express != "" {
		c.OrderBy = append(c.OrderBy, OrderBy{Express: express})
	}
	return c
}

func (c *Condition) SetGroup(cols ...string) *Condition {
	c.GroupBy = cols
	return c
}

func (c *Condition) AddRelation(name string, cond *Condition) *Condition {
	if c.Relations == nil {
		c.Relations = make(map[string]*Condition, 0)
	}
	c.Relations[name] = cond
	return c
}

// Where 查询条件
type Where struct {
	Connect  string // AND, OR
	Filters  []*Filter
	SubWhere *Where
}

func (w *Where) Condition(fields ...string) *Condition {
	return &Condition{Where: w, Fields: fields}
}

func OrWhere(sub *Where, filters ...*Filter) *Where {
	return &Where{
		Connect:  "OR",
		SubWhere: sub,
		Filters:  filters,
	}
}

func AndWhere(sub *Where, filters ...*Filter) *Where {
	return &Where{
		Connect:  "AND",
		SubWhere: sub,
		Filters:  filters,
	}
}

// CMap 简单条件
type CMap map[string]interface{}

func (cm CMap) Where(sub *Where) *Where {
	fs := make([]*Filter, 0)
	for k, v := range cm {
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Bool && rv.IsZero() {
			continue // skip zero filter, but bool
		}
		if rv.Kind() == reflect.Slice {
			fs = append(fs, &Filter{Field: k, Value: v, Symbol: "in"})
		} else {
			fs = append(fs, &Filter{Field: k, Value: v, Symbol: "="})
		}
	}
	return AndWhere(sub, fs...)
}

func (cm CMap) And(filters ...*Filter) *Where {
	return cm.Where(AndWhere(nil, filters...))
}

func (cm CMap) Condition(fields ...string) *Condition {
	return cm.Where(nil).Condition(fields...)
}

// TimeRange 时间区间
type TimeRange struct {
	Start int64
	End   int64
}

func (t *TimeRange) Filters(field string) Filters {
	return Filters{
		{Field: field, Value: t.Start, Symbol: ">=", UseZero: false},
		{Field: field, Value: t.End, Symbol: "<=", UseZero: false},
	}
}

// Filter 过滤项
type Filter struct {
	Field   string      // 字段名
	Value   interface{} // 查询值
	Symbol  string      // 比较符
	Express string      // 表达式
	UseZero bool        // 查空值
}

func (f *Filter) split(sep string) interface{} {
	if v, ok := f.Value.(string); ok {
		return strings.Split(v, sep)
	}
	return f.Value
}

func (f *Filter) string() string { return f.Value.(string) }

func (f *Filter) Condition(fields ...string) *Condition {
	return AndWhere(nil, f).Condition(fields...)
}

// Filters 过滤集合
type Filters []*Filter

func (fs Filters) Or(fs2 ...*Filter) *Where {
	if len(fs2) == 0 {
		return OrWhere(nil, fs...)
	}
	return OrWhere(AndWhere(nil, fs2...), fs...)
}

func (fs Filters) And() *Where { return AndWhere(nil, fs...) }

// Paging 分页
type Paging struct {
	PageNo   int // 当前页码，默认1
	PageSize int // 每页数量，<0 仅查询数量，=0 查询全部
	noCount  bool
}

func (p *Paging) NoCount() { p.noCount = true }

func (p *Paging) onlyCount() bool  { return p.PageSize < 0 }
func (p *Paging) onlySelect() bool { return p.PageSize == 0 || p.noCount }

// OrderBy 排序
type OrderBy struct {
	Field     string // 排序字段
	Express   string // 排序表达式
	Direction string // ASC, DESC
}

// Locker 乐观锁
type Locker struct {
	Field string // 字段
	Value int64  // 比较值
	Next  func() int64
}
