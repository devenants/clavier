# clavier

[![Go Report Card](https://goreportcard.com/badge/github.com/devenants/clavier)](https://goreportcard.com/report/github.com/devenants/clavier)
[![Go Reference](https://pkg.go.dev/badge/github.com/devenants/clavier.svg)](https://pkg.go.dev/github.com/devenants/clavier)
[![Go Build](https://github.com/devenants/clavier/actions/workflows/ci.yml/badge.svg)](https://pkg.go.dev/github.com/devenants/clavier)
![license](https://img.shields.io/badge/license-Apache--2.0-green.svg)

A highly available service discovery framework applied to clients with load balancing and health check functions.

# Overview
In order to enhance the availability of service access, Clavier provides a scalable capability framework from multiple perspectives.

The capabilities provided by Clavier are as follows:
* Extensible Service Discovery Framework
* An easily extensible health check framework
* A highly flexible client-side active load balancing framework

Based on the frameworks, Clavier integrates some specific implementations, including dns-based service discovery, IP list-based polling scheduling strategy, and tcp/http/grpc-based health check process.

The health check strategy is as follows:
* TCP: Trying to establish a connection can only confirm that the peer service has responded to the three-way handshake, which does not represent the true availability of the service.
* HTTP: Use a short link to execute an HTTP/1.1 request, and judge that the return code not less than 500.
* GRPC: Execute the health check request of grpc, and the check return code is SERVING status.

Based on the provided cases, some more scenario-based functions can be easily realized.

# Example
A simple dns-based service discovery case that enables tcp health check.
```
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
			ScoutConfig: &scout.HelperConfig{
				Model: "tcp",
				Data: &tcp.TcpCheckerConfig{
					ConnectTimeout: 1000,
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

	fmt.Printf("%s current\n", time.Now())

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
	time.Sleep(3 * time.Second)

	//get endpoint
	endpoint, err := l.GetEndpoint()
	if err != nil {
		fmt.Printf("first get endpoint failed %v %v\n", l, err)
		return
	}
	fmt.Printf("first endpoint = %v\n", endpoint)

	//get endpoint again
	endpoint, err = l.GetEndpoint()
	if err != nil {
		fmt.Printf("second get endpoint failed %v %v\n", l, err)
		return
	}
	fmt.Printf("second endpoint = %v\n", endpoint)
}
```

## Contributing
- Fork it
- Create your feature branch (`git checkout -b my-new-feature`)
- Commit your changes (`git commit -am 'Add some feature'`)
- Push to the branch (`git push origin my-new-feature`)
- Create new Pull Request
