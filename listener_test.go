package clavier

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/devenants/clavier/discovery"
	"github.com/devenants/clavier/discovery/dns"
	_ "github.com/devenants/clavier/discovery/dns"
	"github.com/devenants/clavier/filter"
	_ "github.com/devenants/clavier/filter/default"
	_ "github.com/devenants/clavier/filter/round_robin"
	"github.com/devenants/clavier/scout"
)

func TestListener(t *testing.T) {
	conf := &ListenerConfig{
		Group: "localhost",
		Discovery: &DiscoveryConfig{
			Model: "dns",
			Config: discovery.ModelConfig{
				Data: &dns.DnsConfig{
					Port: "80",
				},
			},
		},
		Entry: &EntryManagerConfig{
			FilterConfig: &FilterMangerConfig{
				Model:      "round-robin",
				StatusJump: true,
				Config: filter.ModelConfig{
					Data: nil,
				},
			},
			ScoutConfig: &scout.HelperConfig{
				Model: "none",
			},
		},
	}

	ctx := context.Background()

	lis, err := NewListener(ctx, conf, nil)
	require.Equal(t, err, nil, "")
	require.NotEqual(t, lis, nil, "")

	time.Sleep(4 * time.Second)

	dst := lis.ListEndpoints()
	require.NotEqual(t, len(dst), 0, "")

	lis.GetEndpoint()

	ctx.Done()
}
