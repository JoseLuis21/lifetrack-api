package middleware

import (
	"net/http"

	muxhandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/neutrinocorp/lifetrack-api/pkg/transport/resiliency"
)

// InjectHTTP injects to the given HTTP router predefined middlewares such as logging, CORS, compression and rate limit
func InjectHTTP(router *mux.Router) {
	// router.Use(muxhandlers.RecoveryHandler(muxhandlers.RecoveryLogger(logger)))
	router.Use(muxhandlers.CORS(
		muxhandlers.AllowedMethods([]string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPatch,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		}),
		muxhandlers.AllowedOrigins([]string{"*"}),
	))
	router.Use(muxhandlers.CompressHandler)
	router.Use(resiliency.HTTPRateLimit)
}
