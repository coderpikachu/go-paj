package querylog

import (
	"context"
	"log"
	"paj/my_orm"
)

type MiddlewareBuilder struct {
	logFunc func(sql string, args []any)
}

func (m *MiddlewareBuilder) LogFunc(logFunc func(sql string, args []any)) *MiddlewareBuilder {
	m.logFunc = logFunc
	return m
}

func NewBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		logFunc: func(sql string, args []any) {
			log.Println(sql, args)
		},
	}
}

func (m *MiddlewareBuilder) Build() my_orm.Middleware {
	return func(next my_orm.HandleFunc) my_orm.HandleFunc {
		return func(ctx context.Context, qc *my_orm.QueryContext) *my_orm.QueryResult {
			q, err := qc.Query()
			if err != nil {
				return &my_orm.QueryResult{
					Err: err,
				}
			}
			m.logFunc(q.SQL, q.Args)
			return next(ctx, qc)
		}
	}
}
