package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type Producer struct {
	P *kafka.Producer
}

func NewProducer() *Producer {
	kfk, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})

	if err != nil {
		logrus.WithFields(
			logrus.Fields{

				"package": "Producer",
				"func":    "NewProducer",
				"method":  "NewProducer",
			}).Fatalln(err)
	}
	return &Producer{P: kfk}
}
func (p *Producer) Produce(topic string, Value []byte, Key []byte) error {
	return p.P.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic,
			Partition: kafka.PartitionAny},
		Value: Value,
		Key:   Key,
	}, nil)
}
func (p *Producer) Flush(timeoutMs int) int {
	return p.P.Flush(timeoutMs)
}
