package service

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/config"
	"github.com/zorg113/golang_dipl/atibruteforce/model/entity"
	mock_service "github.com/zorg113/golang_dipl/atibruteforce/store/adapters/mocks"
)

func Test_Bucket(t *testing.T) {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	controller := gomock.NewController(t)
	defer controller.Finish()
	blackListMockStor := mock_service.NewMockBlackListStore(controller)
	blacklist := NewBlackList(blackListMockStor, &logger)

	whiteListMockStor := mock_service.NewMockWhiteListStore(controller)
	whitelist := NewWhiteList(whiteListMockStor, &logger)
	cfg, err := config.NewConfig("./../../config/conf.yaml")
	require.NoError(t, err)
	auth := NewAuthorization(blacklist, whitelist, &cfg, &logger)
	req := entity.Request{
		Login:    "test",
		Password: "123",
		IP:       "127.0.0.1",
	}
	blackListMockStor.EXPECT().GetIPs().Return([]entity.IPNetwork{}, nil).AnyTimes()
	whiteListMockStor.EXPECT().GetIPs().Return([]entity.IPNetwork{}, nil).AnyTimes()
	for i := 0; i < 10; i++ {
		res, err := auth.Authorization(req)
		require.NoError(t, err)
		require.True(t, res)
	}
	res, err := auth.Authorization(req)
	require.NoError(t, err)
	require.False(t, res)
	res = auth.ResetIPBucket("127.0.0.1")
	require.True(t, res)
	res = auth.ResetLoginInBucket("test")
	require.True(t, res)
	res = auth.ResetPasswordInBucket("test")
	require.False(t, res)
	res, err = auth.Authorization(req)
	require.NoError(t, err)
	require.True(t, res)
}
