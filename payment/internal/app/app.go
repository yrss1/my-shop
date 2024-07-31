package app

import (
	"context"
	"flag"
	"fmt"
	"github.com/yrss1/my-shop/tree/main/payment/internal/config"
	"github.com/yrss1/my-shop/tree/main/payment/internal/handler"
	"github.com/yrss1/my-shop/tree/main/payment/internal/provider/epay"
	"github.com/yrss1/my-shop/tree/main/payment/internal/repository"
	"github.com/yrss1/my-shop/tree/main/payment/internal/service/epayment"
	"github.com/yrss1/my-shop/tree/main/payment/pkg/server"
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
	EpayClient, err := epay.New(epay.Credentials{
		URL:            configs.EPAY.URL,
		Login:          configs.EPAY.Login,
		Password:       configs.EPAY.Password,
		OAuthURL:       configs.EPAY.OAuthURL,
		PaymentPageURL: configs.EPAY.PaymentPageURL,
		GlobalToken:    epay.TokenResponse{},
	})
	if err != nil {
		fmt.Printf("ERR_INIT_CLIENTS: %v", err)
		return
	}

	epayService, err := epayment.New(
		epayment.WithPaymentRepository(repositories.Payment),
		epayment.WithEpayClient(EpayClient),
	)
	if err != nil {
		fmt.Printf("ERR_INIT_EPAY_SERVICE: %v", err)
		return
	}

	handlers, err := handler.New(
		handler.Dependencies{
			Configs:     configs,
			EpayService: epayService,
		},
		handler.WithHTTPHandler())
	if err != nil {
		fmt.Printf("ERR_INIT_HANDLERS: %v", err)
		return
	}

	servers, err := server.New(server.WithHTTPServer(handlers.HTTP, configs.APP.Port))
	if err != nil {
		fmt.Printf("ERR_RUN_SERVERS: %v", err)
		return
	}
	if err = servers.Run(); err != nil {
		fmt.Printf("ERR_RUN_SERVERS: %v", err)
		return
	}
	fmt.Println("http server started on http://localhost:" + configs.APP.Port)

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
