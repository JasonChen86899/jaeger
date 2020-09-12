package grpc

import (
	"flag"
	"fmt"
	"github.com/jaegertracing/jaeger/pkg/config/tlscfg"
	"github.com/jaegertracing/jaeger/ports"
	"github.com/spf13/viper"
)

const (
	ReportTailBasedSamplingOpen   = "report.tail-based-sampling.open"
	// CollectorGRPCHostPort is the flag for collector gRPC port
	ReportTailBasedSamplingGRPCHostPort   = "report.tail-based-sampling.grpc-server.host-port"
)

var tlsFlagsConfig = tlscfg.ServerFlagsConfig{
	Prefix:       "report.tail-based-sampling.grpc",
	ShowEnabled:  true,
	ShowClientCA: true,
}

type TailBasedSamplingOptions struct {
	Open bool
	// CollectorGRPCHostPort is the host:port address that the reporter service listens in on for gRPC requests
	GRPCHostPort string
	// TLS configures secure transport
	TLS tlscfg.Options
}

// AddFlags adds flags for CollectorOptions
func AddFlags(flags *flag.FlagSet) {
	flags.Bool(ReportTailBasedSamplingOpen, true, fmt.Sprintf("Set the reporter tail-based-sampling server open or close"))
	AddOTELJaegerFlags(flags)
}

// AddOTELJaegerFlags adds flags that are exposed by OTEL Jaeger receier
func AddOTELJaegerFlags(flags *flag.FlagSet) {
	flags.String(ReportTailBasedSamplingGRPCHostPort, ports.PortToHostPort(ports.ReportTailBasedSamplingGRPC), "The port (e.g. :9411) of the report tail-based-sampling server")
	tlsFlagsConfig.AddFlags(flags)
}

// InitFromViper initializes CollectorOptions with properties from viper
func (cOpts *TailBasedSamplingOptions) InitFromViper(v *viper.Viper) *TailBasedSamplingOptions {
	cOpts.Open = v.GetBool(ReportTailBasedSamplingOpen)
	cOpts.GRPCHostPort = v.GetString(ReportTailBasedSamplingGRPCHostPort)
	cOpts.TLS = tlsFlagsConfig.InitFromViper(v)
	return cOpts
}