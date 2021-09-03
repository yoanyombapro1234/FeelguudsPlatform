package authentication_handler

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	core_metrics "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-metrics"
)

type Telemetry struct {
	MicroServiceMetrics *ServiceMetrics
}

// NewTelemetry initializes a new instance of the telemetry object
func NewTelemetry(engine *core_metrics.CoreMetricsEngine, serviceName string) *Telemetry {
	return &Telemetry{
		MicroServiceMetrics: NewServiceMetrics(engine, serviceName),
	}
}

type ServiceMetrics struct {
	ServiceName string
	// tracks the number of grpc requests partitioned by name and status code
	// used for monitoring and alerting (RED method)
	GRPCRequestCounter *core_metrics.CounterVec
	// tracks the latencies associated with a GRPC requests by operation name
	// used for horizontal pod auto-scaling (Kubernetes HPA v2)
	GRPCRequestLatencyCounter *core_metrics.HistogramVec
	// tracks the number of times there was a failure or success when trying to extract id from the request url
	ExtractIdOperationCounter *core_metrics.CounterVec
	// tracks the number of times there was a failure or success when trying to extract id from the request url
	RemoteOperationStatusCounter *core_metrics.CounterVec
	// tracks the latency of various remote operations
	RemoteOperationsLatencyCounter *core_metrics.HistogramVec
	// tracks the number of invalid requests processed by the service
	InvalidRequestParametersCounter *core_metrics.CounterVec
	// tracks the number of failed casting operations captured by the service
	CastingOperationFailureCounter *core_metrics.CounterVec
	// tracks the number of failed request decoding operations for the service
	DecodeRequestStatusCounter *core_metrics.CounterVec
}

// NewServiceMetrics returns a pointer reference to a metrics objects encapsulating all registered counters for this service
func NewServiceMetrics(engine *core_metrics.CoreMetricsEngine, serviceName string) *ServiceMetrics {
	return &ServiceMetrics{
		ServiceName:                     serviceName,
		GRPCRequestCounter:              NewGRPCRequestCounter(engine, serviceName),
		GRPCRequestLatencyCounter:       NewGRPCRequestLatencyCounter(engine, serviceName),
		ExtractIdOperationCounter:       NewExtractIdOperationCounter(engine, serviceName),
		RemoteOperationStatusCounter:    NewRemoteOperationStatusCounter(engine, serviceName),
		RemoteOperationsLatencyCounter:  NewRemoteOperationLatencyCounter(engine, serviceName),
		InvalidRequestParametersCounter: NewInvalidRequestParametersCounter(engine, serviceName),
		CastingOperationFailureCounter:  NewCastingOperationFailureCounter(engine, serviceName),
		DecodeRequestStatusCounter:      NewDecodeRequestStatusCounter(engine, serviceName),
	}
}

// NewGRPCRequestCounter returns a counter instance capturing the number of grpd requests
func NewGRPCRequestCounter(engine *core_metrics.CoreMetricsEngine, serviceName string) *core_metrics.CounterVec {
	newCounter := core_metrics.NewCounterVec(&core_metrics.CounterOpts{
		Namespace: serviceName,
		Subsystem: "GRPC",
		Name:      fmt.Sprintf("%s_http_requests_total", serviceName),
		Help:      "How many GRPC requests processed partitioned by name and status code",
	}, []string{"name", "code"})

	engine.RegisterMetric(newCounter)
	return newCounter
}

// NewGRPCRequestLatencyCounter returns a counter instance capturing the request latency of a grpc operation
func NewGRPCRequestLatencyCounter(engine *core_metrics.CoreMetricsEngine, serviceName string) *core_metrics.HistogramVec {
	newCounter := core_metrics.NewHistogramVec(&core_metrics.HistogramOpts{
		Namespace:         serviceName,
		Subsystem:         "GRPC",
		Name:              fmt.Sprintf("%s_http_requests_latencies", serviceName),
		Help:              "Seconds spent serving GRPC requests.",
		ConstLabels:       nil,
		Buckets:           prometheus.DefBuckets,
		DeprecatedVersion: "",
		StabilityLevel:    "",
	}, []string{"method", "path", "status"})
	engine.RegisterMetric(newCounter)
	return newCounter
}

