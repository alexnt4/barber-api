package http

import (
	"github.com/alexnt4/barber-api/internal/service"
	"github.com/gin-gonic/gin"
)

func NewRouter(apptSvc *service.AppointmentService, prodSvc *service.ProductService) *gin.Engine {
	r := gin.Default()

	// Middleware de CORS basico
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/healt", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "barberia-api",
		})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Apointments routes
		appts := v1.Group("/appointments")
		{
			apptHandler := NewAppoinmentHandler(apptSvc)
			appts.POST("", apptHandler.Create)
			appts.GET("", apptHandler.List)
			appts.GET("/:id", apptHandler.Get)
			appts.PUT("/:id", apptHandler.Update)
			appts.DELETE("/:id", apptHandler.Delete)
			appts.GET("/:id/total", apptHandler.GetTotal)
		}

		// Product routes
		products := v1.Group("/products")
		{
			prodHandler := NewProductHandler(prodSvc)
			products.POST("", prodHandler.Create)
			products.GET("", prodHandler.List)
			products.GET("/:id", prodHandler.Get)
			products.PUT("/:id", prodHandler.Update)
			products.DELETE("/:id", prodHandler.Delete)
		}
	}

	return r
}
