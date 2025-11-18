package middleware

import (
	"github.com/Maximumsoft-Co-LTD/otelgo/eto"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

func OTELMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := eto.Propagate().FromHTTPRequest(c.Request)

		ctx, span := eto.Trace().
			Name(c.FullPath()).
			FromContext(ctx).
			Kind(trace.SpanKindServer).
			Attr("http.method", c.Request.Method).
			Attr("http.route", c.FullPath()).
			Start()
		defer span.End()

		c.Request = c.Request.WithContext(ctx)

		c.Next()

		eto.Propagate().
			FromContext(ctx).
			ToHTTPResponse(c.Writer)
	}
}
