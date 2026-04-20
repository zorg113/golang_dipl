package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/rs/zerolog"
)

// проверка входа администратора.
func AdminAuth(apiKey string, log *zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("X-Admin-Key")
			if subtle.ConstantTimeCompare([]byte(key), []byte(apiKey)) != 1 {
				log.Warn().
					Str("ip", r.RemoteAddr).
					Str("path", r.URL.Path).Msg("unauthorized admin request")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
