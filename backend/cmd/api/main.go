package main

import (
	"context"
	"crypto/subtle"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/lregnier/design-youtube/backend/internal/api"
	"github.com/lregnier/design-youtube/backend/internal/config"
	"github.com/lregnier/design-youtube/backend/internal/handler"
	"github.com/lregnier/design-youtube/backend/internal/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := store.New(cfg)
	if err != nil {
		log.Fatalf("store: %v", err)
	}

	h := handler.New(cfg, db)

	secretMw := uploadSecretMiddleware(cfg.UploadSecret)
	strictHandler := api.NewStrictHandlerWithOptions(h, []api.StrictMiddlewareFunc{secretMw}, api.StrictHTTPServerOptions{})

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	api.HandlerFromMux(strictHandler, r)

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server: %v", err)
	}
}

var uploadOps = map[string]bool{
	"InitUpload":     true,
	"ConfirmChunk":   true,
	"CompleteUpload": true,
}

func uploadSecretMiddleware(secret string) api.StrictMiddlewareFunc {
	const msg = "missing or invalid upload secret"
	return func(f api.StrictHandlerFunc, operationID string) api.StrictHandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request, req any) (any, error) {
			if !uploadOps[operationID] {
				return f(ctx, w, r, req)
			}
			provided := r.Header.Get("X-Upload-Secret")
			if subtle.ConstantTimeCompare([]byte(provided), []byte(secret)) != 1 {
				switch operationID {
				case "InitUpload":
					return api.InitUpload401JSONResponse{Error: msg}, nil
				case "ConfirmChunk":
					return api.ConfirmChunk401JSONResponse{Error: msg}, nil
				case "CompleteUpload":
					return api.CompleteUpload401JSONResponse{Error: msg}, nil
				}
			}
			return f(ctx, w, r, req)
		}
	}
}
