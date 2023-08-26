package kafka

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	proto2 "github.com/golang/protobuf/proto"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"time"
	"transaction/internal/proto"
)

var (
	bootstrapservers = "localhost"
	groupid          = "myGroup"
	autooffsetreset  = "earliest"
)

type Consumer struct {
	C     *kafka.Consumer
	cache *cache.Cache
}

func NewConsumer(Cache *cache.Cache) Consumer {
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
	return Consumer{C: kfk, cache: Cache}
}

func (C Consumer) ConsumerStart(ctx context.Context) error {
	C.C.SubscribeTopics([]string{"Add", "Sub"}, nil)
	// A signal handler or similar could be used to set this to false to break the loop.

	go func(ctx context.Context) {
		defer C.C.Close()
		for {
			msg, err := C.C.ReadMessage(time.Millisecond)
			if err == nil {
				fmt.Println("get")
				var message proto.Message
				err := proto2.Unmarshal(msg.Value, &message)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Printf("Message on %v: msg:%v:\n", msg.TopicPartition.Topic, message)
			} else if !err.(kafka.Error).IsTimeout() {
				// The client will automatically try to recover from all errors.
				// Timeout is not considered an error because it is raised by
				// ReadMessage in absence of messages.
				fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			}
			select {
			case <-ctx.Done():
				logrus.WithFields(logrus.Fields{"func": "ConsumerStart"}).Fatalf("faild Consumer")
				return
			default:
			}
		}
	}(ctx)
	//
	return nil

}
