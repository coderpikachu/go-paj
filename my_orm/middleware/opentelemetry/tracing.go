package opentelemetry

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"paj/my_orm"
)

const defaultInstrumentationName = "middleware/opentelemetry"

type MiddlewareBuilder struct {
	Tracer trace.Tracer
}

func (b MiddlewareBuilder) Build() my_orm.Middleware {
	if b.Tracer == nil {
		b.Tracer = otel.GetTracerProvider().Tracer(defaultInstrumentationName)
	}
	return func(next my_orm.HandleFunc) my_orm.HandleFunc {
		return func(ctx context.Context, qc *my_orm.QueryContext) *my_orm.QueryResult {
			tbl := qc.Model.TableName
			reqCtx, span := b.Tracer.Start(ctx, qc.Type+"-"+tbl, trace.WithAttributes())
			defer span.End()
			span.SetAttributes(attribute.String("component", "my_orm"))
			q, err := qc.Query()
			if err != nil {
				span.RecordError(err)
			}
			span.SetAttributes(attribute.String("table", tbl))
			if q != nil {
				span.SetAttributes(attribute.String("sql", q.SQL))
			}
			return next(reqCtx, qc)
		}
	}
}
