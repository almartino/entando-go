package tracing

import (
	"context"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace/noop"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

// Middleware must be used to provide tracing capabilities to the MS.
//
// This interface embed io.Closer which must be deferred (usually in the main) to
// correctly shut down the tracer provider.
type Middleware interface {
	io.Closer
	// Middleware start a span for each request received by the http.Handler.
	// The resulting span will be injected into the request context for tracing propagation purposes.
	//
	// If tracing is disabled h is left untouched.
	Middleware(h http.Handler) http.Handler
}

type zipkinMiddleware struct {
	cleanup        func() error
	excludedPathRE *regexp.Regexp
	enabled        bool
}

// New initializes the global tracer and return a Middleware.
// This middleware is configured by using the following env variables:
//
//   - OTEL_EXPORTERURL: the url which points to the zipkin endpoint: [default: "http://localhost:9411/api/v2/spans"]
//   - OTEL_FILTEREDROUTES: a valid regular expression which when matches a valid request url path skips the span generation [default: "/api/health|/metrics"]
//   - OTEL_SERVICENAME: the microservice name [default: "ms"]
//   - OTEL_ENABLED: whether to enable tracing [default: false]
func New() Middleware {
	var c otelConfiguration
	envconfig.MustProcess("OTEL", &c)
	var logger = log.New(os.Stderr, "zk", log.Ldate|log.Ltime|log.Llongfile)
	shutdown, err := c.initTracer(logger)
	if err != nil {
		panic(err)
	}
	return zipkinMiddleware{
		cleanup: func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			return shutdown(ctx)
		},
		excludedPathRE: regexp.MustCompile(c.FilteredRoutes),
		enabled:        c.Enabled,
	}
}

func (o zipkinMiddleware) Middleware(h http.Handler) http.Handler {
	if !o.enabled {
		return h
	}
	tr := otel.GetTracerProvider().Tracer("http-server")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if o.excludedPathRE.MatchString(path) {
			h.ServeHTTP(w, r)
			return
		}
		spanName := fmt.Sprintf("%s %s", r.Method, r.URL.Path)
		ctx, span := tr.Start(r.Context(), spanName)
		defer span.End()
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (o zipkinMiddleware) Close() error {
	return o.cleanup()
}

type otelConfiguration struct {
	ExporterURL    string `default:"http://localhost:9411/api/v2/spans"`
	FilteredRoutes string `default:"/api/health|/metrics"`
	ServiceName    string `default:"ms"`
	Enabled        bool   `default:"false"`
}

// initTracer creates a new trace provider instance and registers it as global trace provider.
// when Enabled is false the noop tracer is registered.
func (o otelConfiguration) initTracer(logger *log.Logger) (func(context.Context) error, error) {
	if !o.Enabled {
		tp := noop.NewTracerProvider()
		otel.SetTracerProvider(tp)
		return func(ctx context.Context) error {
			return nil
		}, nil
	}
	exporter, err := zipkin.New(
		o.ExporterURL,
		zipkin.WithLogger(logger),
	)
	if err != nil {
		return nil, err
	}
	batch := trace.NewBatchSpanProcessor(exporter)
	tp := trace.NewTracerProvider(
		trace.WithSpanProcessor(batch),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(o.ServiceName),
		)),
	)
	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}
