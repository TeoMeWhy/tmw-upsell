package main

import (
	"points_mgmt/api"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.GET("/customers", func(c *gin.Context) {
		api.GetCustomer(c)
	})

	r.POST("/customers", func(c *gin.Context) {
		api.PostCustomer(c)
	})

	r.PUT("/addPoints", func(c *gin.Context) {
		api.PutAddUserPoints(c)
	})

	r.Run(":8080")
}
