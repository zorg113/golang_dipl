package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpenConfig(t *testing.T) {
	var conf Config
	conf, err := NewConfig("config.yaml")
	require.NoError(t, err)
	require.Equal(t, "HiWorld", conf.DbData.DbPassword)

}
