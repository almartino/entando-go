package main

import (
	"embed"
	"github.com/almartino/entando-go"
	"github.com/almartino/entando-go/examples/conference/internal/conference"
	"github.com/almartino/entando-go/examples/conference/internal/dbconn"
	"github.com/almartino/entando-go/oidc"
	"github.com/almartino/entando-go/tracing"
	"github.com/debyten/apierr"
	"github.com/debyten/httplayer"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"os"
)

var defaultCors = cors.New(cors.Options{
	AllowCredentials: true,
	AllowedMethods:   []string{http.MethodPost, http.MethodGet, http.MethodPatch, http.MethodPut, http.MethodDelete},
	AllowedHeaders:   []string{"Authorization", "Content-Type"},
	ExposedHeaders:   []string{"X-Total-Count"},
})

//go:embed migrations
var migrations embed.FS

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	apierr.DefaultDBNotFoundHandler = func(err error) bool {
		return errors.Is(err, gorm.ErrRecordNotFound)
	}
	// initialize tracing
	otlp := tracing.New()
	defer otlp.Close()
	// initialize database connection
	db := dbconn.New(migrations)
	// build the conference api server
	conferenceServer := conference.NewServer(logger, db)
	mux := entando.NewMux()                          // build the standard Entando serve mux
	oidcMW := oidc.New(oidc.PrincipalParser, logger) // keycloak oidc middleware with the default PrincipalParser
	httplayer.NewServiceBuilder(oidcMW).
		Add(conferenceServer).
		MountTo(mux)
	logger.Info("starting service")
	handler := otlp.Middleware(mux)        // register the otlp zipkin middleware
	handler = defaultCors.Handler(handler) // register cors
	handler = entando.ContextPath(handler) // register servlet context path middleware
	if err := http.ListenAndServe(":8081", handler); err != nil {
		panic(err)
	}
}
