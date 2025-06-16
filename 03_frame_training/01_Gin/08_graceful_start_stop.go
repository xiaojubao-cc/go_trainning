package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 实现服务优雅启停
func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	//severSimpleStop(router)
	severWithContextStop(router)

}

func severSimpleStop(router *gin.Engine) {
	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		<-signals
		log.Println("server is shutting down")
		if err := server.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	if err := server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("Server closed under request")
		} else {
			log.Fatal("Server closed unexpect")
		}
	}
}

func severWithContextStop(router *gin.Engine) {
	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Println("Server closed under request")
			} else {
				log.Fatal("Server closed unexpect")
			}
		}
	}()
	<-ctx.Done()

	stop()
	// 预留时间让服务器处理完请求
	context, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	if err := server.Shutdown(context); err != nil {
		log.Fatal(err)
	}
}
