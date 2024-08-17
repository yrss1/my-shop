package app

import (
	"context"
	"flag"
	"fmt"
	"github.com/yrss1/my-shop/user/internal/config"
	"github.com/yrss1/my-shop/user/internal/handler"
	"github.com/yrss1/my-shop/user/internal/repository"
	"github.com/yrss1/my-shop/user/internal/service/userService"
	"github.com/yrss1/my-shop/user/pkg/log"
	"github.com/yrss1/my-shop/user/pkg/server"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	logger := log.LoggerFromContext(context.Background())

	configs, err := config.New()
	if err != nil {
		logger.Error("ERR_INIT_CONFIGS", zap.Error(err))
		return
	}

	repositories, err := repository.New(repository.WithPostgresStore(configs.POSTGRES.DSN))
	if err != nil {
		logger.Error("ERR_INIT_REPOSITORIES", zap.Error(err))
		return
	}

	userService, err := userService.New(
		userService.WithUserRepository(repositories.User),
	)
	if err != nil {
		logger.Error("ERR_INIT_USER_SERVICE", zap.Error(err))
		return
	}

	handlers, err := handler.New(
		handler.Dependencies{
			Configs:     configs,
			UserService: userService,
		},
		handler.WithHTTPHandler(),
		handler.WithGRPCHandler())
	if err != nil {
		logger.Error("ERR_INIT_HANDLERS", zap.Error(err))
		return
	}

	servers, err := server.New(
		server.WithHTTPServer(handlers.HTTP, configs.APP.Port),
		server.WithGRPCServer(handlers.GRPCServer, configs.APP.GRPCPort),
	)
	if err != nil {
		logger.Error("ERR_INIT_SERVERS", zap.Error(err))
		return
	}
	if err = servers.Run(); err != nil {
		logger.Error("ERR_INIT_SERVERS", zap.Error(err))
		return
	}
	logger.Info("http server started on http://localhost:" + configs.APP.Port + "/swagger/index.html")
	logger.Info("grpc server started on http://localhost:" + configs.APP.GRPCPort)

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the httpServer gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	fmt.Println("gracefully shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	if err = servers.Stop(ctx); err != nil {
		panic(err)
	}

	fmt.Println("running cleanup tasks...")

	fmt.Println("server was successful shutdown.")
}
