package services

import (
	"SE-CarRentalService/services/proto"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type ConverterClient struct {
	conn   *grpc.ClientConn
	client proto.CurrencyConverterServiceClient
}

func NewConverterClientRetry() *ConverterClient {
	for {
		converter, err := NewConverterClient()
		if err == nil {
			return converter
		}

		fmt.Printf("Failed to connect to currency converter, retrying in 3 seconds: %v\n", err)
		time.Sleep(3 * time.Second)
	}
}

func NewConverterClient() (*ConverterClient, error) {
	host := os.Getenv("GRPCHOST")
	port := os.Getenv("GRPCPORT")

	if host == "" {
		return nil, fmt.Errorf("GRPCHOST is not set")
	}

	if port == "" {
		return nil, fmt.Errorf("GRPCPORT is not set")
	}

	target := host + ":" + port

	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}

	return &ConverterClient{
		conn:   conn,
		client: proto.NewCurrencyConverterServiceClient(conn),
	}, nil
}

func (c *ConverterClient) Close() error {
	return c.conn.Close()
}

func (c *ConverterClient) ConvertCurrency(initialCurrency string, initialAmount float64, targetCurrency string) (float64, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	username := os.Getenv("BE_USER")
	password := os.Getenv("BE_PASSWORD")

	if username == "" {
		return 0, "", fmt.Errorf("BE_USER is not set")
	}

	if password == "" {
		return 0, "", fmt.Errorf("BE_PASSWORD is not set")
	}

	authValue := "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authValue)

	reply, err := c.client.ConvertCurrency(ctx, &proto.ConvertCurrencyRequest{
		InitialCurrency: initialCurrency,
		InitialAmount:   initialAmount,
		TargetCurrency:  targetCurrency,
	})
	if err != nil {
		return 0, "", fmt.Errorf("ConvertCurrency failed: %w", err)
	}

	return reply.ConvertedAmount, reply.TargetCurrency, nil
}
