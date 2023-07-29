package data

import (
	"encoding/json"
	"os"
)

type Center struct {
	Redis      Redis
	Postgres   Postgres
	Rabbitmq   Rabbitmq
	Zincsearch Zincsearch
}

func (c *Center) prepare() *Center {
	if c.Redis.Address != "" {
		c.Redis.Cache()
	}
	if c.Postgres.Address != "" {
		c.Postgres.GetConn("postgres")
	}
	if c.Rabbitmq.Address != "" {
		c.Rabbitmq.GetConn("/")
	}
	return c
}

func LoadFromJsonFile(path string) *Center {
	var center = new(Center)
	bts, err := os.ReadFile(path)
	throwError(err, json.Unmarshal(bts, center))
	return center.prepare()
}
