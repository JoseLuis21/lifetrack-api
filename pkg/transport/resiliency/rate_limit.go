package resiliency

import (
	"net/http"

	"golang.org/x/time/rate"
)

var httpLimiter = rate.NewLimiter(1, 50)

// HTTPRateLimit returns an http handler that performs request rate limiting
func HTTPRateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !httpLimiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
