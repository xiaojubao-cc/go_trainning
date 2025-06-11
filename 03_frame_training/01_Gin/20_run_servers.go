package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// 运行多个服务

func main() {
	router := gin.New()
	server1 := createHttpServer(router)
	server2 := createHttpServer(router)
	go func() {
		server1.ListenAndServe()
	}()
	go func() {
		server2.ListenAndServe()
	}()
}

func createHttpServer(handler http.Handler) *http.Server {
	server := &http.Server{
		Addr:           ":8080",
		Handler:        handler,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return server
}