// NewExtractIdOperationCounter returns an instance of the status of the extract id operation counter
func NewExtractIdOperationCounter(engine *core_metrics.CoreMetricsEngine, serviceName string) *core_metrics.CounterVec {
	// tracks the number of times there was a failure or success when trying to extract id from the request url
	newCounter := core_metrics.NewCounterVec(&core_metrics.CounterOpts{
		Namespace: serviceName,
		Subsystem: "HTTP",
		Name:      fmt.Sprintf("%s_status_of_extract_id_operation_from_requests_total", serviceName),
		Help:      "The status of the extract the id operation from the HTTP requests processed partitioned by operation name and operation status",
	}, []string{"operation_name", "status"})
	engine.RegisterMetric(newCounter)
	return newCounter
}

// NewRemoteOperationStatusCounter returns an instance of a counter capturing the status of an rpc operation
func NewRemoteOperationStatusCounter(engine *core_metrics.CoreMetricsEngine, serviceName string) *core_metrics.CounterVec {
	newCounter := core_metrics.NewCounterVec(&core_metrics.CounterOpts{
		Namespace: serviceName,
		Subsystem: "HTTP",
		Name:      fmt.Sprintf("%s_status_of_remote_operation_total", serviceName),
		Help:      "A count of the status all remote operations operation",
	}, []string{"operation_name", "status"})
	engine.RegisterMetric(newCounter)
	return newCounter
}

// NewRemoteOperationLatencyCounter returns an instance of the rpc operation latency counter
func NewRemoteOperationLatencyCounter(engine *core_metrics.CoreMetricsEngine, serviceName string) *core_metrics.HistogramVec {
	newCounter := core_metrics.NewHistogramVec(&core_metrics.HistogramOpts{
		Namespace:         serviceName,
		Subsystem:         "HTTP",
		Name:              fmt.Sprintf("%s_remote_operation_requests_latencies", serviceName),
		Help:              "Seconds spent serving remote operations HTTP requests.",
		ConstLabels:       nil,
		Buckets:           prometheus.DefBuckets,
		DeprecatedVersion: "",
		StabilityLevel:    "",
	}, []string{"operation", "status"})
	engine.RegisterMetric(newCounter)
	return newCounter
}

// NewInvalidRequestParametersCounter returns an instance of the invalid request parameters counter
func NewInvalidRequestParametersCounter(engine *core_metrics.CoreMetricsEngine, serviceName string) *core_metrics.CounterVec {
	newCounter := core_metrics.NewCounterVec(&core_metrics.CounterOpts{
		Namespace: serviceName,
		Subsystem: "HTTP",
		Name:      fmt.Sprintf("%s_invalid_request_parameters_total", serviceName),
		Help:      "A count of the total number of invalid request parameter count",
	}, []string{"operation_name"})
	engine.RegisterMetric(newCounter)
	return newCounter
}

// NewCastingOperationFailureCounter returns an instance of the casting operation failure counter
func NewCastingOperationFailureCounter(engine *core_metrics.CoreMetricsEngine, serviceName string) *core_metrics.CounterVec {
	newCounter := core_metrics.NewCounterVec(&core_metrics.CounterOpts{
		Namespace: serviceName,
		Subsystem: "HTTP",
		Name:      fmt.Sprintf("%s_casting_operation_failure_total", serviceName),
		Help:      "A count of the total number of failed casts from interface to object",
	}, []string{"operation_name"})
	engine.RegisterMetric(newCounter)
	return newCounter
}

// NewDecodeRequestStatusCounter returns an instance of the request status counter
func NewDecodeRequestStatusCounter(engine *core_metrics.CoreMetricsEngine, serviceName string) *core_metrics.CounterVec {
	newCounter := core_metrics.NewCounterVec(&core_metrics.CounterOpts{
		Namespace: serviceName,
		Subsystem: "HTTP",
		Name:      fmt.Sprintf("%s_decoder_request_op_counter_total", serviceName),
		Help:      "A count of the status of all decode operations",
	}, []string{"operation_name", "status"})
	engine.RegisterMetric(newCounter)
	return newCounter
}
