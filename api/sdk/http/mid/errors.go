package mid

import (
	"context"
	"net/http"

	"github.com/dubey22rohit/togo-service/app/sdk/mid"
	"github.com/dubey22rohit/togo-service/foundation/logger"
	"github.com/dubey22rohit/togo-service/foundation/web"
)

func Errors(log *logger.Logger) web.MidFunc {
	midFunc := func(ctx context.Context, r *http.Request, next mid.HandlerFunc) mid.Encoder {
		return mid.Errors(ctx, log, next)
	}

	return addMidFunc(midFunc)
}
