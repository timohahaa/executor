package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/timohahaa/executor/config"
	v1 "github.com/timohahaa/executor/internal/controllers/http/v1"
	"github.com/timohahaa/executor/internal/repository"
	"github.com/timohahaa/executor/internal/service"
	"github.com/timohahaa/executor/pkg/httpserver"
	"github.com/timohahaa/executor/pkg/logger"
	"github.com/timohahaa/executor/pkg/validator"
	"github.com/timohahaa/postgres"
)

func Run(configFilePath string) {
	cfg, err := config.New(configFilePath)
	if err != nil {
		log.Fatalf("error reading config: %v\n", err)
	}

	mainLogger := logger.New("internal.log", cfg.Server.LogPath)
	httpLogger := logger.New("requests.log", cfg.Server.LogPath)

	// database
	mainLogger.Info("initializing postgres connection...")
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxConnPoolSize(cfg.PG.ConnPoolSize))
	if err != nil {
		mainLogger.WithFields(logrus.Fields{"error": err}).Fatal("error connecting to postgres")
	}

	mainLogger.Info("initializing redis connection...")
	opts, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		mainLogger.WithFields(logrus.Fields{"error": err}).Fatal("error connecting to redis")
	}
	rdb := redis.NewClient(opts)

	// транспортный слой
	mainLogger.Info("initializing repositories...")
	commandRepo := repository.NewCommandRepository(pg, rdb, cfg.Redis.PidTTL, mainLogger)

	// слой БЛ
	mainLogger.Info("initializing services...")
	commandService := service.NewCommandService(commandRepo, cfg.General.DefaultShellPath, mainLogger)

	// handlers and routes
	mainLogger.Info("initializing handlers and routes...")
	handler := v1.NewRouter(commandService, httpLogger)
	handler.Validator = validator.New()

	mainLogger.Infof("starting http server...")
	server := httpserver.New(handler, httpserver.Port(cfg.Server.Port))

	// gracefull shutdown
	mainLogger.Info("configuring gracefull shutdown...")
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	mainLogger.WithFields(logrus.Fields{"port": cfg.Server.Port}).Info("server started!")

	<-shutdownChan

	mainLogger.Info("shutting down...")
	err = server.Shutdown()
	if err != nil {
		mainLogger.WithFields(logrus.Fields{"error": err}).Fatal("error shutting down the server")
	}
}
