package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/httpapi/middleware"
)

func TestAdminAuth(t *testing.T) {
	log := zerolog.New(os.Stderr)
	const testKey = "test-secret-key"

	// Следующий хендлер — просто отвечает 200
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.AdminAuth(testKey, &log)(next)

	t.Run("No auth key — 401", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/auth/blacklist", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("Bad key— 401", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/auth/blacklist", nil)
		req.Header.Set("X-Admin-Key", "wrong-key")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("Good key — 200", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/auth/blacklist", nil)
		req.Header.Set("X-Admin-Key", testKey)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	})
}
