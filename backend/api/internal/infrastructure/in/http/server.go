package http

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/lregnier/design-youtube/api/gen/api"
)

func newRouter(h *Handler, uploadSecret string, corsAllowedOrigins []string) http.Handler {
	secretMw := UploadSecretMiddleware(uploadSecret)
	strictHandler := api.NewStrictHandlerWithOptions(h, []api.StrictMiddlewareFunc{secretMw}, api.StrictHTTPServerOptions{})

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: corsAllowedOrigins,
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "X-Upload-Secret"},
	}))
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	api.HandlerFromMux(strictHandler, r)

	return r
}

type Server struct {
	addr   string
	router http.Handler
}

func NewServer(h *Handler, uploadSecret string, corsAllowedOrigins []string, addr string) *Server {
	return &Server{addr: addr, router: newRouter(h, uploadSecret, corsAllowedOrigins)}
}

func (s *Server) Start() error {
	log.Println("listening on " + s.addr)
	return http.ListenAndServe(s.addr, s.router)
}
