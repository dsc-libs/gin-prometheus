package gp

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func Test(t *testing.T)  {
	r := gin.New()
	gp := GetGP()
	gp.Use(r)

	err := gp.AddDefaultMetrics()
	if err != nil {
		return
	}

	r.GET("/:id/:ttt", func(c *gin.Context) {
		c.JSON(200, "Hello world!")
	})

	r.Run(":8989")

}