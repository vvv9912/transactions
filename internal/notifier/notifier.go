package notifier

import (
	"context"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"time"
	"transaction/internal/kafka"
	"transaction/internal/model"
)

// var _ Notifier = (*notifier)(nil)
//
//	type KafkaProducer interface {
//		Produce(topic string, Value []byte, Key []byte) error
//		Flush(timeoutMs int) int
//	}
type TransactionStorager interface {
	//GetStatus(ctx context.Context, numTransaction string) (int, error)
	GetNotifyTrans(ctx context.Context) ([]model.Transactions, error)
}

type notifier struct {
	KafkaProduce *kafka.Producer
	DbTrans      TransactionStorager
	Cache        *cache.Cache
}

func NewNotifier(DbTrans TransactionStorager) *notifier {
	//Один продюссер
	n := &notifier{
		DbTrans: DbTrans,
		Cache:   nil,
	}
	producer := kafka.NewProducer()
	n.KafkaProduce = producer
	n.KafkaProduce.Topic = "Transaction"
	return n

}
func (n *notifier) NotifyPending(ctx context.Context) error {
	//тут считываем из бд и отправляем
	trans, err := n.DbTrans.GetNotifyTrans(ctx)
	if err != nil {
		logrus.WithFields(
			logrus.Fields{

				"package": "notifier",
				"func":    "NotifyPending",
				"method":  "GetNotifyTrans",
			}).Fatalln(err)
		return err
	}
	_ = trans
	return nil
}

func (n *notifier) SendNotification(Value []byte, Key []byte) error {
	n.sendNotifier(Value, Key)
	return nil
}

func (n *notifier) sendNotifier(Value []byte, Key []byte) error {
	err := n.KafkaProduce.Produce(Value, Key)
	if err != nil {
		logrus.WithFields(
			logrus.Fields{

				"package": "notifier",
				"func":    "SendNotification",
				"method":  "Produce",
			}).Fatalln(err)
		return err
	}
	return nil
}

func (n *notifier) StartNotifyCron(ctx context.Context) error {
	go func() {
		err := n.NotifyPending(ctx)
		if err != nil {
			logrus.WithFields(
				logrus.Fields{

					"package": "notifier",
					"func":    "GetNotifyCron",
					"method":  "NotifyPending",
				}).Fatalln(err)
			return
		}
		select {
		case <-ctx.Done():
			n.KafkaProduce.Close()
			return
		default:

		}
		time.Sleep(time.Minute)

	}()
	return nil
}
