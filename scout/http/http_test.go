package http

import (
	"testing"

	"github.com/devenants/clavier/scout"
	"github.com/stretchr/testify/require"
)

func TestHttpWatcher(t *testing.T) {
	sc := &scout.ModelConfig{
		Data: &HttpCheckerConfig{
			RequestTimeout: 1000,
			URL:            "/",
			Method:         "GET",
		},
	}

	c, err := NewHttpChecker(sc)
	require.Equal(t, err, nil, "")

	n := c.Model()
	require.Equal(t, n, modelName, "")

	e := &testEntry{
		name: "192.168.11.2",
		idx:  0,
	}
	err = c.Register(e)
	require.Equal(t, err, nil, "")

	c.UnRegister(e.name)

	c.Done()
}
