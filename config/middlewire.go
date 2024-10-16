package config

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/labstack/echo/v4"
)

func TraceMiddlewire(tp *trace.TracerProvider) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tracer := tp.Tracer("example-tracer")
			ctx, span := tracer.Start(c.Request().Context(), "handling request")
			span.AddEvent("event: handling / request")
			span.SetAttributes(attribute.String("handler", "root"))

			defer span.End()
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
