package data

import (
	"context"
	"fmt"
	"github.com/zinclabs/sdk-go-zincsearch"
)

type Zincsearch struct {
	Address  string
	Username string
	Password string
}

func (z *Zincsearch) api() *client.APIClient {
	return client.NewAPIClient(&client.Configuration{
		Debug:   false,
		Servers: client.ServerConfigurations{{URL: z.Address}},
	})
}

func (z *Zincsearch) ctx(p context.Context) context.Context {
	auth := client.BasicAuth{UserName: z.Username, Password: z.Password}
	return context.WithValue(p, client.ContextBasicAuth, auth)
}

func (z *Zincsearch) Search(ctx context.Context, index, query string) ([]map[string]interface{}, int, error) {
	q := client.MetaZincQuery{
		Query: &client.MetaQuery{
			Match: &map[string]client.MetaMatchQuery{
				"name": {
					Query: &query,
				},
			},
		},
	}
	res, _, err := z.api().Search.Search(z.ctx(ctx), index).Query(q).Execute()
	fmt.Println(*res.Hits.Total.Value, err)
	return nil, int(*res.Hits.Total.Value), err
}

func (z *Zincsearch) Delete(ctx context.Context, index, id string) error {
	_, _, err := z.api().Document.Delete(z.ctx(ctx), index, id).Execute()
	return err
}

func (z *Zincsearch) IndexWithID(ctx context.Context, index, id string, record map[string]interface{}) error {
	_, _, err := z.api().Document.IndexWithID(z.ctx(ctx), index, id).Document(record).Execute()
	return err
}
