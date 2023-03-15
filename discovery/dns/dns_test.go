package dns

import (
	"testing"

	"github.com/devenants/clavier/discovery"
	"github.com/stretchr/testify/require"
)

func TestResolver(t *testing.T) {
	r, err := NewDnsResolver(&discovery.ModelConfig{
		Data: &DnsConfig{
			Port: "80",
		},
	})

	require.Equal(t, err, nil, "")
	require.NotEqual(t, r, nil, "")

	m := r.Model()
	require.Equal(t, m, "dns", "")

	dst, err := r.Lookup("www.baidu.com", nil)
	require.Equal(t, err, nil, "")
	require.NotEqual(t, len(dst), 1, "")
}
