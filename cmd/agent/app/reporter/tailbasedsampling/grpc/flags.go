package grpc

import (
	"flag"
	"fmt"
	"github.com/jaegertracing/jaeger/pkg/config/tlscfg"
	"github.com/jaegertracing/jaeger/ports"
	"github.com/spf13/viper"
)

const (
	reportTailBasedSamplingOpen   = "reporter.tail-based-sampling.open"
)

var tlsFlagsConfig = tlscfg.ServerFlagsConfig{
	Prefix:       "reporter.tail-based-sampling.grpc",
	ShowEnabled:  true,
	ShowClientCA: true,
}

type TailBasedSamplingOptions struct {
	Open bool
	// GRPCHostPort is the host:port address that the reporter service listens on for gRPC requests
	GRPCHostPort string
	// TLS configures secure transport
	TLS tlscfg.Options
}

// AddFlags adds flags for reporter TailBasedSamplingOptions
func AddFlags(flags *flag.FlagSet) {
	flags.Bool(reportTailBasedSamplingOpen, true, fmt.Sprintf("Set the reporter tail-based-sampling server open or close"))
	AddGrpcFlags(flags)
}

// AddOTELJaegerFlags adds flags that with grpc server used by reporter TailBasedSamplingOptions
func AddGrpcFlags(flags *flag.FlagSet) {
	tlsFlagsConfig.AddFlags(flags)
}

// InitFromViper initializes CollectorOptions with properties from viper
func (cOpts *TailBasedSamplingOptions) InitFromViper(v *viper.Viper) *TailBasedSamplingOptions {
	cOpts.Open = v.GetBool(reportTailBasedSamplingOpen)
	cOpts.GRPCHostPort = ports.PortToHostPort(ports.ReportTailBasedSamplingGRPC)
	cOpts.TLS = tlsFlagsConfig.InitFromViper(v)
	return cOpts
}