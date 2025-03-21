package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommon(t *testing.T) {
	ip := "127.88.99.22"
	mask := "255.255.0.0"
	rez, err := GetPrefix(ip, mask)
	require.NoError(t, err)
	require.Equal(t, "127.88.0.0", rez)
}
