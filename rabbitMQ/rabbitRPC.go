package rabbitMQ

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func CallRabbitRPC[TRequest any, TResponse any](
	ch *amqp.Channel,
	queueName string,
	request TRequest,
	timeout time.Duration,
) (TResponse, error) {
	var response TResponse

	replyQueue, err := ch.QueueDeclare(
		"",
		false,
		true,
		true,
		false,
		nil,
	)
	if err != nil {
		return response, err
	}

	msgs, err := ch.Consume(
		replyQueue.Name,
		"",
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return response, err
	}

	correlationID := uuid.New().String()

	body, err := json.Marshal(request)
	if err != nil {
		return response, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = ch.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: correlationID,
			ReplyTo:       replyQueue.Name,
			Body:          body,
		},
	)
	if err != nil {
		return response, err
	}

	for {
		select {
		case msg := <-msgs:
			if msg.CorrelationId != correlationID {
				continue
			}

			err := json.Unmarshal(msg.Body, &response)
			if err != nil {
				return response, err
			}

			return response, nil

		case <-ctx.Done():
			return response, fmt.Errorf("rabbitmq rpc timeout")
		}
	}
}
