package webserver

import (
	"adobo/config"
	"adobo/generated/oapi"
	oapi_public "adobo/generated/oapi/public"
	"adobo/internal/webserver/handler"
	"adobo/internal/webserver/middleware"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{app.EnvVars().AppBaseUrl()},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Get("/health", handler.GetHealth)

	// Public API routes (rate limited)
	authRateLimiter := middleware.NewRateLimiter(1, 5) // 1 req/sec, burst of 5
	r.Group(func(r chi.Router) {
		r.Use(authRateLimiter.Middleware())
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
	s := &http.Server{
		Handler: ws.router,
		Addr:    ws.serverAddr,
	}

	go func() {
		log.Print("WebServer listening on " + ws.serverAddr)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Print("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Print("Server exited")
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
