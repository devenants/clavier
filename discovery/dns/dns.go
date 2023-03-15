package dns

import (
	"fmt"
	"net"

	"github.com/devenants/clavier/discovery"
	"github.com/devenants/clavier/types"
)

const (
	modelName = "dns"
)

type DnsResolver struct {
	config *DnsConfig
}

func NewDnsResolver(config *discovery.ModelConfig) (*DnsResolver, error) {
	var conf *DnsConfig
	var ok bool
	if config.Data != nil {
		conf, ok = config.Data.(*DnsConfig)
		if !ok {
			return nil, fmt.Errorf("dns resolver model data invalid %v", config)
		}
	}

	return &DnsResolver{
		config: conf,
	}, nil
}

func (r *DnsResolver) Model() string {
	return modelName
}

func (r *DnsResolver) Lookup(name string, option interface{}) ([]*types.Endpoint, error) {
	dst, err := net.LookupHost(name)
	if err != nil {
		return nil, err
	}

	endpoints := make([]*types.Endpoint, 0)
	for _, item := range dst {
		end := &types.Endpoint{
			Host: item,
			Port: r.config.Port,
		}
		endpoints = append(endpoints, end)
	}

	return endpoints, nil
}

func init() {
	discovery.Register(modelName, func(conf *discovery.ModelConfig) (discovery.DiscoveryModel, error) {
		return NewDnsResolver(conf)
	})
}
