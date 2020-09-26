package grpc

import (
	"fmt"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/jaegertracing/jaeger/pkg/config/tlscfg"
	"github.com/jaegertracing/jaeger/pkg/discovery/grpcresolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// ConnBuilder Struct to hold configurations
type ConnBuilder struct {
	MaxRetry uint
	TLS      tlscfg.Options
}

// NewConnBuilder creates a new grpc connection builder.
func NewConnBuilder(options *TailBasedSamplingOptions) *ConnBuilder {
	return &ConnBuilder{
		MaxRetry: options.MaxRetry,
		TLS:      options.TLS,
	}
}

// CreateConnection creates the gRPC connection
func (b *ConnBuilder) createConnection(dialTarget string, logger *zap.Logger) (*grpc.ClientConn, error) {
	var dialOptions []grpc.DialOption
	if b.TLS.Enabled { // user requested a secure connection
		logger.Info("Collector requested secure grpc connection to agent")
		tlsConf, err := b.TLS.Config(logger)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS config: %w", err)
		}

		creds := credentials.NewTLS(tlsConf)
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(creds))
	} else { // insecure connection
		logger.Info("Collector requested insecure grpc connection to agent")
		dialOptions = append(dialOptions, grpc.WithInsecure())
	}

	dialOptions = append(dialOptions, grpc.WithDefaultServiceConfig(grpcresolver.GRPCServiceConfig))
	dialOptions = append(dialOptions, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(grpc_retry.WithMax(b.MaxRetry))))
	return grpc.Dial(dialTarget, dialOptions...)
}

func (b *ConnBuilder) GetConnection(agentIP string, logger *zap.Logger) (*grpc.ClientConn, error) {
	// TODO
	return b.createConnection(agentIP, logger)
}
