package grpc

import (
	"context"
	"github.com/jaegertracing/jaeger/cmd/agent/app/reporter/tailbasedsampling/buffer"
	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/proto-gen/api_v2"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	maxSpanCountInChunk = 10

	msgTraceNotFound = "trace not found"
)

// GRPCHandler implements the gRPC endpoint of the query service.
type Handler struct {
	buffer *buffer.WindowBuffer
	logger *zap.Logger
}

// NewGRPCHandler returns a GRPCHandler
func NewGRPCHandler(buffer *buffer.WindowBuffer, logger *zap.Logger) *Handler {
	gH := &Handler{
		buffer: buffer,
		logger: logger,
	}

	return gH
}

func (g *Handler) sendSpanChunks(spans []*model.Span, sendFn func(*api_v2.SpansResponseChunk) error) error {
	chunk := make([]model.Span, 0, len(spans))
	for i := 0; i < len(spans); i += maxSpanCountInChunk {
		chunk = chunk[:0]
		for j := i; j < len(spans) && j < i+maxSpanCountInChunk; j++ {
			chunk = append(chunk, *spans[j])
		}
		if err := sendFn(&api_v2.SpansResponseChunk{Spans: chunk}); err != nil {
			g.logger.Error("failed to send response to client", zap.Error(err))
			return err
		}
	}
	return nil
}

// GetTrace is the gRPC handler to fetch traces based on trace-id.
func (g *Handler) GetTrace(r *api_v2.GetTraceRequest, stream api_v2.QueryService_GetTraceServer) error {
	spans, ok := g.buffer.Get(r.TraceID.String())
	if !ok {
		g.logger.Error(msgTraceNotFound)
		return status.Errorf(codes.NotFound, "%s", msgTraceNotFound)
	}

	return g.sendSpanChunks(spans, stream.Send)
}

// ArchiveTrace is the gRPC handler to archive traces.
func (g *Handler) ArchiveTrace(ctx context.Context, r *api_v2.ArchiveTraceRequest) (*api_v2.ArchiveTraceResponse, error) {
	return nil, nil
}

// FindTraces is the gRPC handler to fetch traces based on TraceQueryParameters.
func (g *Handler) FindTraces(r *api_v2.FindTracesRequest, stream api_v2.QueryService_FindTracesServer) error {
	return nil
}

// GetServices is the gRPC handler to fetch services.
func (g *Handler) GetServices(ctx context.Context, r *api_v2.GetServicesRequest) (*api_v2.GetServicesResponse, error) {
	return nil, nil
}

// GetOperations is the gRPC handler to fetch operations.
func (g *Handler) GetOperations(
	ctx context.Context,
	r *api_v2.GetOperationsRequest,
) (*api_v2.GetOperationsResponse, error) {
	return nil, nil
}

// GetDependencies is the gRPC handler to fetch dependencies.
func (g *Handler) GetDependencies(ctx context.Context, r *api_v2.GetDependenciesRequest) (*api_v2.GetDependenciesResponse, error) {
	return nil, nil
}
