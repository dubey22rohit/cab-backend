package mux

import (
	"context"
	"embed"
	"net/http"

	"github.com/dubey22rohit/togo-service/api/sdk/http/mid"
	"github.com/dubey22rohit/togo-service/foundation/logger"
	"github.com/dubey22rohit/togo-service/foundation/web"
	"go.opentelemetry.io/otel/trace"
)

// Options represent optional parameters.
type Options struct {
	corsOrigin []string
	static     *embed.FS
	staticDir  string
}

// WithCORS provides configuration options for CORS.
func WithCORS(origins []string) func(opts *Options) {
	return func(opts *Options) {
		opts.corsOrigin = origins
	}
}

// WithFileServer provides configuration options for file server.
func WithFileServer(static embed.FS, dir string) func(opts *Options) {
	return func(opts *Options) {
		opts.static = &static
		opts.staticDir = dir
	}
}

/*
TODO: create service specific configs as services grow
*/

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build string
	Log   *logger.Logger
	// DB *sqlx.DB
	Tracer trace.Tracer
	// App specific configs
}

// RouteAdder defines behavior that sets the routes to bind for an instance
// of the service.
type RouteAdder interface {
	Add(app *web.App, cfg Config)
}

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI(cfg Config, routeAdder RouteAdder, options ...func(opts *Options)) http.Handler {
	logger := func(ctx context.Context, msg string, args ...any) {
		cfg.Log.Info(ctx, msg, args...)
	}

	app := web.NewApp(
		logger,
		cfg.Tracer,
		mid.Errors(cfg.Log),
	)

	var opts Options
	for _, option := range options {
		option(&opts)
	}

	if len(opts.corsOrigin) > 0 {
		app.EnableCORS(opts.corsOrigin)
	}

	routeAdder.Add(app, cfg)

	if opts.static != nil {
		app.FileServer(*opts.static, opts.staticDir, http.NotFound)
	}

	return app
}
