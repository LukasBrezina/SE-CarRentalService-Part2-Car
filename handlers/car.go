package handlers

import (
	"SE-CarRentalService/rabbitMQ"
	"SE-CarRentalService/services"
	"SE-CarRentalService/types"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

type CarHandler struct{}

var RabbitConnection *amqp.Connection
var RabbitChannel *amqp.Channel
var GRCPConnection *services.ConverterClient

func NewCarHandler(grcpconnection *services.ConverterClient) *CarHandler {
	GRCPConnection = grcpconnection
	return &CarHandler{}
}

// GetCars godoc
// @Summary Get all cars
// @Description Returns all cars, optionally converted to the account currency
// @Tags car
// @Produce json
// @Success 200 {array} types.Car
// @Failure 400 {object} map[string]string
// @Router /cars [get]
func (h *CarHandler) GetCars(c *gin.Context) {

	response, err := verifyToken(RabbitChannel, c.Request.Header.Get("Authorization"))
	if err != nil {
		log.Println("Error verifying token:", err)
		c.JSON(400, gin.H{"error": "invalid token"})
		return
	}

	account := response.Account
	valid := response.Valid
	var currency string
	if valid {
		currency = account.Currency
	} else {
		currency = "EUR"
	}

	cars, err := getCarsFromDatabase()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var carsToSend []types.Car
	for i := range cars {
		newPrice, _, err := convertCurrency("USD", float64(cars[i].Price), currency)
		if err != nil {
			log.Println(err)
			log.Println(cars[i])
			continue
		}
		log.Println(newPrice)
		cars[i].Price = float32(newPrice)
		carsToSend = append(carsToSend, cars[i])
	}
	c.JSON(200, carsToSend)
}

// GetCar godoc
// @Summary Get car by id
// @Description Returns one car by id
// @Tags car
// @Produce json
// @Param id path int true "Car ID"
// @Success 200 {object} types.Car
// @Failure 400 {object} map[string]string
// @Router /car/{id} [get]
func (h *CarHandler) GetCar(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	response, _ := verifyToken(RabbitChannel, c.Request.Header.Get("Authorization"))

	account := response.Account
	valid := response.Valid

	var currency string
	if valid {
		currency = account.Currency
	} else {
		currency = "EUR"
	}

	car, err := getCarByID(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	newPrice, _, _ := convertCurrency("USD", float64(car.Price), currency)

	car.Price = float32(newPrice)

	c.JSON(200, car)
}

// CreateCar godoc
// @Summary Create car
// @Description Creates a new car
// @Tags car
// @Accept json
// @Produce json
// @Param body body types.CreateCarRequest true "Car data"
// @Success 200 {object} types.Car
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /car [post]
func (h *CarHandler) CreateCar(c *gin.Context) {

	response, _ := verifyToken(RabbitChannel, c.Request.Header.Get("Authorization"))

	account := response.Account
	valid := response.Valid
	strErr := response.Error
	if !valid || !account.IsAdmin {
		message := "invalid token"
		if strErr != "" {
			message = strErr
		}

		c.JSON(401, gin.H{"error": message})
		return
	}

	var newCar types.CreateCarRequest

	if err := c.ShouldBindJSON(&newCar); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	createdCar, err := createCar(newCar, account.Currency)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	c.JSON(200, createdCar)
}

// UpdateCar godoc
// @Summary Update car
// @Description Updates a car by id
// @Tags car
// @Accept json
// @Produce json
// @Param id path int true "Car ID"
// @Param body body types.Car true "Updated car"
// @Success 200 {object} types.Car
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /car/{id} [put]
func (h *CarHandler) UpdateCar(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	response, _ := verifyToken(RabbitChannel, c.Request.Header.Get("Authorization"))

	account := response.Account
	valid := response.Valid
	strErr := response.Error
	if !valid {
		message := "invalid token"
		if strErr != "" {
			message = strErr
		}

		c.JSON(401, gin.H{"error": message})
		return
	}

	var updatedCar types.Car
	if err := c.ShouldBindJSON(&updatedCar); err != nil {
		c.JSON(400, gin.H{"error": "invalid body"})
		return
	}

	err = updateCar(id, updatedCar, account)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	updatedCar.ID = id

	c.JSON(200, updatedCar)
}

// DeleteCar godoc
// @Summary Delete car
// @Description Deletes a car by id
// @Tags car
// @Produce json
// @Param id path int true "Car ID"
// @Success 200 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /car/{id} [delete]
func (h *CarHandler) DeleteCar(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
	}

	response, _ := rabbitMQ.CheckTokenViaRabbit(RabbitChannel, c.Request.Header.Get("Authorization"))

	account := response.Account
	valid := response.Valid

	if !valid || !account.IsAdmin {
		c.JSON(401, gin.H{"error": "invalid token"})
		return
	}
	err = deleteCar(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, gin.H{"id": id})
}
