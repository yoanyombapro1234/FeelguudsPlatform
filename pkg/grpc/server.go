package grpc

import (
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	logger *zap.Logger
	config *Config
}

type Config struct {
	Port                        int    `mapstructure:"GRPC_PORT"`
	ServiceName                 string `mapstructure:"GRPC_SERVICE_NAME"`
	RpcDeadline                 int    `mapstructure:"GRPC_RPC_DEADLINE_IN_MS"`
	RpcRetries                  int    `mapstructure:"GRPC_RPC_RETRIES"`
	RpcRetryTimeout             int    `mapstructure:"GRPC_RPC_RETRY_TIMEOUT_IN_MS"`
	RpcRetryBackoff             int    `mapstructure:"GRPC_RPC_RETRY_BACKOFF_IN_MS"`
	EnableTls                   bool   `mapstructure:"GRPC_ENABLE_TLS"`
	CertificatePath             string `mapstructure:"GRPC_CERT_PATH"`
	EnableDelayMiddleware       bool   `mapstructure:"ENABLE_RANDOM_DELAY"`
	EnableRandomErrorMiddleware bool   `mapstructure:"ENABLE_RANDOM_RANDOM_ERROR"`
	MinRandomDelay              int    `mapstructure:"RANDOM_DELAY_MIN_IN_MS"`
	MaxRandomDelay              int    `mapstructure:"RANDOM_DELAY_MAX_IN_MS"`
	DelayUnit                   string `mapstructure:"RANDOM_DELAY_UNIT"`
	Version                     string `mapstructure:"VERSION"`
	MetricAddr                  string `mapstructure:"METRIC_CONNECTION_ADDRESS"`
}

func NewServer(config *Config, logger *zap.Logger) (*Server, error) {
	srv := &Server{
		logger: logger,
		config: config,
	}

	return srv, nil
}

func (s *Server) ListenAndServe() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", s.config.Port))
	if err != nil {
		s.logger.Fatal("failed to listen", zap.Int("port", s.config.Port))
	}

	srv := grpc.NewServer()
	server := health.NewServer()
	reflection.Register(srv)
	grpc_health_v1.RegisterHealthServer(srv, server)
	server.SetServingStatus(s.config.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)

	if err := srv.Serve(listener); err != nil {
		s.logger.Fatal("failed to serve", zap.Error(err))
	}
}
