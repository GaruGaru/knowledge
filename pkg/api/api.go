package api

import (
	"github.com/garugaru/knowledge/pkg/data"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	muxprom "gitlab.com/msvechla/mux-prometheus/pkg/middleware"
	"net/http"
	"time"
)

const (
	DefaultServeAddr     = "0.0.0.0:8000"
	DefaultServerTimeout = 15 * time.Second
)

type ServeOpts struct {
	Addr    string
	Timeout time.Duration
}

type Config struct {
	EnableMetrics bool
}

type Api struct {
	catalog data.Catalog
	config  Config
}

func New(config Config, catalog data.Catalog) *Api {
	return &Api{catalog: catalog, config: config}
}

func (a Api) router() *mux.Router {
	router := mux.NewRouter()

	router.PathPrefix("/catalog").Handler(a.catalogRouter())
	router.HandleFunc("/healthz", a.healthz).Methods(http.MethodGet)

	if a.config.EnableMetrics {
		router.Use(muxprom.NewDefaultInstrumentation().Middleware)
		router.Path("/metrics").Handler(promhttp.Handler())
	}

	return router
}

func (a Api) Server(opts ServeOpts) *http.Server {
	if len(opts.Addr) == 0 {
		opts.Addr = DefaultServeAddr
	}

	if opts.Timeout == 0 {
		opts.Timeout = DefaultServerTimeout
	}

	return &http.Server{
		Handler:      a.router(),
		Addr:         opts.Addr,
		WriteTimeout: opts.Timeout,
		ReadTimeout:  opts.Timeout,
	}
}
