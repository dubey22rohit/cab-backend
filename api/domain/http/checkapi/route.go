package checkapi

import (
	"net/http"

	"github.com/dubey22rohit/togo-service/app/domain/checkapp"
	"github.com/dubey22rohit/togo-service/foundation/logger"
	"github.com/dubey22rohit/togo-service/foundation/web"
)

type Config struct {
	Build string
	Log   *logger.Logger
	// DB *sqlx.DB
}

func Routes(app *web.App, cfg Config) {
	const version = "v1"

	api := newAPI(checkapp.NewApp(cfg.Build, cfg.Log))
	app.HandlerFuncNoMid(http.MethodGet, version, "/liveness", api.liveness)
}
