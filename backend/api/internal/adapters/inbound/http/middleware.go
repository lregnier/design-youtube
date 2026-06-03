package http

import (
	"context"
	"crypto/subtle"
	"net/http"

	"github.com/lregnier/design-youtube/api/internal/api"
)

var uploadOps = map[string]bool{
	"InitUpload":     true,
	"ConfirmChunk":   true,
	"CompleteUpload": true,
}

func UploadSecretMiddleware(secret string) api.StrictMiddlewareFunc {
	const msg = "missing or invalid upload secret"
	return func(f api.StrictHandlerFunc, operationID string) api.StrictHandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request, req any) (any, error) {
			if !uploadOps[operationID] {
				return f(ctx, w, r, req)
			}
			if subtle.ConstantTimeCompare([]byte(r.Header.Get("X-Upload-Secret")), []byte(secret)) != 1 {
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
