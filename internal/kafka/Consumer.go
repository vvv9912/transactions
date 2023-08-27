package kafka

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	proto2 "github.com/golang/protobuf/proto"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"time"
	"transaction/internal/model"
	"transaction/internal/proto"
)

var (
	bootstrapservers = "localhost"
	groupid          = "myGroup"
	autooffsetreset  = "earliest"
)

const (
	StateAdddb string = "Add"
	StateSubdb        = "Sub"
)

type UsersStorager interface {
	AddUsers(ctx context.Context, Users model.Users) (int64, error)
	CheckId(ctx context.Context, idUser int64) (int64, error)
	AddAccountById(ctx context.Context, id int64, account float64) error
	SetAccountById(ctx context.Context, id int64, account float64) error
	GetAccount(ctx context.Context, idUser int64) (float64, error)
}
type Consumer struct {
	UsersStorage UsersStorager
	C            *kafka.Consumer
	cache        *cache.Cache
}

func NewConsumer(Cache *cache.Cache, storager UsersStorager) Consumer {
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
	return Consumer{C: kfk, cache: Cache, UsersStorage: storager}
}

func (C Consumer) ConsumerStart(ctx context.Context) error {
	err := C.C.SubscribeTopics([]string{"Add", "Sub"}, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{"func": "ConsumerStart"}).Fatalf("Add to topics: %v", err)
		return err
	}
	go func(ctx context.Context) {
		for {
			msg, err := C.C.ReadMessage(time.Millisecond)
			if err == nil {
				fmt.Println("get")

				var message proto.Message
				err := proto2.Unmarshal(msg.Value, &message)

				if err != nil {
					fmt.Println(err)
				}
				//addBD
				//		a := 12
				//		b, _ := C.cache.Get(message.Numtrans)
				//		C.cache.Set(message.Numtrans, a, cache.DefaultExpiration)
				//		c, _ := C.cache.Get(message.Numtrans)
				switch *msg.TopicPartition.Topic {
				case StateAdddb:

					id, err := C.UsersStorage.CheckId(ctx, message.Id)
					if err != nil {
						if err.Error() == "sql: no rows in result set" {
							err = nil
							logrus.WithField("sql: no rows in result set, add users", nil).Warning(err)
							id, err = C.UsersStorage.AddUsers(ctx, model.Users{UserID: message.Id, Account: 0})
							if err != nil {
								return
							}
							logrus.WithField("Add user with id:", message.Id)
							id, err = C.UsersStorage.CheckId(ctx, message.Id)
							if err != nil {
								return
							}
						} else {
							logrus.WithFields(logrus.Fields{"package": "Consumer", "func": "ConsumerStart", "method": "CheckId"}).Fatalln(err)
							return
						}
					}
					err = C.UsersStorage.AddAccountById(ctx, id, message.Account)
					if err != nil {
						logrus.WithFields(logrus.Fields{"package": "Consumer", "func": "ConsumerStart", "method": "AddAccountById"}).Fatalln(err)
						return
					}
					C.cache.Set(message.Numtrans, model.Caches{ID: id, NumberTransaction: message.Numtrans, Status: 1}, cache.DefaultExpiration)
					// Проверка есть ли такой id
					// если нет создать и сразу добавить

					// Обновить баланс
				case StateSubdb:
					id, err := C.UsersStorage.CheckId(ctx, message.Id)
					if err != nil {
						if err.Error() == "sql: no rows in result set" {
							err = nil
							logrus.WithField("sql: no rows in result set, sub users", nil).Warning(err)
						} else {

							logrus.WithFields(logrus.Fields{"package": "Consumer", "func": "ConsumerStart", "method": "CheckId"}).Fatalln(err)
							return
						}
					} else {
						account, err := C.UsersStorage.GetAccount(ctx, id)
						if account >= message.Account {
							err = C.UsersStorage.SetAccountById(ctx, id, (account - message.Account))
							if err != nil {
								logrus.WithFields(logrus.Fields{"package": "Consumer", "func": "ConsumerStart", "method": "AddAccountById"}).Fatalln(err)
								return
							}
							C.cache.Set(message.Numtrans, model.Caches{ID: id, NumberTransaction: message.Numtrans, Status: 2}, cache.DefaultExpiration) //транзакция успешна
						} else {
							//добавить в бд транзакций
							logrus.Infof("account<message.Account")
							C.cache.Set(message.Numtrans, model.Caches{ID: id, NumberTransaction: message.Numtrans, Status: 3}, cache.DefaultExpiration) //транзакция неуспешна
						}

					}
					//Проверка есть ли такой id
					// если нет ошибка транзакции
					// Обновить баланс
				}

				// Добавить транзакцию в бд
				// Добавить в кэш как отработанную
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
