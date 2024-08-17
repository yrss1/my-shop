package app

import (
	"context"
	"flag"
	"fmt"
	"github.com/yrss1/my-shop/api-gateway/intrenal/config"
	"github.com/yrss1/my-shop/api-gateway/intrenal/handler"
	"github.com/yrss1/my-shop/api-gateway/pkg/log"
	"github.com/yrss1/my-shop/api-gateway/pkg/server"
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
		fmt.Printf("ERR_INIT_CONFIGS: %v", err)
		return
	}

	handlers, err := handler.New(
		handler.Dependencies{
			Configs: configs,
		},
		handler.WithHTTPHandler())
	if err != nil {
		logger.Error("ERR_INIT_HANDLERS", zap.Error(err))
		return
	}

	servers, err := server.New(
		server.WithHTTPServer(handlers.HTTP, configs.APP.Port),
	)
	if err != nil {
		logger.Error("ERR_INIT_SERVERS", zap.Error(err))
		return
	}
	if err = servers.Run(); err != nil {
		logger.Error("ERR_INIT_SERVERS", zap.Error(err))
		return
	}

	logger.Info("http server started on http://localhost:" + configs.APP.Port)

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", configs.APP.Timeout, "the duration for which the httpServer gracefully wait for existing connections to finish - e.g. 15s or 1m")
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
