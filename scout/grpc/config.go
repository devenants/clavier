package grpc

import (
	"google.golang.org/grpc"
)

type GrpcCheckerConfig struct {
	Service      string
	DialOptions  []grpc.DialOption
	CheckTimeout int
}
