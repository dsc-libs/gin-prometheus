package gp

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
	"time"
)

func (gp *GP) HandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()
		status := strconv.Itoa(c.Writer.Status())
		elapsed := float64(time.Since(start)) / float64(time.Second)
		responseSize := float64(c.Writer.Size())

		url := c.Request.URL.Path

		if m, exist := gp.Metric("requests_total"); exist {
			m.IncBy([]string{}, 1.0)
		}
		if m, exist := gp.Metric("requests_total_by_uri"); exist {
			m.IncBy([]string{url, status}, 1.0)
		}
		if m, exist := gp.Metric("request_duration_seconds"); exist {
			m.Observe([]string{url}, elapsed)
		}
		if m, exist := gp.Metric("response_size_bytes"); exist {
			m.Observe([]string{}, responseSize)
		}
	}
}

func (gp *GP) Use(e *gin.Engine) {
	e.Use(gp.HandlerFunc())
	e.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})
}

func (gp *GP) AddDefaultMetrics() error {
	return gp.RegisterMetrics([]*Metric{
		&Metric{
			Type:        Counter,
			Name:        "requests_total",
			Description: "server received request num",
			Labels:      nil,
		},
		&Metric{
			Type:        Counter,
			Name:        "requests_total_by_uri",
			Description: "server received request num by uri",
			Labels:      []string{"uri", "code"},
		},
		&Metric{
			Type:        Histogram,
			Name:        "request_duration_seconds",
			Description: "the time server took to handle the request",
			Labels:      []string{"uri"},
			Buckets:     []float64{0.1, 0.3, 1.2, 5, 10},
		},
		&Metric{
			Type:        Summary,
			Name:        "response_size_bytes",
			Description: "response sizes in bytes",
		},
	})
}
