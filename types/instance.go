package types

import "fmt"

type Endpoints []Endpoint

type Endpoint struct {
	Host   string
	Port   string
	Status bool
}

func (e *Endpoint) ToString() string {
	if len(e.Port) == 0 {
		return e.Host
	}

	return fmt.Sprintf("%s:%s", e.Host, e.Port)
}
