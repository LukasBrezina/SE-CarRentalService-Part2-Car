package db

import (
	"SE-CarRentalService/services"
	"SE-CarRentalService/types"
	"database/sql"
	"fmt"
)

func GetCarsFromDatabase() ([]types.Car, error) {
	database := DATABASE

	rows, err := database.Query(`
        SELECT id, model, brand, collectAt, accountId, year, price, ps
        FROM car
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carList []types.Car

	for rows.Next() {
		car := types.Car{}

		err := rows.Scan(
			&car.ID,
			&car.Model,
			&car.Brand,
			&car.CollectAt,
			&car.AccountId,
			&car.Year,
			&car.Price,
			&car.PS,
		)
		if err != nil {
			return nil, err
		}
		fmt.Println(car.Price)
		carList = append(carList, car)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return carList, nil
}

func GetCarByID(id int) (types.Car, error) {
	database := DATABASE

	row := database.QueryRow("SELECT id, model, brand, collectat, accountid, year, price, ps FROM car WHERE id = $1", id)

	currentCar := types.Car{}
	if err := row.Scan(
		&currentCar.ID,
		&currentCar.Model,
		&currentCar.Brand,
		&currentCar.CollectAt,
		&currentCar.AccountId,
		&currentCar.Year,
		&currentCar.Price,
		&currentCar.PS); err != nil {
		return types.Car{}, err

	}
	return currentCar, nil
}

func CreateCar(c types.CreateCarRequest, currency string) (types.Car, error) {
	database := DATABASE
	var newCar types.Car
	price, _ := services.ConvertCurrency(currency, c.Price, "USD")
	err := database.QueryRow(`
		INSERT INTO car (model, brand, collectAt, accountId, year, price, ps)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, model, brand, collectAt, accountId, year, price, ps
	`,
		c.Model,
		c.Brand,
		c.CollectAt,
		0,
		c.Year,
		price,
		c.PS,
	).Scan(
		&newCar.ID,
		&newCar.Model,
		&newCar.Brand,
		&newCar.CollectAt,
		&newCar.AccountId,
		&newCar.Year,
		&newCar.Price,
		&newCar.PS,
	)
	if err != nil {
		return types.Car{}, err
	}

	return newCar, nil
}
func UpdateCar(id int, c types.Car, account types.Account) error {
	database := DATABASE
	var result sql.Result
	var err error
	price, _ := services.ConvertCurrency(account.Currency, c.Price, "USD")
	if account.IsAdmin {
		result, err = database.Exec(`
        UPDATE car
        SET model = $1, brand = $2, collectAt = $3, accountId = $4, year = $5, price = $6, ps = $7
        WHERE id = $8
    `,
			c.Model,
			c.Brand,
			c.CollectAt,
			c.AccountId,
			c.Year,
			price,
			c.PS,
			id,
		)
	} else {
		result, err = database.Exec(`
		UPDATE car
		SET accountId = CASE
			WHEN accountId = 0 THEN $1
			WHEN accountId = $2 THEN 0
			ELSE accountId
		END
		WHERE id = $3
`,
			account.ID,
			account.ID,
			id,
		)
	}

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("car not found")
	}

	return nil
}

func DeleteCar(id int) error {
	database := DATABASE
	_, err := database.Exec("DELETE FROM car WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
