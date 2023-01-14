package test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/letsfire/go-pkgs/data"
	"github.com/letsfire/go-pkgs/data/orm"
)

var ctx context.Context
var postgres orm.Interface

func init() {
	pg := data.Postgres{
		Address:  "mlszp.com:5432",
		Username: "postgres",
		Password: "Mangkaixin666!",
	}
	postgres = orm.Postgres(pg.GetConn("manghi"))
	ctx = context.WithValue(context.TODO(), "request_id", time.Now().Unix())
}

type Wallet struct {
	tableName struct{} `pg:"pcenter.wallet,alias:w"`
	Id        string   `pg:"id,pk"`
	UserId    string
	RoleCode  string
	RoleName  string
	KindCode  string
	KindName  string
	Amount    int64
	Freeze    int64
	Version   int
}

func TestPostgresORM(t *testing.T) {
	cond := orm.Condition{
		Where: &orm.Where{
			Connect: "OR",
			Filters: []*orm.Filter{
				{Field: "id", Value: "2"},
				{Field: "amount", Value: "3", Symbol: ">="},
			},
			SubWhere: &orm.Where{
				Filters: []*orm.Filter{
					{Field: "freeze", Value: "2"},
					{Field: "version", Value: "1"},
				},
			},
		},
		Fields: []string{"*"},
		Paging: orm.Paging{
			PageNo:   1,
			PageSize: 2,
		},
		OrderBy: []orm.OrderBy{
			{Field: "id", Direction: "asc"},
			{Field: "amount", Direction: "desc"},
		},
	}
	var wms = make([]*Wallet, 0)
	fmt.Println(postgres.Search(ctx, &cond, &wms))
	fmt.Println(len(wms))
}

func TestUpdateLock(t *testing.T) {
	cmp := orm.CMap{"id": "mc8ivhipuv"}
	var wallet = make([]*Wallet, 0)
	var wg = new(sync.WaitGroup)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			postgres.UpdateLock(ctx, cmp.Where(nil), &wallet, func(res interface{}, uw orm.UnitWorkInterface) error {
				var model = *(res.(*[]*Wallet))
				for i := range model {
					model[i].Amount += 1
					model[i].Freeze += 2
					uw.Update(model[i], "amount")
					uw.Update(model[i], "freeze")
				}
				return uw.ExecuteCheckAffected(len(model) * 2)
			})
		}()
	}
	wg.Wait()
	fmt.Println(wallet[0].Amount)
}

func TestDelete(t *testing.T) {
	ws := []*Wallet{{Id: "1"}}
	fmt.Println(postgres.Delete(ctx, &ws))
}
