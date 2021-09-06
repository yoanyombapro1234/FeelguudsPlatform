package api

import (
	"time"

	"github.com/gorilla/mux"
	core_logging "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-logging"
)

func NewMockServer() *Server {
	config := &Config{
		Port:                      "9898",
		HttpServerShutdownTimeout: 5 * time.Second,
		HttpServerTimeout:         30 * time.Second,
		BackendURL:                []string{},
		ConfigPath:                "/config",
		DataPath:                  "/data",
		HttpClientTimeout:         30 * time.Second,
		UIColor:                   "blue",
		UIPath:                    ".ui",
		UIMessage:                 "Greetings",
		Hostname:                  "localhost",
	}

	logInstance := core_logging.New("info")
	defer logInstance.ConfigureLogger()
	log := logInstance.Logger

	return &Server{
		router: mux.NewRouter(),
		logger: log,
		config: config,
	}
}
