package checkapi

import (
	"context"
	"net/http"

	"github.com/dubey22rohit/togo-service/app/domain/checkapp"
	"github.com/dubey22rohit/togo-service/foundation/web"
)

type api struct {
	checkApp *checkapp.App
}

func newAPI(checkapp *checkapp.App) *api {
	return &api{
		checkApp: checkapp,
	}
}

func (api *api) liveness(ctx context.Context, r *http.Request) web.Encoder {
	return api.checkApp.Liveness()
}
