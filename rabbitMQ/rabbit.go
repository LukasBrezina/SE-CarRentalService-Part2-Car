package rabbitMQ

import (
	"SE-CarRentalService/types"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func CheckTokenViaRabbit(ch *amqp.Channel, token string) (types.TokenResponse, error) {
	return CallRabbitRPC[types.TokenRequest, types.TokenResponse](
		ch,
		"auth.check_token",
		types.TokenRequest{
			Token: token,
		},
		5*time.Second,
	)
}

func Connect() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, err
	}

	return conn, ch, nil
}
