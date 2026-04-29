package handlers

import (
	"SE-CarRentalService/db"
)

var getCarsFromDatabase = db.GetCarsFromDatabase
var getCarByID = db.GetCarByID
var createCar = db.CreateCar
var updateCar = db.UpdateCar
var deleteCar = db.DeleteCar
