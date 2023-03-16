package clavier

import (
	"context"
	"testing"
	"time"

	"github.com/devenants/clavier/discovery"
	"github.com/devenants/clavier/discovery/dns"
	"github.com/devenants/clavier/filter"
	_ "github.com/devenants/clavier/filter/default"
	_ "github.com/devenants/clavier/filter/round_robin"
	"github.com/devenants/clavier/scout"
	"github.com/devenants/clavier/scout/custom"
	sgrpc "github.com/devenants/clavier/scout/grpc"
	"github.com/devenants/clavier/scout/http"
	"github.com/devenants/clavier/scout/tcp"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

var (
	c *Clavier
)

func TestNewClavier(t *testing.T) {
	ctx := context.Background()
	c, err = NewClavier(ctx)
	require.Equal(t, err, nil, "")
	require.NotEqual(t, c, nil, "")

	t.Run("testAvailableTest", testListenerTest)
	t.Run("testEndpointTest", testEndpointTest)
	t.Run("testDelListenerTest", testDelListenerTest)
	t.Run("testScoutCustom", testScoutCustom)
	t.Run("testScoutHttp", testScoutHttp)
	t.Run("testScoutGrpc", testScoutGrpc)
	t.Run("testScoutTcp", testScoutTcp)
}

func testListenerTest(t *testing.T) {
	conf := &ListenerConfig{
		Group: "www.baidu.com",
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

	l, err := c.AddListener("baidu", conf)
	require.Equal(t, err, nil, "")
	require.NotEqual(t, l, nil, "")

	time.Sleep(3 * time.Second)
}

func testEndpointTest(t *testing.T) {
	time.Sleep(3 * time.Second)
	al, err := c.ListEndpoints("baidu")
	require.Equal(t, err, nil, "")
	require.NotEqual(t, len(al), 0, "")

	time.Sleep(4 * time.Second)

	addr, err := c.GetEndpoint("baidu")
	require.Equal(t, err, nil, "")
	require.NotEqual(t, addr, nil, "")
}

func testDelListenerTest(t *testing.T) {
	count := c.ListenerCount()
	require.Equal(t, count, 1, "")

	conf := &ListenerConfig{
		Group: "www.taobao.com",
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

	l, err := c.AddListener("taobao", conf)
	require.Equal(t, err, nil, "")
	require.NotEqual(t, l, nil, "")

	time.Sleep(3 * time.Second)

	count = c.ListenerCount()
	require.Equal(t, count, 2, "")

	c.DelListener("taobao")
	count = c.ListenerCount()
	require.Equal(t, count, 1, "")
}

func testScoutCustom(t *testing.T) {
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
				Model: "custom",
				Data: &custom.CustomCheckerConfig{
					Probe: func(_ interface{}) (interface{}, error) {
						return true, nil
					},
				},
			},
		},
	}

	l, err := c.AddListener("baidu", conf)
	require.Equal(t, err, nil, "")
	require.NotEqual(t, l, nil, "")

	time.Sleep(6 * time.Second)

	_, err = c.GetEndpoint("baidu")
	require.Equal(t, err, nil, "")
}

func testScoutHttp(t *testing.T) {
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
				Model: "http",
				Data: &http.HttpCheckerConfig{
					RequestTimeout: 1000,
					URL:            "/",
					Method:         "GET",
				},
			},
		},
	}

	l, err := c.AddListener("baidu", conf)
	require.Equal(t, err, nil, "")
	require.NotEqual(t, l, nil, "")

	time.Sleep(6 * time.Second)

	_, err = c.GetEndpoint("baidu")
	require.Equal(t, err, nil, "")
}

func testScoutGrpc(t *testing.T) {
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
				Model: "grpc",
				Data: &sgrpc.GrpcCheckerConfig{
					Service: "HealthTest",
					DialOptions: []grpc.DialOption{
						grpc.WithInsecure(),
					},
				},
			},
		},
	}

	l, err := c.AddListener("baidu", conf)
	require.Equal(t, err, nil, "")
	require.NotEqual(t, l, nil, "")

	time.Sleep(6 * time.Second)

	_, err = c.GetEndpoint("baidu")
	require.Equal(t, err, nil, "")
}

func testScoutTcp(t *testing.T) {
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
				Model: "tcp",
				Data: &tcp.TcpCheckerConfig{
					ConnectTimeout: 1000,
				},
			},
		},
	}

	l, err := c.AddListener("baidu", conf)
	require.Equal(t, err, nil, "")
	require.NotEqual(t, l, nil, "")

	time.Sleep(6 * time.Second)

	_, err = c.GetEndpoint("baidu")
	require.Equal(t, err, nil, "")
}
