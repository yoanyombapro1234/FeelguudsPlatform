package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/wailsapp/wails/lib/logger"
	core_auth_sdk "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-auth-sdk"
	core_logging "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-logging"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/authentication_handler"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/helper"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant"
	"github.com/yoanyombapro1234/FeelguudsPlatform/pkg/api"
	"github.com/yoanyombapro1234/FeelguudsPlatform/pkg/grpc"
	"github.com/yoanyombapro1234/FeelguudsPlatform/pkg/signals"
	"github.com/yoanyombapro1234/FeelguudsPlatform/pkg/version"
	"go.uber.org/zap"
)

func main() {
	// flags definition
	fs := pflag.NewFlagSet("default", pflag.ContinueOnError)
	configureEnvironmentVariables(fs)

	// capture goroutines waiting on synchronization primitives
	runtime.SetBlockProfileRate(1)
	versionFlag := fs.BoolP("ENABLE_VERSION_FROM_FILE", "v", false, "get version number")

	// parse flags
	ParseFlags(fs, versionFlag)
	LoadEnvVariables(fs)
	LoadServiceConfigsFromFile()

	serviceName := viper.GetString("SERVICE_NAME")
	logLevel := viper.GetString("LOG_LEVEL")

	logInstance := core_logging.New(logLevel)
	defer logInstance.ConfigureLogger()
	log := logInstance.Logger

	authenticationComponent := InitializeAuthenticationComponent(log, serviceName)
	merchantAccountComponent := InitializeMerchantAccountComponent(log, authenticationComponent)

	// start stress tests if any
	numStressedCpus := viper.GetInt("NUMBER_OF_STRESSED_CPU")
	dataInMemForStressTestInMb := viper.GetInt("DATA_LOADED_IN_MEMORY_FOR_STRESS_TEST_IN_MB")
	beginStressTest(numStressedCpus, dataInMemForStressTestInMb, log)
	ValidatePorts(fs)
	ValidateDelayOptions(log)
	grpcCfg, srvCfg := LoadServerConfigs(log)

	// start gRPC server
	if grpcCfg.Port > 0 {
		grpcSrv, _ := grpc.NewServer(&grpcCfg, log)
		go grpcSrv.ListenAndServe()
	}

	// log version and port
	log.Info(fmt.Sprintf("Starting %s", serviceName),
		zap.String("version", viper.GetString("VERSION")),
		zap.String("revision", viper.GetString("REVISION")),
		zap.String("port", srvCfg.Port),
	)

	// start HTTP server
	srv, _ := api.NewServer(&srvCfg, log, authenticationComponent, merchantAccountComponent)
	stopCh := signals.SetupSignalHandler()
	srv.ListenAndServe(stopCh)
}

func InitializeMerchantAccountComponent(log *zap.Logger, authenticationComponent *authentication_handler.AuthenticationComponent) *merchant.
AccountComponent {
	host := viper.GetString("MERCHANT_COMPONENT_HOST")
	port := viper.GetInt("MERCHANT_COMPONENT_PORT")
	user := viper.GetString("MERCHANT_COMPONENT_USER")
	password := viper.GetString("MERCHANT_COMPONENT_PASSWORD")
	dbname := viper.GetString("MERCHANT_COMPONENT_DB_NAME")
	stripeApiKey := viper.GetString("STRIPE_API_KEY")
	httpTimeout := viper.GetDuration("HTTP_REQUEST_TIMEOUT_IN_MS")
	refreshUrl := viper.GetString("REFRESH_URL")
	returnUrl := viper.GetString("REFRESH_URL")

	maxDBConnAttempts := viper.GetInt("MAX_DB_CONNECTION_ATTEMPTS")
	maxRetriesPerDBConnectionAttempt := viper.GetInt("MAX_DB_CONNECTION_ATTEMPTS_RETRIES")
	maxDBRetryTimeout := viper.GetDuration("MAX_DB_RETRY_TIMEOUT")
	maxDBSleepInterval := viper.GetDuration("DB_RETRY_SLEEP_INTERVAL")

	params := merchant.AccountParams{
		AuthenticationComponent: authenticationComponent,
		DatabaseConnectionParams: &helper.DatabaseConnectionParams{
			Host:         host,
			User:         user,
			Password:     password,
			DatabaseName: dbname,
			Port:         port,
		},
		DatabaseConnectionMetadataParams: &merchant.DatabaseConnectionMetadataParams{
			MaxDatabaseConnectionAttempts:  maxDBConnAttempts,
			MaxRetriesPerConnectionAttempt: maxRetriesPerDBConnectionAttempt,
			RetryTimeout:                   maxDBRetryTimeout,
			RetrySleepInterval:             maxDBSleepInterval,
		},
		Logger:       log,
		StripeApiKey: &stripeApiKey,
		RefreshUrl:   &refreshUrl,
		ReturnUrl:    &returnUrl,
		HttpTimeout:  httpTimeout,
	}

	log.Info("successfully initialized merchant account component")
	return merchant.NewMerchantAccountComponent(&params)
}

