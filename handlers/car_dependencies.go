package handlers

import (
	"SE-CarRentalService/db"
	"SE-CarRentalService/rabbitMQ"
	"SE-CarRentalService/services"
	"errors"
)

var GRCPConnection *services.ConverterClient

var getCarsFromDatabase = db.GetCarsFromDatabase
var getCarByID = db.GetCarByID
var createCar = db.CreateCar
var updateCar = db.UpdateCar
var deleteCar = db.DeleteCar
var verifyToken = rabbitMQ.CheckTokenViaRabbit

var convertCurrency = func(from string, amount float64, to string) (float64, string, error) {
	if GRCPConnection == nil {
		return 0, "", errors.New("GRCPConnection is nil")
	}

	return GRCPConnection.ConvertCurrency(from, amount, to)
}
