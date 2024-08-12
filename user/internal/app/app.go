package app

import (
	"context"
	"flag"
	"fmt"
	"github.com/yrss1/my-shop/tree/main/user/internal/config"
	"github.com/yrss1/my-shop/tree/main/user/internal/handler"
	"github.com/yrss1/my-shop/tree/main/user/internal/repository"
	"github.com/yrss1/my-shop/tree/main/user/internal/service/shop"
	"github.com/yrss1/my-shop/tree/main/user/pkg/server"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	configs, err := config.New()
	if err != nil {
		fmt.Printf("ERR_INIT_CONFIGS: %v", err)
		return
	}

	repositories, err := repository.New(repository.WithPostgresStore(configs.POSTGRES.DSN))
	if err != nil {
		fmt.Printf("ERR_INIT_REPOSITORIES: %v", err)
		return
	}

	shopService, err := shop.New(
		shop.WithUserRepository(repositories.User),
	)
	if err != nil {
		fmt.Printf("ERR_INIT_SHOP_SERVICE: %v", err)
		return
	}
	fmt.Println("--------------------------------")
	fmt.Println(configs.APP.GRPCPort)
	fmt.Println("--------------------------------")
	handlers, err := handler.New(
		handler.Dependencies{
			Configs:     configs,
			ShopService: shopService,
		},
		handler.WithHTTPHandler(),
		handler.WithGRPCHandler())
	if err != nil {
		fmt.Printf("ERR_INIT_HANDLERS: %v", err)
		return
	}

	servers, err := server.New(
		server.WithHTTPServer(handlers.HTTP, configs.APP.Port),
		server.WithGRPCServer(handlers.GRPCServer, configs.APP.GRPCPort),
	)
	if err != nil {
		fmt.Printf("ERR_RUN_SERVERS: %v", err)
		return
	}
	if err = servers.Run(); err != nil {
		fmt.Printf("ERR_RUN_SERVERS: %v", err)
		return
	}
	fmt.Println("http server started on http://localhost:" + configs.APP.Port)
	fmt.Println("grpc server started on http://localhost:" + configs.APP.GRPCPort)

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