func configureEnvironmentVariables(fs *pflag.FlagSet) {
	fs.Int("HTTP_PORT", 9897, "HTTP port")
	fs.Int("HTTPS_PORT", 9898, "HTTPS port")
	fs.Int("METRICS_PORT", 9899, "metrics port")
	fs.String("GRPC_SERVICE_NAME", "FEELGUUDS_PLATFORM", "service name")
	fs.Int("GRPC_PORT", 9896, "gRPC port")
	fs.Int("GRPC_RPC_DEADLINE_IN_MS", 5, "gRPC deadline in milliseconds")
	fs.Int("GRPC_RPC_RETRIES", 2, "gRPC max operation retries in the face of errors")
	fs.Int("GRPC_RPC_RETRY_TIMEOUT_IN_MS", 100, "gRPC max timeout of retry operation in milliseconds")
	fs.Int("GRPC_RPC_RETRY_BACKOFF_IN_MS", 20, "gRPC backoff in between failed retry operations in milliseconds")
	fs.String("LOG_LEVEL", "info", "log level debug, info, warn, error, flat or panic")
	fs.StringSlice("BACKEND_SERVICE_URLS", []string{}, "backend service URL")
	fs.Duration("HTTP_CLIENT_TIMEOUT_IN_MINUTES", 2*time.Minute, "client timeout duration")
	fs.Duration("HTTP_SERVER_TIMEOUT_IN_SECONDS", 30*time.Second, "server read and write timeout duration")
	fs.Duration("HTTP_SERVER_SHUTDOWN_TIMEOUT_IN_SECONDS", 5*time.Second, "server graceful shutdown timeout duration")
	fs.String("DATA_PATH", "/data", "data local path")
	fs.String("CONFIG_PATH", "", "config dir path")
	fs.String("CERT_PATH", "/go/src/github.com/yoanyombapro1234/FeelguudsPlatform/certificate/cert",
		"certificate path for HTTPS port")
	fs.String("CONFIG_FILE", "config.yaml", "config file name")
	fs.String("UI_PATH", "./ui", "UI local path")
	fs.String("UI_LOGO", "", "UI logo")
	fs.String("UI_COLOR", "#34577c", "UI color")
	fs.String("UI_MESSAGE", fmt.Sprintf("greetings from service v%v", version.VERSION), "UI message")
	fs.Bool("ENABLE_H2C", false, "allow upgrading to H2C")
	fs.Bool("ENABLE_RANDOM_DELAY", false, "between 0 and 5 seconds random delay by default")
	fs.String("RANDOM_DELAY_UNIT", "s", "either s(seconds) or ms(milliseconds")
	fs.Int("RANDOM_DELAY_MIN", 0, "min for random delay: 0 by default")
	fs.Int("RANDOM_DELAY_MAX", 5, "max for random delay: 5 by default")
	fs.Bool("ENABLE_RANDOM_RANDOM_ERROR", false, "1/3 chances of a random response error")
	fs.Bool("SET_SERVICE_UNHEALTHY", false, "when set, healthy state is never reached")
	fs.Bool("SET_SERVICE_UNREADY", false, "when set, ready state is never reached")
	fs.Bool("ENABLE_CPU_STRESS_TEST", false, "enable cpu stress tests")
	fs.Bool("ENABLE_MEMORY_STRESS_TEST", false, "enable memory stress tests")
	fs.Int("NUMBER_OF_STRESSED_CPU", 0, "number of CPU cores with 100 load")
	fs.Int("DATA_LOADED_IN_MEMORY_FOR_STRESS_TEST_IN_MB", 0, "MB of data to load into memory")
	fs.String("CACHE_SERVER_ADDRESS", "", "Redis address in the format <host>:<port>")
	// authentication service specific flags
	fs.String("AUTHN_USERNAME", "feelguuds", "username of authentication client")
	fs.String("AUTHN_PASSWORD", "feelguuds", "password of authentication client")
	fs.String("AUTHN_ISSUER_BASE_URL", "http://localhost", "authentication service issuer")
	fs.String("AUTHN_ORIGIN", "http://localhost", "origin of auth requests")
	fs.String("AUTHN_DOMAINS", "localhost", "authentication service domains")
	fs.String("PRIVATE_BASE_URL", "http://authentication-service",
		"authentication service private url. should be local host if these are not running on docker containers. "+
			"However if running in docker container with a configured docker network, the url should be equal to the service name")
	fs.String("AUTHN_PUBLIC_BASE_URL", "http://localhost", "authentication service public endpoint")
	fs.String("AUTHN_INTERNAL_PORT", "3000", "authentication service port")
	fs.String("AUTHN_EXTERNAL_PORT", "8000", "authentication service external port")
	fs.Bool("ENABLE_AUTHN_PRIVATE_INTEGRATION", true, "enables communication with authentication service")

	// retry specific configurations
	fs.Int("HTTP_MAX_RETRIES", 5, "max retries to perform on failed http calls")
	fs.Duration("HTTP_MIN_RETRY_WAIT_TIME_IN_MS", 150*time.Millisecond, "minimum time to wait between failed calls for retry")
	fs.Duration("HTTP_MAX_RETRY_WAIT_TIME_IN_MS", 300*time.Millisecond, "maximum time to wait between failed calls for retry")
	fs.Duration("HTTP_REQUEST_TIMEOUT_IN_MS", 500*time.Millisecond, "time until a request is seen as timing out")
	// logging specific configurations
	fs.String("SERVICE_NAME", "FEELGUUDS_PLATFORM", "service name")
	fs.Int("DOWNSTREAM_SERVICE_CONNECTION_LIMIT", 8, "max retries to perform while attempting to connect to downstream services")

	// merchant component database connection configurations
	fs.String("MERCHANT_COMPONENT_HOST", "merchant_component_db", "database host string")
	fs.Int("MERCHANT_COMPONENT_PORT", 5432, "database port")
	fs.String("MERCHANT_COMPONENT_USER", "merchant_component", "database user string")
	fs.String("MERCHANT_COMPONENT_PASSWORD", "merchant_component", "database password string")
	fs.String("MERCHANT_COMPONENT_DB_NAME", "merchant_component", "database name")
	fs.String("REFRESH_URL", "http://localhost/v1/merchant-account/refresh-url", "refresh url used as part of stripe onboarding")
	fs.String("RETURN_URL", "http://localhost/v1/merchant-account/return-url", "return url used as part of stripe onboarding")
	fs.Int("MAX_DB_CONNECTION_ATTEMPTS", 2, "max database connection attempts")
	fs.Int("MAX_DB_CONNECTION_ATTEMPTS_RETRIES", 2, "max database connection attempts")
	fs.Duration("MAX_DB_RETRY_TIMEOUT", 500*time.Millisecond, "max time until a db connection request is seen as timing out")
	fs.Duration("DB_RETRY_SLEEP_INTERVAL", 100*time.Millisecond, "max time to sleep in between db connection attempts")

	// shopper component database connection configurations
	fs.String("SHOPPER_COMPONENT_HOST", "shopper_component_db", "database host string")
	fs.Int("SHOPPER_COMPONENT_PORT", 5432, "database port")
	fs.String("SHOPPER_COMPONENT_USER", "shopper_component", "database user string")
	fs.String("SHOPPER_COMPONENT_PASSWORD", "shopper_component", "database password string")
	fs.String("SHOPPER_COMPONENT_DB_NAME", "shopper_component", "database name")

	// stripe specific secrets
	fs.String("STRIPE_API_KEY", "", "stripe api key")
}

