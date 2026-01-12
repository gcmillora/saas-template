package webserver

import (
	"adobo/config"
	"adobo/generated/oapi"
	oapi_public "adobo/generated/oapi/public"
	"adobo/internal/webserver/handler"
	"adobo/internal/webserver/middleware"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

type Webserver struct {
	router     *chi.Mux
	serverAddr string
}

func (ws *Webserver) Router() *chi.Mux {
	return ws.router
}

func NewWebserver(app *config.App) *Webserver {
	handler := handler.NewHandler(app)
	serverAddr := ":" + app.EnvVars().ServerPort()

	r := chi.NewRouter()
	r.Use(middleware.NewLoggerMiddleware())
	r.Use(chimiddleware.Recoverer)
	r.Get("/health", handler.GetHealth)

	// Authenticated API routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.NewContextInjectorMiddleware())
		baseURL := "/api/public/v1"
		strictHandler := oapi_public.NewStrictHandlerWithOptions(
			handler,
			[]oapi_public.StrictMiddlewareFunc{},
			oapi_public.StrictHTTPServerOptions{
				RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
					http.Error(w, err.Error(), http.StatusBadRequest)
				},
				ResponseErrorHandlerFunc: middleware.HandleErrorWithLog(app),
			},
		)
		oapi_public.HandlerFromMuxWithBaseURL(strictHandler, r, baseURL)
	})

	// Authenticated API routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.NewContextInjectorMiddleware())
		r.Use(middleware.NewAuthMiddleware(app))
		baseURL := "/api/v1"
		strictHandler := oapi.NewStrictHandlerWithOptions(
			handler,
			[]oapi.StrictMiddlewareFunc{},
			oapi.StrictHTTPServerOptions{
				RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
					http.Error(w, err.Error(), http.StatusBadRequest)
				},
				ResponseErrorHandlerFunc: middleware.HandleErrorWithLog(app),
			},
		)
		oapi.HandlerFromMuxWithBaseURL(strictHandler, r, baseURL)
	})

	return &Webserver{
		router:     r,
		serverAddr: serverAddr,
	}
}

func (ws *Webserver) Start() {

	log.Print("WebServer listening on " + ws.serverAddr)

	s := &http.Server{
		Handler: ws.router,
		Addr:    ws.serverAddr,
	}

	log.Fatal(s.ListenAndServe())
}

func (ws *Webserver) PrintRoutes() {
	err := chi.Walk(
		ws.router,
		func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
			fmt.Printf("[%s]: '%s' has %d middlewares\n", method, route, len(middlewares))
			return nil
		},
	)

	if err != nil {
		log.Panicln(err)
	}
}
