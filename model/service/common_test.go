package service_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zorg113/golang_dipl/atibruteforce/model/service"
)

func TestCommon(t *testing.T) {
	ip := "127.88.99.22"
	mask := "255.255.0.0"
	rez, err := service.GetPrefix(ip, mask)
	require.NoError(t, err)
	require.Equal(t, "127.88.0.0", rez)
}
