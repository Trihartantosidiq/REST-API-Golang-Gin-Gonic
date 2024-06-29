package main

import (
	"backendGO/controllers"
	"backendGO/database"
	"backendGO/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	database.InitDB()

	router := gin.Default()

	// User routes
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)

	// Product routes
	protected := router.Group("/products")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/", controllers.GetProducts)
		protected.GET("/:id", controllers.GetProductByID)
		protected.POST("/", controllers.CreateProduct)
		protected.PUT("/:id", controllers.UpdateProduct)
		protected.DELETE("/:id", controllers.DeleteProduct)
	}

	router.Run(":8080")
}
