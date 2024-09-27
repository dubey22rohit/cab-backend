package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/dubey22rohit/togo-service/api/sdk/http/debug"
	"github.com/dubey22rohit/togo-service/foundation/logger"
)

var build = "develop"

func main() {
	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "********* SEND ALERT *********")
		},
	}

	traceIdFunc := func(ctx context.Context) string {
		/**
		* TODO: Implement
		 */
		return ""
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "AUTH", traceIdFunc, events)

	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "error starting service", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	/*
	* -------------------------------------------------------------------------
	* NUMCPU
	 */

	log.Info(ctx, "startup", "NUMCPU", runtime.NumCPU())

	/*
	* -------------------------------------------------------------------------
	* Configuration
	 */

	config := struct {
		conf.Version
		Web struct {
			ReadTimeout        time.Duration `conf:"default:5s"`
			WriteTimeout       time.Duration `conf:"default:10s"`
			IdleTimeout        time.Duration `conf:"default:120s"`
			ShutdownTimeout    time.Duration `conf:"default:20s"`
			APIHost            string        `conf:"default:0.0.0.0:6000"`
			DebugHost          string        `conf:"default:0.0.0.0:6010"`
			CORSAllowedOrigins []string      `conf:"default:*"`
		}
		Auth struct {
			KeysEnvVar string
			KeysFolder string `conf:"default:zarf/keys/"`
			ActiveKID  string `conf:"default:54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"`
			Issuer     string `conf:"default:service project"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "Auth",
		},
	}

	const prefix = "AUTH"
	help, err := conf.Parse(prefix, &config)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config %s: %w", prefix, err)
	}

	/*
	* -------------------------------------------------------------------------
	* App Starting
	 */

	log.Info(ctx, "starting auth service", "version", config.Build)
	defer log.Info(ctx, "auth shutdown complete")

	out, err := conf.String(&config)
	if err != nil {
		return fmt.Errorf("error generating config for output: %w", err)
	}
	log.Info(ctx, "auth startup", "config", out)

	log.BuildInfo(ctx)

	expvar.NewString("build").Set(config.Build)

	/*
		* -------------------------------------------------------------------------
		* Database Support
		TODO
	*/

	/*
		* -------------------------------------------------------------------------
		* Initialize Authentication Support
		TODO
	*/

	/*
		* -------------------------------------------------------------------------
		* Initialize Tracing Support
		TODO
	*/

	/*
	* -------------------------------------------------------------------------
	* Start Debug Service
	 */

	go func() {
		log.Info(ctx, "debug service startup", "status", "debug v1 router started", "host", config.Web.DebugHost)
		if err := http.ListenAndServe(config.Web.DebugHost, debug.Mux()); err != nil {
			log.Error(ctx, "shutdown", "status", "debug v1 router closed", "host", config.Web.DebugHost, "msg", err)
		}
	}()

	/*
	* -------------------------------------------------------------------------
	* Start API Service
	 */

	log.Info(ctx, "startup", "status", "initializing V1 API support")
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

	/*
		TODO: Make cfgMux and http server
	*/

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", nil)
		// serverErrors <- api.ListenAndServe()

	}()

	select {
	// This blocks until the server errors channel receives an error.
	// This could be due to the server being unable to start or an error
	// occurred while the server was running.
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	// This line blocks until a signal is received from the OS. The signal is
	// generated when the process receives a SIGTERM or SIGINT signal. This
	// allows the program to gracefully shut down when a signal is received.
	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)
	}

	return nil
}