// LoadServerConfigs loads server configurations (grpc and http)
func LoadServerConfigs(log *zap.Logger) (grpc.Config, api.Config) {
	// load gRPC server config
	var grpcCfg grpc.Config
	if err := viper.Unmarshal(&grpcCfg); err != nil {
		log.Panic("config unmarshal failed", zap.Error(err))
	}

	// load HTTP server config
	var srvCfg api.Config
	if err := viper.Unmarshal(&srvCfg); err != nil {
		log.Panic("config unmarshal failed", zap.Error(err))
	}

	return grpcCfg, srvCfg
}

// InitializeAuthenticationComponent configures the authentication component
func InitializeAuthenticationComponent(log *zap.Logger, serviceName string) *authentication_handler.AuthenticationComponent {
	authUsername := viper.GetString("AUTHN_USERNAME")
	authPassword := viper.GetString("AUTHN_PASSWORD")
	audience := viper.GetString("AUTHN_DOMAINS")
	privateURL := viper.GetString("PRIVATE_BASE_URL") + ":" + viper.GetString("AUTHN_INTERNAL_PORT")

	origin := viper.GetString("AUTHN_ORIGIN")
	issuer := viper.GetString("AUTHN_ISSUER_BASE_URL") + ":" + viper.GetString("AUTHN_EXTERNAL_PORT")
	httpTimeout := 300 * time.Millisecond

	return authentication_handler.NewAuthenticationComponent(&authentication_handler.AuthenticationParams{
		AuthConfig: &core_auth_sdk.Config{
			Issuer:         issuer,
			PrivateBaseURL: privateURL,
			Audience:       audience,
			Username:       authUsername,
			Password:       authPassword,
			KeychainTTL:    0,
		},
		AuthConnectionConfig: &core_auth_sdk.RetryConfig{
			MaxRetries:       viper.GetInt("HTTP_MAX_RETRIES"),
			MinRetryWaitTime: viper.GetDuration("HTTP_MIN_RETRY_WAIT_TIME_IN_MS"),
			MaxRetryWaitTime: viper.GetDuration("HTTP_MAX_RETRY_WAIT_TIME_IN_MS"),
			RequestTimeout:   viper.GetDuration("HTTP_REQUEST_TIMEOUT_IN_MS"),
		},
		Logger: log,
		Origin: origin,
	}, serviceName, httpTimeout)
}

