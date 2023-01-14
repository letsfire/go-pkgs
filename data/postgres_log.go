package data

import (
	"context"
	"errors"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog/log"
)

// postgresLogger sql日志
type postgresLogger struct {
	slowLine time.Duration
}

func (pl *postgresLogger) BeforeQuery(ctx context.Context, qe *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (pl *postgresLogger) AfterQuery(ctx context.Context, qe *pg.QueryEvent) error {
	err := qe.Err
	sql_, _ := qe.FormattedQuery()
	cost := time.Since(qe.StartTime)
	if err == nil && cost > pl.slowLine {
		err = errors.New("slow query")
	}
	log.Err(err).Str("sql", string(sql_)).Dur("cost", cost).
		Interface("request_id", ctx.Value("request_id")).Msg("")
	return nil
}
