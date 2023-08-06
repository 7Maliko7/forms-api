package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitoc "github.com/go-kit/kit/tracing/opencensus"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/7Maliko7/forms-api/internal/config"
	"github.com/7Maliko7/forms-api/internal/middleware"
	"github.com/7Maliko7/forms-api/internal/service"
	formsvc "github.com/7Maliko7/forms-api/internal/service/forms"
	"github.com/7Maliko7/forms-api/internal/transport"
	httptransport "github.com/7Maliko7/forms-api/internal/transport/http"
	"github.com/7Maliko7/forms-api/pkg/broker/driver/rabbitmq"
	"github.com/7Maliko7/forms-api/pkg/db/driver/postgres"
	"github.com/7Maliko7/forms-api/pkg/oc"
	"github.com/7Maliko7/forms-api/pkg/storage"
	"github.com/7Maliko7/forms-api/pkg/storage/driver/filesystem"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "c", "", "Custom config path")
	flag.Parse()
}

func main() {
	var logger log.Logger
	{
		logger = log.NewJSONLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"svc", "forms-api",
			"ts", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	appConfig, err := config.New(configPath)
	if err != nil {
		level.Error(logger).Log(err.Error())
		os.Exit(1)
	}

	var db *sql.DB
	{
		db, err = sql.Open("postgres", appConfig.RWDB.ConnectionString)
		if err != nil {
			level.Error(logger).Log("msg", err.Error())
			os.Exit(1)
		}
	}
	defer db.Close()

	var broker *rabbitmq.Broker
	{
		broker, err = rabbitmq.NewBroker(appConfig.Broker.ConnectionString)
		if err != nil {
			level.Error(logger).Log("msg", err.Error())
			os.Exit(1)
		}
	}
	defer broker.Close()

	var Storage storage.Storager
	Storage = &filesystem.FileSystem{DefaultPermissions: appConfig.Storage.Filesystem.DefaultPermissions, FilePath: appConfig.Storage.Filesystem.Path}

	var svc service.Service
	{
		repository, err := postgres.New(db)
		if err != nil {
			level.Error(logger).Log("exit", err.Error())
			os.Exit(1)
		}

		svc = formsvc.NewService(repository, logger, broker, Storage)
		svc = middleware.LoggingMiddleware(logger)(svc)
	}

	var endpoints transport.Endpoints
	{
		endpoints = transport.MakeEndpoints(svc)
		endpoints = transport.Endpoints{
			Save:        oc.ServerEndpoint("Save")(endpoints.Save),
			GetForm:     oc.ServerEndpoint("GetForm")(endpoints.GetForm),
			GetFormList: oc.ServerEndpoint("GetFormList")(endpoints.GetFormList),
		}
	}

	var h http.Handler
	{
		ocTracing := kitoc.HTTPServerTrace()
		serverOptions := []kithttp.ServerOption{ocTracing}
		h = httptransport.NewService(endpoints, serverOptions, logger)
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		level.Info(logger).Log("transport", "HTTP", "addr", appConfig.ListenAddress)
		server := &http.Server{
			Addr:    appConfig.ListenAddress,
			Handler: h,
		}
		errs <- server.ListenAndServe()
	}()

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	level.Error(logger).Log("exit", <-errs)
}
