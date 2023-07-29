package data

import (
	"context"
	"github.com/zinclabs/sdk-go-zincsearch"
)

type Zincsearch struct {
	Address  string
	Username string
	Password string
}

func (z *Zincsearch) Api() *client.APIClient {
	return client.NewAPIClient(&client.Configuration{
		Debug:   false,
		Servers: client.ServerConfigurations{{URL: z.Address}},
	})
}

func (z *Zincsearch) Ctx(p context.Context) context.Context {
	auth := client.BasicAuth{UserName: z.Username, Password: z.Password}
	return context.WithValue(p, client.ContextBasicAuth, auth)
}

func (z *Zincsearch) Search(ctx context.Context, index, keyword string, fields ...string) ([]map[string]interface{}, int, error) {
	q := client.MetaZincQuery{
		Query: &client.MetaQuery{
			MultiMatch: &client.MetaMultiMatchQuery{
				Fields: fields, Query: &keyword,
			},
		},
		Sort: []string{"_sorce"},
	}
	if r, _, e := z.Api().Search.Search(z.Ctx(ctx), index).Query(q).Execute(); e != nil {
		return nil, 0, e
	} else {
		res := make([]map[string]interface{}, int(*r.Hits.Total.Value))
		for idx, hit := range r.Hits.Hits {
			res[idx] = map[string]interface{}{
				"id": *hit.Id,
			}
			delete(hit.Source, "@timestamp")
			for key, val := range hit.Source {
				res[idx][key] = val
			}
		}
		return res, int(*r.Hits.Total.Value), nil
	}
}

func (z *Zincsearch) Delete(ctx context.Context, index, id string) error {
	_, _, err := z.Api().Document.Delete(z.Ctx(ctx), index, id).Execute()
	return err
}

func (z *Zincsearch) IndexWithID(ctx context.Context, index, id string, record map[string]interface{}) error {
	_, _, err := z.Api().Document.IndexWithID(z.Ctx(ctx), index, id).Document(record).Execute()
	return err
}
