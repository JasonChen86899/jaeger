package grpc

import (
	"context"
	"strings"
	"sync"

	"go.uber.org/zap"

	"github.com/jaegertracing/jaeger/cmd/collector/app/processor"
	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/proto-gen/api_v2"
)

const (
	servicesIPsTagKey = "tag.services.ips"

	traceErrorTagKey         = "tag.tail_based_sampling.error"
	traceParentErrorTagValue = 1
	traceSelfErrorTagValue   = 0

	defaultRequestChannelSize = 1024 * 10
)

// GRPCHandler get span with traceID from agent.
type TailBasedSamplingHandler struct {
	logger      *zap.Logger
	connBuilder *ConnBuilder

	once        sync.Once
	requestChan chan *model.Span

	spanProcessor processor.SpanProcessor
}

func NewTailBasedSamplingHandler(logger *zap.Logger, builder *ConnBuilder, processor processor.SpanProcessor) *TailBasedSamplingHandler {
	return &TailBasedSamplingHandler{
		logger:      logger,
		connBuilder: builder,

		once:        sync.Once{},
		requestChan: make(chan *model.Span, defaultRequestChannelSize),

		spanProcessor: processor,
	}
}

func (h *TailBasedSamplingHandler) FilterSpans(spans []*model.Span) {
	h.once.Do(func() {
		go h.daemonProcessRequest()
	})

	// filter span for next process
	for _, span := range spans {
		for _, kv := range span.Tags {
			if kv.Key == traceErrorTagKey {
				if kv.VInt64 == traceSelfErrorTagValue {
					h.requestChan <- span
					break
				}
			}
		}
	}
}

func (h *TailBasedSamplingHandler) daemonProcessRequest() {
	for {
		select {
		case span := <-h.requestChan:
			for _, kv := range span.Tags {
				if kv.Key == servicesIPsTagKey {
					ips := kv.VStr
					ipList := strings.Split(ips, ",")
					// remove it self ip
					ipList = ipList[:len(ipList)-1]

					// do grpc request
					spans := make([]*model.Span, 0, len(ipList))
					w := sync.WaitGroup{}
					w.Add(len(ipList))
					for _, ip := range ipList {
						agentIP := ip
						go func() {
							agentSpan := h.grpcRequest(agentIP, span.TraceID)
							spans = append(spans, agentSpan)
							w.Done()
						}()
					}
					w.Wait()

					// batch process
					_, err := h.spanProcessor.ProcessSpans(spans, processor.SpansOptions{
						InboundTransport: processor.GRPCTransport,
						SpanFormat:       processor.ProtoSpanFormat,
					})
					if err != nil {
						h.logger.Error("tail based sampling handler process spans", zap.Error(err))
					}
				}
			}
		}
	}
}

func (h *TailBasedSamplingHandler) grpcRequest(agentIP string, traceID model.TraceID) *model.Span {
	gConn, err := h.connBuilder.GetConnection(agentIP)
	if err != nil {
		h.logger.Error("Collector with tail based sampling build grpc connection to agent failed",
			zap.String("agentIP", agentIP),
			zap.String("traceID", traceID.String()),
			zap.Error(err))
		return nil
	}

	res, err := api_v2.NewQueryServiceClient(gConn).GetTrace(context.Background(), &api_v2.GetTraceRequest{
		TraceID: traceID,
	})

	if err != nil {
		h.logger.Error("Collector with tail based sampling send grpc request to agent failed",
			zap.String("agentIP", agentIP),
			zap.String("traceID", traceID.String()),
			zap.Error(err))
		return nil
	}

	spansChunk, err := res.Recv()
	if err != nil {
		h.logger.Error("Collector with tail based sampling get grpc response from agent failed",
			zap.String("agentIP", agentIP),
			zap.String("traceID", traceID.String()),
			zap.Error(err))
		return nil
	}

	return &spansChunk.Spans[0]
}
