package checkapp

import (
	"os"
	"runtime"

	"github.com/dubey22rohit/togo-service/foundation/logger"
)

// App manages the set of app layer api functions for the check domain.
type App struct {
	build string
	log   *logger.Logger
	// db    *sqlx.DB
}

// NewApp constructs a check app API for use.
func NewApp(build string, log *logger.Logger) *App {
	return &App{
		build: build,
		log:   log,
	}
}

// Liveness returns simple status info if the service is alive. If the
// app is deployed to a Kubernetes cluster, it will also return pod, node, and
// namespace details via the Downward API. The Kubernetes environment variables
// need to be set within your Pod/Deployment manifest.
func (a *App) Liveness() Info {
	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}

	info := Info{
		Status:     "up",
		Build:      a.build,
		Host:       host,
		Name:       os.Getenv("KUBERNETES_NAME"),
		PodIP:      os.Getenv("KUBERNETES_POD_IP"),
		Node:       os.Getenv("KUBERNETES_NODE_NAME"),
		Namespace:  os.Getenv("KUBERNETES_NAMESPACE"),
		GOMAXPROCS: runtime.NumCPU(),
	}

	// This handler provides a free timer loop.

	return info
}
