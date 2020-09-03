// Copyright (c) 2018 The Jaeger Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package reporter

import (
	"flag"
	"fmt"
	tail_based_sampling "github.com/jaegertracing/jaeger/cmd/agent/app/reporter/tailbasedsampling/grpc"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/jaegertracing/jaeger/cmd/all-in-one/setupcontext"
	"github.com/jaegertracing/jaeger/cmd/flags"
)

const (
	// Reporter type
	reporterType = "reporter.type"
	// AgentTagsDeprecated is a configuration property name for adding process tags to incoming spans.
	AgentTagsDeprecated = "jaeger.tags"
	agentTags           = "agent.tags"
	// GRPC is name of gRPC reporter.
	GRPC Type = "grpc"

	// tail-based Sampling
	tailBasedSampling = "tail-based-sampling.open"
)

// Type defines type of reporter.
type Type string

// Options holds generic reporter configuration.
type Options struct {
	ReporterType Type
	AgentTags    map[string]string

	TailBasedSampling *tail_based_sampling.TailBasedSamplingOptions
}

// AddFlags adds flags for Options.
func AddFlags(flags *flag.FlagSet) {
	flags.String(reporterType, string(GRPC), fmt.Sprintf("Reporter type to use e.g. %s", string(GRPC)))
	if !setupcontext.IsAllInOne() {
		flags.String(AgentTagsDeprecated, "", "(deprecated) see --"+agentTags)
		flags.String(agentTags, "", "One or more tags to be added to the Process tags of all spans passing through this agent. Ex: key1=value1,key2=${envVar:defaultValue}")
	}

	tail_based_sampling.AddFlags(flags)
}

// InitFromViper initializes Options with properties retrieved from Viper.
func (b *Options) InitFromViper(v *viper.Viper, logger *zap.Logger) *Options {
	b.ReporterType = Type(v.GetString(reporterType))
	if !setupcontext.IsAllInOne() {
		if len(v.GetString(AgentTagsDeprecated)) > 0 {
			logger.Warn("Using deprecated configuration", zap.String("option", AgentTagsDeprecated))
			b.AgentTags = flags.ParseJaegerTags(v.GetString(AgentTagsDeprecated))
		}
		if len(v.GetString(agentTags)) > 0 {
			b.AgentTags = flags.ParseJaegerTags(v.GetString(agentTags))
		}
	}

	b.TailBasedSampling = new(tail_based_sampling.TailBasedSamplingOptions).InitFromViper(v)
	return b
}
