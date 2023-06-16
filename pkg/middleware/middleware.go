package middleware

import (
	"compress/flate"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"net/http"
)

func DefaultMiddleware(r *chi.Mux) http.Handler {
	compressor := middleware.NewCompressor(flate.DefaultCompression)
	r.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Recoverer,
		middleware.RequestID,
		middleware.RealIP,
		middleware.StripSlashes,
		compressor.Handler,
		cors.Handler(cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins: []string{"https://*", "http://*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		},
		))
	return r
}