// ValidateDelayOptions validates random delay options
func ValidateDelayOptions(log *zap.Logger) {
	// validate random delay options
	if viper.GetInt("RANDOM_DELAY_MAX") < viper.GetInt("RANDOM_DELAY_MIN") {
		err := errors.New("`--random-delay-max` should be greater than `--random-delay-min`")
		log.Fatal("please fix configurations", zap.Error(err))
	}

	switch delayUnit := viper.GetString("RANDOM_DELAY_UNIT"); delayUnit {
	case
		"s",
		"ms":
		break
	default:
		err := errors.New("random-delay-unit` accepted values are: s|ms")
		log.Fatal("please fix configurations", zap.Error(err))
	}
}

// ValidatePorts ensures http and https ports are valid
func ValidatePorts(fs *pflag.FlagSet) {
	// validate port
	if _, err := strconv.Atoi(viper.GetString("HTTP_PORT")); err != nil {
		port, _ := fs.GetInt("HTTP_PORT")
		viper.Set("HTTP_PORT", strconv.Itoa(port))
	}

	// validate secure port
	if _, err := strconv.Atoi(viper.GetString("HTTPS_PORT")); err != nil {
		securePort, _ := fs.GetInt("HTTPS_PORT")
		viper.Set("HTTPS_PORT", strconv.Itoa(securePort))
	}
}

