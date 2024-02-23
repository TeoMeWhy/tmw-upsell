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

	r.PUT("/customers", func(c *gin.Context) {
		api.PutCustomer(c)
	})

	r.DELETE("/customers", func(c *gin.Context) {
		api.DeleteCustomer(c)
	})

	r.PUT("/email", func(c *gin.Context) {
		api.PutCustomerEmail(c)
	})

	r.PUT("/addPoints", func(c *gin.Context) {
		api.PutAddCustomerPoints(c)
	})

	r.GET("/transactions", func(c *gin.Context) {
		api.GetCustomerTransactions(c)
	})

	r.Run(":8080")
}
