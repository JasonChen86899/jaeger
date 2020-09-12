package grpc

const (
	TraceErrorTagKey             = "tag.tail_based_sampling.error"
	TraceParentErrorTagValue     = 1
	TraceSelfErrorTagValue       = 0

	TraceErrorBaggageKey         = "baggage.tail_based_sampling.error"
	TraceParentErrorBaggageValue = "1"
	TraceSelfErrorBaggageValue   = "0"
)
