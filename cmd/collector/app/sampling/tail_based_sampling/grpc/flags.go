package grpc

import (
	"flag"
	"fmt"

	"github.com/jaegertracing/jaeger/pkg/config/tlscfg"
	"github.com/spf13/viper"
)

const (
	collectorTailBasedSamplingOpen = "collector.tail-based-sampling.open"

	gRPCPrefix      = "collector.tail-based-sampling.grpc"
	retry           = gRPCPrefix + ".retry.max"
	defaultMaxRetry = 3
)

var tlsFlagsConfig = tlscfg.ServerFlagsConfig{
	Prefix:       "collector.tail-based-sampling.grpc",
	ShowEnabled:  true,
	ShowClientCA: true,
}

type TailBasedSamplingOptions struct {
	Open     bool
	MaxRetry uint
	TLS      tlscfg.Options
}

func NewTailBasedSamplingOptions() *TailBasedSamplingOptions {
	return &TailBasedSamplingOptions{}
}

// AddFlags adds flags for Collector TailBasedSamplingOptions
func AddFlags(flags *flag.FlagSet) {
	flags.Uint(retry, defaultMaxRetry, "Sets the maximum number of retries for a call")
	flags.Bool(collectorTailBasedSamplingOpen, true, fmt.Sprintf("Set the collector tail-based-sampling server open or close"))
	tlsFlagsConfig.AddFlags(flags)
}

// InitFromViper initializes Collector TailBasedSamplingOptions with properties from viper
func (cOpts *TailBasedSamplingOptions) InitFromViper(v *viper.Viper) *TailBasedSamplingOptions {
	cOpts.Open = v.GetBool(collectorTailBasedSamplingOpen)
	cOpts.MaxRetry = uint(v.GetInt(retry))
	cOpts.TLS = tlsFlagsConfig.InitFromViper(v)
	return cOpts
}
