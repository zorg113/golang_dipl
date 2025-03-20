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
	"github.com/zorg113/golang_dipl/atibruteforce/internal/common"
	"github.com/zorg113/golang_dipl/atibruteforce/model/entity"
	"github.com/zorg113/golang_dipl/atibruteforce/model/service"
	mock_service "github.com/zorg113/golang_dipl/atibruteforce/store/adapters/mocks"
)

func TestWhiteList_AddIP(t *testing.T) {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockStor := mock_service.NewMockWhiteListStore(controller)
	cases :=
		[]struct {
			name    string
			network entity.IpNetwork
		}{
			{name: "test valid ip mask", network: entity.IpNetwork{
				Ip:   "127.0.0.1",
				Mask: "255.255.0.0",
			}},
			{name: "test invalid ip mask", network: entity.IpNetwork{
				Ip:   "127.0.0.1",
				Mask: "256.255.0.0",
			}},
		}

	for _, testCase := range cases {
		prefix, _ := service.GetPrefix(testCase.network.Ip, testCase.network.Mask)
		mockStor.EXPECT().AddIP(prefix, testCase.network.Mask).Return(nil).MaxTimes(1)
		mockStor.EXPECT().AddIP(prefix, testCase.network.Mask).Return(common.IpAlreadyExist).AnyTimes()
	}
	whiteListService := service.NewWhiteList(mockStor, &logger)
	whiteList := NewWhiteList(whiteListService, &logger)

	router := mux.NewRouter()
	router.HandleFunc("/auth/whitelist", whiteList.AddIP).Methods("POST")

	ip := cases[0].network

	body, err := json.Marshal(ip)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/auth/whitelist", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	ss := httptest.NewRecorder()
	router.ServeHTTP(ss, req)

	require.Equal(t, http.StatusNoContent, ss.Code)

	req, err = http.NewRequest("POST", "/auth/whitelist", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	ss = httptest.NewRecorder()
	router.ServeHTTP(ss, req)

	require.Equal(t, http.StatusBadRequest, ss.Code)
	expected := "IP address already exists"
	require.Equal(t, expected, ss.Body.String())

	ivalidMask := cases[1].network
	body, err = json.Marshal(ivalidMask)
	require.NoError(t, err)

	req, err = http.NewRequest("POST", "/auth/whitelist", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	ss = httptest.NewRecorder()
	router.ServeHTTP(ss, req)

	require.Equal(t, http.StatusBadRequest, ss.Code)
}

func TestWhiteList_DeleteIP(t *testing.T) {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockStor := mock_service.NewMockWhiteListStore(controller)
	cases :=
		[]struct {
			name    string
			network entity.IpNetwork
		}{
			{name: "test valid ip mask", network: entity.IpNetwork{
				Ip:   "127.0.0.1",
				Mask: "255.255.0.0",
			}},
			{name: "test invalid ip mask", network: entity.IpNetwork{
				Ip:   "127.0.0.1",
				Mask: "256.255.0.0",
			}},
		}

	for _, testCase := range cases {
		prefix, _ := service.GetPrefix(testCase.network.Ip, testCase.network.Mask)
		mockStor.EXPECT().DeleteIP(prefix, testCase.network.Mask).Return(nil).AnyTimes()
	}
	whiteListService := service.NewWhiteList(mockStor, &logger)
	whiteList := NewWhiteList(whiteListService, &logger)

	router := mux.NewRouter()
	router.HandleFunc("/auth/whitelist", whiteList.DeleteIP).Methods("DELETE")

	ip := cases[0].network

	body, err := json.Marshal(ip)
	require.NoError(t, err)

	req, err := http.NewRequest("DELETE", "/auth/whitelist", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	ss := httptest.NewRecorder()
	router.ServeHTTP(ss, req)

	require.Equal(t, http.StatusNoContent, ss.Code)

	ivalidMask := cases[1].network
	body, err = json.Marshal(ivalidMask)
	require.NoError(t, err)

	req, err = http.NewRequest("DELETE", "/auth/whitelist", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	ss = httptest.NewRecorder()
	router.ServeHTTP(ss, req)

	require.Equal(t, http.StatusBadRequest, ss.Code)

}

func TestWhiteList_GetIPs(t *testing.T) {

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockStor := mock_service.NewMockWhiteListStore(controller)
	cases := []entity.IpNetwork{{
		Ip:   "127.0.0.1",
		Mask: "255.255.0.0",
	}, {
		Ip:   "127.9.0.1",
		Mask: "255.255.0.0",
	}}

	mockStor.EXPECT().GetIPs().Return(cases, nil).AnyTimes()

	whiteListService := service.NewWhiteList(mockStor, &logger)
	whiteList := NewWhiteList(whiteListService, &logger)

	router := mux.NewRouter()
	router.HandleFunc("/auth/whitelist", whiteList.GetIPs).Methods("GET")

	req, err := http.NewRequest("GET", "/auth/whitelist", bytes.NewReader(nil))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	ss := httptest.NewRecorder()
	router.ServeHTTP(ss, req)

	require.Equal(t, http.StatusOK, ss.Code)

	var ipList []entity.IpNetwork
	err = json.Unmarshal(ss.Body.Bytes(), &ipList)
	require.NoError(t, err)
	require.Equal(t, cases, ipList)
}
