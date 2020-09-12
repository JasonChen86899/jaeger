package grpc

import (
	"context"
	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/proto-gen/api_v2"
	"go.uber.org/zap"
	"strings"
	"sync"
)

const (
	servicesIPsTagKey = "tag.services.ips"

	traceErrorTagKey         = "tag.tail_based_sampling.error"
	traceParentErrorTagValue = 1
	traceSelfErrorTagValue   = 0

	defaultRequestChannelSize = 1024 * 10
)

// GRPCHandler get span with traceID from agent.
type Handler struct {
	logger      *zap.Logger
	connBuilder *ConnBuilder

	once        sync.Once
	requestChan chan *model.Span
}

func NewTailBasedSamplingHandler() *Handler {
	return &Handler{
		logger:      nil,
		connBuilder: nil,

		once:        sync.Once{},
		requestChan: make(chan *model.Span, defaultRequestChannelSize),
	}
}

func (h *Handler) filterSpans(spans []*model.Span) {
	h.once.Do(func() {
		go h.processRequest()
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

func (h *Handler) processRequest() {
	for {
		select {
		case span := <-h.requestChan:
			for _, kv := range span.Tags {
				if kv.Key == servicesIPsTagKey {
					ips := kv.VStr
					ipList := strings.Split(ips, ",")
					ipList = ipList[:len(ipList)-1]

					// do grpc request
					for _, ip := range ipList {
						go func() {
							agentSpan := h.grpcRequest(ip, span.TraceID)

						}()
					}
				}
			}
		}
	}
}

func (h *Handler) grpcRequest(agentIP string, traceID model.TraceID) *model.Span {
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
