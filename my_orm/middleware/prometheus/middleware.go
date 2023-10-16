package prometheus

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"paj/my_orm"
	"time"
)

type MiddlewareBuilder struct {
	Name        string
	Subsystem   string
	ConstLabels map[string]string
	Help        string
}

func (m MiddlewareBuilder) Build() my_orm.Middleware {
	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:        m.Name,
		Subsystem:   m.Subsystem,
		ConstLabels: m.ConstLabels,
		Help:        m.Help,
	}, []string{"type", "table"})

	return func(next my_orm.HandleFunc) my_orm.HandleFunc {
		return func(ctx context.Context, qc *my_orm.QueryContext) *my_orm.QueryResult {
			startTime := time.Now()
			defer func() {
				endTime := time.Now()
				typ := "unknown"
				// 原生查询才会走到这里
				tblName := "unknown"
				if qc.Model != nil {
					typ = qc.Model.TableName
					tblName = qc.Model.TableName
				}
				summaryVec.WithLabelValues(typ, tblName).
					Observe(float64(endTime.Sub(startTime).Milliseconds()))
			}()
			return next(ctx, qc)
		}
	}
}
