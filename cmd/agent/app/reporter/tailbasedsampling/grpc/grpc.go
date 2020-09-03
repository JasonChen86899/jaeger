package grpc

import (
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/jaegertracing/jaeger/pkg/config/tlscfg"
)

// GRPCServerParams to construct a new Jaeger Collector gRPC Server
type ServerParams struct {
	TLSConfig     tlscfg.Options
	HostPort      string
	Handler       *Handler
	Logger        *zap.Logger
	OnError       func(error)
}

func StartGRPCServer(params *ServerParams) (*grpc.Server, error) {
	var server *grpc.Server

	if params.TLSConfig.Enabled {
		// user requested a server with TLS, setup creds
		tlsCfg, err := params.TLSConfig.Config(params.Logger)
		if err != nil {
			return nil, err
		}

		creds := credentials.NewTLS(tlsCfg)
		server = grpc.NewServer(grpc.Creds(creds))
	} else {
		// server without TLS
		server = grpc.NewServer()
	}

	listener, err := net.Listen("tcp", params.HostPort)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on gRPC port: %w", err)
	}

	if err := serveGRPC(server, listener, params); err != nil {
		return nil, err
	}

	return server, nil
}

func serveGRPC(server *grpc.Server, listener net.Listener, params *ServerParams) error {
	params.Logger.Info("Starting jaeger-agent tail-based Sampling gRPC server", zap.String("grpc.host-port", params.HostPort))
	go func() {
		if err := server.Serve(listener); err != nil {
			params.Logger.Error("Could not launch gRPC service", zap.Error(err))
			if params.OnError != nil {
				params.OnError(err)
			}
		}
	}()

	return nil
}