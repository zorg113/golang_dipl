package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/config"
	"github.com/zorg113/golang_dipl/atibruteforce/model/entity"
	"github.com/zorg113/golang_dipl/atibruteforce/model/service"
	mock_service "github.com/zorg113/golang_dipl/atibruteforce/store/adapters/mocks"
)

func TestAuthorization(t *testing.T) {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	controller := gomock.NewController(t)
	defer controller.Finish()
	blackListMockStor := mock_service.NewMockBlackListStore(controller)
	blacklist := service.NewBlackList(blackListMockStor, &logger)

	whiteListMockStor := mock_service.NewMockWhiteListStore(controller)
	whitelist := service.NewWhiteList(whiteListMockStor, &logger)
	cfg, err := config.NewConfig("./../../../../config/conf.yaml")
	if err != nil {
		require.NoError(t, err)
	}
	auth := service.NewAuthorization(blacklist, whitelist, &cfg, &logger)
	authHandler := NewAuthorization(auth, &logger)
	cases := []struct {
		name    string
		request entity.Request
	}{
		{name: "test request", request: entity.Request{
			Login:    "admin",
			Password: "password",
			Ip:       "127.0.0.1",
		}},
	}
	blackListMockStor.EXPECT().GetIPs().Return([]entity.IpNetwork{}, nil).AnyTimes()
	whiteListMockStor.EXPECT().GetIPs().Return([]entity.IpNetwork{}, nil).AnyTimes()
	router := mux.NewRouter()
	router.HandleFunc("/auth/check", authHandler.AuthorizationHanler).Methods("POST")
	request := cases[0].request
	body, err := json.Marshal(request)
	require.NoError(t, err)
	req, err := http.NewRequest("POST", "/auth/check", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	ss := httptest.NewRecorder()
	router.ServeHTTP(ss, req)
	require.Equal(t, http.StatusOK, ss.Code)
	s := ss.Body.String()
	require.Equal(t, "ok=true", s)

}
