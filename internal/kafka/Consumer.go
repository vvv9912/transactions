package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

var (
	bootstrapservers = "localhost:9092"
	groupid          = "myGroup"
	autooffsetreset  = "earliest"
)

type Consumer struct {
	C     *kafka.Consumer
	Cache *cache.Cache
}

func NewConsumer() Consumer {
	//&

	kfk, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapservers,
		"group.id":          groupid,
		"auto.offset.reset": autooffsetreset,
	})
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"package": "Consumer",
				"func":    "NewConsumer",
				"method":  "NewConsumer",
			}).Fatalln(err)
	}
	return Consumer{C: kfk}
}

func (C *Consumer) ConsumerStart() error {
	//
	return nil
}
