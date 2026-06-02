package kafkaservice

import (
	"context"
	"strconv"

	"github.com/Shoyeb45/server/pkg/config"
	"github.com/Shoyeb45/server/pkg/logger"
	"github.com/Shoyeb45/server/pkg/shared"
	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	writer *kafka.Writer
}

var kafkaClient *KafkaClient

func KafkaInit() error {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:     []string{config.Cfg.KafkaAddress},
		Topic:       config.Cfg.KafkaTopic,
		Balancer:    &kafka.LeastBytes{},
		MaxAttempts: 3,
	})

	kafkaClient = &KafkaClient{writer: writer}
	logger.Log.Info("kafka connected successfully")
	return nil
}

func Produce[T any](userID int32, data T) error {
	message, err := shared.ToString(data)
	
	if err != nil {
		return err
	}

	return kafkaClient.writer.WriteMessages(context.Background(),
		kafka.Message{Value: []byte(message), Key: []byte(strconv.Itoa(int(userID)))},
	)
}

func Close() error {
	return kafkaClient.writer.Close()
}
