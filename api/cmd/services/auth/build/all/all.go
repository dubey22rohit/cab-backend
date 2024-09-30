package all

import (
	"github.com/dubey22rohit/togo-service/api/domain/http/checkapi"
	"github.com/dubey22rohit/togo-service/api/sdk/http/mux"
	"github.com/dubey22rohit/togo-service/foundation/web"
)

func Routes() add {
	return add{}
}

type add struct{}

func (add) Add(app *web.App, cfg mux.Config) {
	// Construct the business domain packages we need here so we are using the
	// sames instances for the different set of domain apis.
	// delegate := delegate.New(cfg.Log)
	// userBus := userbus.NewBusiness(cfg.Log, delegate, usercache.NewStore(cfg.Log, userdb.NewStore(cfg.Log, cfg.DB), time.Minute))

	checkapi.Routes(app, checkapi.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		// DB:    cfg.DB,
	})

	// authapi.Routes(app, authapi.Config{
	// 	UserBus: userBus,
	// 	Auth:    cfg.Auth,
	// })
}
