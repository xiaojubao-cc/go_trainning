package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"io"
	"net/http"
	"time"
)

/*
context.DataFromReader() 是 Gin 框架中用于流式传输数据的核心方法，
其作用是通过 io.Reader 直接将数据流写入 HTTP 响应体，无需将完整数据加载到内存。
特性                                             说明
零内存压力                  数据直接从 io.Reader 流式传输到 TCP 层，不占用应用内存
支持分块传输编码             当 contentLength=-1 时自动启用 Transfer-Encoding: chunked
实时性                     适用于需要实时生成/传输数据的场景（如 SSE、日志流）
资源自动回收                Gin 会自动关闭 reader（如果实现了 io.Closer 接口）
*/
var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

type rateLimitedReader struct {
	r       io.Reader
	limiter *rate.Limiter
}

func handleDownload(c *gin.Context) {
	response, err := httpClient.Get("https://raw.githubusercontent.com/gin-gonic/logo/master/color.png")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"error": "upstream unavailable",
		})
		return
	}
	//关闭响应流
	defer response.Body.Close()
	contentType := response.Header.Get("Content-Type")
	reader := response.Body

	//限流
	limitedReader := &rateLimitedReader{
		r:       reader,
		limiter: rate.NewLimiter(rate.Limit(1024*1024), 2*1024*1024),
	}

	//使用流式下载图片
	c.DataFromReader(
		http.StatusOK,
		response.ContentLength,
		contentType,
		limitedReader,
		map[string]string{
			"Content-Disposition": `attachment; filename="gopher.png"; filename*=UTF-8''gopher.png`,
		},
	)
}

func (r *rateLimitedReader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	if err != nil {
		return
	}
	if err := r.limiter.WaitN(context.Background(), n); err != nil {
		return 0, err
	}
	return
}

func main() {
	router := gin.Default()
	router.GET("/download", handleDownload)
	router.Run(":8080")
}
