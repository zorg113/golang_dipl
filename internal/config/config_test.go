package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpenConfig(t *testing.T) {
	var conf Config
	conf, err := NewConfig("./../../config/conf.yaml")
	require.NoError(t, err)
	require.Equal(t, "postgres", conf.DBData.DBPassword)
}
