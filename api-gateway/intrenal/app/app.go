package app

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yrss1/my-shop/tree/main/api-gateway/intrenal/config"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func Run() {
	configs, err := config.New()
	if err != nil {
		fmt.Printf("ERR_INIT_CONFIGS: %v", err)
		return
	}
	r := gin.Default()

	services := map[string]string{
		"orders":   configs.API.Order,
		"payments": configs.API.Payment,
		"products": configs.API.Product,
		"users":    configs.API.User,
	}

	for path, target := range services {
		registerRoutes(r, path, target)
	}

	server := &http.Server{
		Addr:    ":" + configs.APP.Port,
		Handler: r,
	}

	go func() {
		fmt.Println("API Gateway started on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("ListenAndServe(): %s", err)
		}
	}()

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the httpServer gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	fmt.Println("gracefully shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server Shutdown Failed:%+s", err)
	}

	fmt.Println("Server exited properly")
}

func registerRoutes(r *gin.Engine, basePath, target string) {
	r.Any(basePath+"/*path", proxy(basePath, target))
}

func proxy(basePath, target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем полный путь запроса
		originalPath := c.Param("path")

		// Удаляем базовый путь из исходного запроса
		trimmedPath := strings.TrimPrefix(originalPath, strings.TrimSuffix(basePath, "/"))

		// Формируем целевой URL с новым путем
		targetURL := fmt.Sprintf("%s/%s?%s", target, strings.TrimLeft(trimmedPath, "/"), c.Request.URL.RawQuery)
		fmt.Printf("Proxying request to: %s\n", targetURL) // Логирование целевого URL

		// Создаем новый запрос
		req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		req.Header = c.Request.Header

		// Выполняем запрос к целевому серверу
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		// Перенаправляем заголовки ответа
		for key, values := range resp.Header {
			for _, value := range values {
				c.Writer.Header().Add(key, value)
			}
		}

		// Перенаправляем статус-код ответа
		c.Writer.WriteHeader(resp.StatusCode)

		// Копируем тело ответа
		_, err = io.Copy(c.Writer, resp.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
	}
}
