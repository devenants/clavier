package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInstanceCreate(t *testing.T) {
	a := &Endpoint{
		Host: "192.168.11.2",
	}

	require.Equal(t, a.ToString(), a.Host, "")

	b := &Endpoint{
		Host: "192.168.11.3",
		Port: "8080",
	}
	require.Equal(t, b.ToString(), "192.168.11.3:8080", "")
}
