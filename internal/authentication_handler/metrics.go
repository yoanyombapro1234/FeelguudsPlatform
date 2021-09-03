package authentication_handler

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type ServiceMetrics struct {
	ServiceName string
	// tracks the number of grpc requests partitioned by name and status code
	// used for monitoring and alerting (RED method)
	RequestCounter *prometheus.CounterVec
	// tracks the latencies associated with a GRPC requests by operation name
	// used for horizontal pod auto-scaling (Kubernetes HPA v2)
	RequestLatency *prometheus.HistogramVec
	// tracks the number of times there was a failure or success when trying to extract id from the request url
	ExtractIdOperationCounter *prometheus.CounterVec
	// tracks the status of rpc operations
	RemoteOperationStatusCounter *prometheus.CounterVec
	// tracks the latency of various remote operations
	RemoteOperationLatencyCounter *prometheus.HistogramVec
	// tracks the number of invalid requests processed by the service
	InvalidRequestParameterCounter *prometheus.CounterVec
	// tracks the number of failed casting operations captured by the service
	CastingOperationFailureCounter *prometheus.CounterVec
	// tracks the number of failed request decoding operations for the service
	DecodeRequestStatusCounter *prometheus.CounterVec
}

// NewServiceMetrics returns a pointer reference to a metrics objects encapsulating all registered counters for this service
func NewServiceMetrics(serviceName string) *ServiceMetrics {
	prometheus.MustRegister(prometheus.NewBuildInfoCollector())

	return &ServiceMetrics{
		ServiceName:                    serviceName,
		RequestCounter:                 NewRequestCounter(serviceName),
		RequestLatency:                 NewRequestLatencyCounter(serviceName),
		ExtractIdOperationCounter:      NewExtractIdOperationCounter(serviceName),
		RemoteOperationStatusCounter:   NewRemoteOperationStatusCounter(serviceName),
		RemoteOperationLatencyCounter:  NewRemoteOperationLatencyCounter(serviceName),
		InvalidRequestParameterCounter: NewInvalidRequestParametersCounter(serviceName),
		CastingOperationFailureCounter: NewCastingOperationFailureCounter(serviceName),
		DecodeRequestStatusCounter:     NewDecodeRequestStatusCounter(serviceName),
	}
}

// NewRequestCounter returns a counter instance capturing the number of requests
func NewRequestCounter(serviceName string) *prometheus.CounterVec {
	var metric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: serviceName,
		Subsystem: "API",
		Name:      fmt.Sprintf("%s_http_requests_total", serviceName),
		Help:      "How many requests processed partitioned by name and status code",
	}, []string{"name", "code"})

	prometheus.MustRegister(metric)
	return metric
}

// NewRequestLatencyCounter returns a counter instance capturing the request latency of a grpc operation
func NewRequestLatencyCounter(serviceName string) *prometheus.HistogramVec {
	var metric = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   serviceName,
		Subsystem:   "API",
		Name:        fmt.Sprintf("%s_http_requests_latencies", serviceName),
		Help:        "Seconds spent serving requests.",
		ConstLabels: nil,
		Buckets:     prometheus.DefBuckets,
	}, []string{"method", "path", "status"})

	prometheus.MustRegister(metric)
	return metric
}

// NewExtractIdOperationCounter returns an instance of the status of the extract id operation counter
func NewExtractIdOperationCounter(serviceName string) *prometheus.CounterVec {
	// tracks the number of times there was a failure or success when trying to extract id from the request url
	var metric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: serviceName,
		Subsystem: "API",
		Name:      fmt.Sprintf("%s_status_of_extract_id_operation_from_requests_total", serviceName),
		Help:      "The status of the extract the id operation from the HTTP requests processed partitioned by operation name and operation status",
	}, []string{"operation_name", "status"})

	prometheus.MustRegister(metric)
	return metric
}

// NewRemoteOperationStatusCounter returns an instance of a counter capturing the status of an rpc operation
func NewRemoteOperationStatusCounter(serviceName string) *prometheus.CounterVec {
	var metric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: serviceName,
		Subsystem: "API",
		Name:      fmt.Sprintf("%s_status_of_remote_operation_total", serviceName),
		Help:      "A count of the status all remote operations operation",
	}, []string{"operation_name", "status"})

	prometheus.MustRegister(metric)
	return metric
}

// NewRemoteOperationLatencyCounter returns an instance of the rpc operation latency counter
func NewRemoteOperationLatencyCounter(serviceName string) *prometheus.HistogramVec {
	var metric = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   serviceName,
		Subsystem:   "API",
		Name:        fmt.Sprintf("%s_remote_operation_requests_latencies", serviceName),
		Help:        "Seconds spent serving remote operations HTTP requests.",
		ConstLabels: nil,
		Buckets:     prometheus.DefBuckets,
	}, []string{"operation_name", "status"})

	prometheus.MustRegister(metric)
	return metric
}

// NewInvalidRequestParametersCounter returns an instance of the invalid request parameters counter
func NewInvalidRequestParametersCounter(serviceName string) *prometheus.CounterVec {
	var metric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: serviceName,
		Subsystem: "API",
		Name:      fmt.Sprintf("%s_invalid_request_parameters_total", serviceName),
		Help:      "A count of the total number of invalid request parameter count",
	}, []string{"operation_name"})

	prometheus.MustRegister(metric)
	return metric
}

// NewCastingOperationFailureCounter returns an instance of the casting operation failure counter
func NewCastingOperationFailureCounter(serviceName string) *prometheus.CounterVec {
	var metric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: serviceName,
		Subsystem: "API",
		Name:      fmt.Sprintf("%s_casting_operation_failure_total", serviceName),
		Help:      "A count of the total number of failed casts from interface to object",
	}, []string{"operation_name"})

	prometheus.MustRegister(metric)
	return metric
}

// NewDecodeRequestStatusCounter returns an instance of the request status counter
func NewDecodeRequestStatusCounter(serviceName string) *prometheus.CounterVec {
	var metric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: serviceName,
		Subsystem: "API",
		Name:      fmt.Sprintf("%s_decoder_request_op_counter_total", serviceName),
		Help:      "A count of the status of all decode operations",
	}, []string{"operation_name", "status"})

	prometheus.MustRegister(metric)
	return metric
}
