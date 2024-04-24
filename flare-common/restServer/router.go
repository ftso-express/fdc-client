package restServer

import (
	"log"
	"net/http"

	swagger "github.com/davidebianchi/gswagger"
	"github.com/davidebianchi/gswagger/support/gorilla"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
	"github.com/swaggest/swgui"
	v3 "github.com/swaggest/swgui/v3"
)

type ErrorHandler struct {
	Handler func(w http.ResponseWriter)
}

type Router interface {
	AddRoute(path string, handler RouteHandler, description ...string)
	WithPrefix(prefix string, tag string) Router
	Finalize()
}

// Default router implementation using gorilla/mux
type defaultRouter struct {
	router *mux.Router
}

func (r *defaultRouter) AddRoute(path string, handler RouteHandler, description ...string) {
	r.router.HandleFunc(path, handler.Handler).Methods(handler.Method)
}

func (r *defaultRouter) WithPrefix(prefix string, tag string) Router {
	return &defaultRouter{
		router: r.router.PathPrefix(prefix).Subrouter(),
	}
}

func (r *defaultRouter) Finalize() {
}

func NewDefaultRouter(mRouter *mux.Router) Router {
	return &defaultRouter{
		router: mRouter,
	}
}

// Router implementation with swagger support
type swaggerRouter struct {
	mRouter *mux.Router
	router  *swagger.Router[gorilla.HandlerFunc, *mux.Route]
	tag     string
}

func NewSwaggerRouter(mRouter *mux.Router, title string, version string) Router {
	router, _ := swagger.NewRouter(gorilla.NewRouter(mRouter), swagger.Options{
		Openapi: &openapi3.T{
			Info: &openapi3.Info{
				Title:   title,
				Version: version,
			},
		},
	})
	return &swaggerRouter{
		mRouter: mRouter,
		router:  router,
		tag:     "",
	}
}

// Add a route to the router and generate openapi definitions from the handler
// The first item in the description parameter is used to set the openapi summary field and
// the second item is used to set the openapi description field
func (r *swaggerRouter) AddRoute(path string, handler RouteHandler, description ...string) {
	swaggerDefinitions := handler.SwaggerDefinitions
	swaggerDefinitions.Tags = []string{r.tag}
	if len(description) > 0 {
		swaggerDefinitions.Summary = description[0]
		if len(description) > 1 {
			swaggerDefinitions.Description = description[1]
		}
	}

	_, err := r.router.AddRoute(handler.Method, path, handler.Handler, swaggerDefinitions)
	if err != nil {
		log.Fatal(err)
	}
}

func (r *swaggerRouter) WithPrefix(prefix string, tag string) Router {
	mSubRouter := r.mRouter.NewRoute().Subrouter()
	subRouter, _ := r.router.SubRouter(gorilla.NewRouter(mSubRouter), swagger.SubRouterOptions{
		PathPrefix: prefix,
	})
	return &swaggerRouter{
		mRouter: mSubRouter,
		router:  subRouter,
		tag:     tag,
	}
}

func (r *swaggerRouter) Finalize() {
	if err := r.router.GenerateAndExposeOpenapi(); err != nil {
		log.Fatal(err)
	}

	config := swgui.Config{
		Title:       "FDC protocol data provider API",
		SwaggerJSON: "/documentation/json",
		BasePath:    "/swagger",
	}

	handler := v3.NewHandlerWithConfig(config)
	r.mRouter.PathPrefix("/swagger").HandlerFunc(handler.ServeHTTP)
}
