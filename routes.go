package main

import (
	_ "SE-CarRentalService/docs"
	"SE-CarRentalService/handlers"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(rabbitConnection *amqp.Connection, RabbitChannel *amqp.Channel) {
	r := gin.Default()

	corsAddress := os.Getenv("CORS_ADDRESS")

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			corsAddress,
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	carHandler := handlers.NewCarHandler(rabbitConnection, RabbitChannel)

	r.GET("/swagger/*any", ginSwagger.CustomWrapHandler(&ginSwagger.Config{
		URL: "doc.json",
	}, swaggerFiles.Handler))

	// GET
	r.GET("/cars", carHandler.GetCars)
	r.GET("/car/:id", carHandler.GetCar)

	// POST
	r.POST("/car", carHandler.CreateCar)

	// UPDATE
	r.PUT("/car/:id", carHandler.UpdateCar)

	// UPDATE
	r.DELETE("/car/:id", carHandler.DeleteCar)

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}

}
