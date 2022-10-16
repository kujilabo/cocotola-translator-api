package middleware

import "go.opentelemetry.io/otel"

var tracer = otel.Tracer("github.com/kujilabo/cocotola-translator-api/src/lib/controller/middleware")
