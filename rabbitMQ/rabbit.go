package rabbitMQ

import (
	"SE-CarRentalService/types"
	"os"
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
	connStr := os.Getenv("RABBIT_URL")
	conn, err := amqp.Dial(connStr)
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
