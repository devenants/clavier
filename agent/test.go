package main

import (
	"context"
	"fmt"
	"time"

	"github.com/devenants/clavier"
	"github.com/devenants/clavier/discovery"
	"github.com/devenants/clavier/discovery/dns"
	"github.com/devenants/clavier/filter"
	"github.com/devenants/clavier/scout"
	"github.com/devenants/clavier/scout/tcp"
)

func main() {
	conf := &clavier.ListenerConfig{
		Group: "www.baidu.com",
		Discovery: &clavier.DiscoveryConfig{
			Model: "dns",
			Config: discovery.ModelConfig{
				Data: &dns.DnsConfig{
					Port: "80",
				},
			},
		},
		Entry: &clavier.EntryManagerConfig{
			FilterConfig: &clavier.FilterMangerConfig{
				Model:  "round-robin",
				Config: filter.ModelConfig{},
			},
			ScoutConfig: &clavier.ScoutMangerConfig{
				Model: "tcp",
				Config: scout.ModelConfig{
					Data: &tcp.TcpCheckerConfig{
						ConnectTimeout: 1000,
					},
				},
			},
		},
	}

	//create a clavier handler
	c, err := clavier.NewClavier(context.Background())
	if err != nil {
		fmt.Printf("create clavier failed %v\n", err)
		return
	}

	//add a listener
	l, err := c.AddListener("www.baidu.com", conf)
	if err != nil {
		fmt.Printf("add listener failed %v %v\n", conf, err)
		return
	}

	//get endpoint will be failed for dns resolve
	_, err = l.GetEndpoint()
	if err == nil {
		fmt.Printf("get endpoint failed %v %v\n", l, err)
		return
	}

	//wait a moment for dns resolve and health check
	time.Sleep(5 * time.Second)

	//get endpoint
	endpoint, err := l.GetEndpoint()
	if err != nil {
		fmt.Printf("get endpoint failed %v %v\n", l, err)
		return
	}
	fmt.Printf("endpoint = %v\n", endpoint)

	//get endpoint again
	endpoint, err = l.GetEndpoint()
	if err != nil {
		fmt.Printf("get endpoint failed %v %v\n", l, err)
		return
	}
	fmt.Printf("endpoint = %v\n", endpoint)
}
