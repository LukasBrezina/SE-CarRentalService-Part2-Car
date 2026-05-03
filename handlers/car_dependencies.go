package handlers

import (
	"SE-CarRentalService/db"
	"SE-CarRentalService/rabbitMQ"
	"SE-CarRentalService/services"
)

var GRCPConnection *services.ConverterClient

var getCarsFromDatabase = db.GetCarsFromDatabase
var getCarByID = db.GetCarByID
var createCar = db.CreateCar
var updateCar = db.UpdateCar
var deleteCar = db.DeleteCar
var verifyToken = rabbitMQ.CheckTokenViaRabbit
var convertCurrency = GRCPConnection.ConvertCurrency
