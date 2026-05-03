package main

import (
	"SE-CarRentalService/db"
	"SE-CarRentalService/handlers"
	"SE-CarRentalService/rabbitMQ"
	"SE-CarRentalService/services"
	"log"

	"github.com/joho/godotenv"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}

	GRCPConnection := services.NewConverterClientRetry()
	db.Connect(GRCPConnection)
	rabbitConnection, rabbitChannel, err := rabbitMQ.Connect()
	handlers.GRCPConnection = GRCPConnection

	if err != nil {
		log.Println(err)
	}
	SetupRouter(rabbitConnection, rabbitChannel)
}
