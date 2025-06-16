package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"time"
)

// 运行多个服务
/*
errgroup.Group:
    1.并发任务管理:通过 errG.Go(func() error { ... }) 启动多个并发任务。自动等待所有任务完成（类似 sync.WaitGroup）。
	2.错误集中处理任意一个任务返回非 nil 错误时，errG.Wait() 会立即返回第一个错误。其他未完成的任务会根据上下文自动取消（需配合 context 使用）。
	3.上下文传播可通过 errgroup.WithContext(ctx) 创建带上下文的 Group。当某个任务失败时，自动触发上下文取消，终止其他任务。
*/
var (
	errG errgroup.Group
)

func main() {
	router1 := gin.New()
	router2 := gin.New()
	server1 := createHttpServer(":8080", router1)
	server2 := createHttpServer(":8081", router2)
	errG.Go(func() error {
		return server1.ListenAndServe()
	})
	errG.Go(func() error {
		return server2.ListenAndServe()
	})
	if err := errG.Wait(); err != nil {
		log.Fatal(err)
	}
}

func createHttpServer(port string, router *gin.Engine) *http.Server {
	router.Use(gin.Recovery())
	router.GET("/startMultipleServer", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "startMultipleServer",
		})
	})
	server := &http.Server{
		Addr:           port,
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return server
}