// LoadEnvVariables binds a set of flags to and loads environment variables
func LoadEnvVariables(fs *pflag.FlagSet) {
	viper.AddConfigPath("/Users/yoanyomba/go/src/github.com/yoanyombapro1234/FeelguudsPlatform")
	viper.BindPFlags(fs)
	viper.RegisterAlias("BACKEND_SERVICE_URLS", "backend-url")
	viper.SetConfigName("service")
	viper.SetConfigType("env")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	hostname, _ := os.Hostname()
	viper.SetDefault("JWT_SECRET", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9")
	viper.SetDefault("UI_LOGO", "https://raw.githubusercontent.com/stefanprodan/podinfo/gh-pages/cuddle_clap.gif")
	viper.Set("HOSTNAME", hostname)
	viper.Set("VERSION", version.VERSION)
	viper.Set("REVISION", version.REVISION)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Error(err.Error())
		return
	}
}

// ParseFlags parses a set of defined flags
func ParseFlags(fs *pflag.FlagSet, versionFlag *bool) {
	err := fs.Parse(os.Args[1:])
	switch {
	case err == pflag.ErrHelp:
		os.Exit(0)
	case err != nil:
		fmt.Fprintf(os.Stderr, "Error: %s\n\n", err.Error())
		fs.PrintDefaults()
		os.Exit(2)
	case *versionFlag:
		fmt.Println(version.VERSION)
		os.Exit(0)
	}
}

// LoadServiceConfigsFromFile loads service configurations from a file
func LoadServiceConfigsFromFile() {
	// load config from file
	p := viper.GetString("CONFIG_PATH")
	f := viper.GetString("CONFIG_FILE")
	log.Info(p, f)
	if _, fileErr := os.Stat(filepath.Join(viper.GetString("CONFIG_PATH"), viper.GetString("CONFIG_FILE"))); fileErr == nil {
		viper.SetConfigName(strings.Split(viper.GetString("CONFIG_FILE"), ".")[0])
		viper.AddConfigPath(viper.GetString("CONFIG_PATH"))
		if readErr := viper.ReadInConfig(); readErr != nil {
			fmt.Printf("Error reading config file, %v\n", readErr)
		}
	} else {
		fmt.Printf("Error to open config file, %v\n", fileErr)
	}
}

// beginStressTest performs cpu and memory stress tests
func beginStressTest(cpus int, mem int, log *zap.Logger) {
	PerformCpuStressTest(cpus, log)
	PerformMemoryStressTest(mem, log)
}

// PerformMemoryStressTest performs memory stress test
func PerformMemoryStressTest(mem int, log *zap.Logger) {
	var stressMemoryPayload []byte
	var memoryStressErr error

	if mem > 0 {
		path := "/tmp/service.data"
		f, err := os.Create(path)

		if err != nil {
			log.Error("memory stress failed", zap.Error(err))
		}

		if err := f.Truncate(1000000 * int64(mem)); err != nil {
			log.Error("memory stress failed", zap.Error(err))
		}

		stressMemoryPayload, memoryStressErr = ioutil.ReadFile(path)
		if err = f.Close(); err != nil {
			log.Error("failed to close file handle", zap.Error(err))
		}

		if err = os.Remove(path); err != nil {
			log.Error(fmt.Sprintf("%s : path - %s", "failed to remove file at path location", path), zap.Error(err))
		}

		if memoryStressErr != nil {
			log.Error("memory stress failed", zap.Error(err))
		}

		log.Info("starting CPU stress", zap.Int("memory", len(stressMemoryPayload)))
	}
}

// PerformCpuStressTest performs a cpu stress test
func PerformCpuStressTest(cpus int, log *zap.Logger) {
	done := make(chan int)
	if cpus > 0 {
		log.Info("starting CPU stress", zap.Int("cores", cpus))
		for i := 0; i < cpus; i++ {
			go func() {
				for {
					select {
					case <-done:
						return
					default:

					}
				}
			}()
		}
	}
}
