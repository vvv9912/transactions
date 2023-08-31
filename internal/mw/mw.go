package mw

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"strconv"
	"transaction/internal/kafka"
	"transaction/internal/model"
)

const (
	StateAdddb string = "Add"
	StateSubdb        = "Sub"
)

type TransactionStorager interface {
	GetStatus(ctx context.Context, numTransaction string) (int, error)
}

type KafkaConsumer interface {
	ConsumerStart(ctx context.Context) error
}

type MW struct {
	Dbusers      kafka.UsersStorager
	DbTrans      TransactionStorager
	KafkaConsume KafkaConsumer
	Cache        *cache.Cache
}

func (M *MW) MwStatus(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id := ctx.QueryParam("id")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Uncorrected id")
		}
		cach, ok := M.Cache.Get(id)
		var statusTrans int
		if !ok {
			statusTrans1, err := M.DbTrans.GetStatus(context.TODO(), id)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Uncorrected id")
			}
			statusTrans = statusTrans1
		} else {
			msg := cach.(model.Caches)
			statusTrans = int(msg.Status)
		}
		ctx.Set("status", statusTrans)
		err := next(ctx)
		if err != nil {
			return err
		}
		return nil
	}
}
func (M *MW) Mw(next echo.HandlerFunc) echo.HandlerFunc {

	return func(ctx echo.Context) error {
		if len(ctx.Request().Header.Get("id")) == 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "Uncorrected id")
		}
		if len(ctx.Request().Header.Get("account")) == 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "Uncorrected account")
		}
		idTransaction := uuid.New()

		ctx.Set("id_transaction", idTransaction)
		err := next(ctx)
		if err != nil {
			return err
		}
		return nil
	}
}
func (M *MW) MwAdd(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		next(ctx)

		ID, err := strconv.Atoi(ctx.Request().Header.Get("id"))
		log.Print("прошло ")
		if err != nil {
			logrus.WithFields(logrus.Fields{"func": "MwAdd"}).Fatalf("Convert ID string to int: %v", err)
			return err
		}
		//account, err := strconv.ParseFloat(ctx.Request().Header.Get("account"), 64)
		if err != nil {
			logrus.WithFields(logrus.Fields{"func": "MwAdd"}).Fatalf("Convert account string to float: %v", err)
			return err
		}
		idTransaction := ctx.Get("id_transaction").(uuid.UUID)

		msgCacheNew := model.Caches{}
		msgCacheNew.NumberTransaction = idTransaction.String()
		msgCacheNew.Status = 0
		msgCacheNew.ID = int64(ID)

		err = M.Cache.Add(idTransaction.String(), msgCacheNew, cache.DefaultExpiration)
		if err != nil {
			logrus.WithFields(logrus.Fields{"func": "MwAdd"}).Fatalf("Add to cache: %v", err)
			return err
		}

		// Отправляем в бд транзакции
		return nil
	}
}
func (M *MW) MwSub(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {

		err := next(ctx)
		if err != nil {
			return err
		}

		ID, err := strconv.Atoi(ctx.Request().Header.Get("id"))
		log.Print("прошло ")
		if err != nil {
			logrus.WithFields(logrus.Fields{"func": "MwSub"}).Fatalf("Convert ID string to int: %v", err)
			return err
		}
		//account, err := strconv.ParseFloat(ctx.Request().Header.Get("account"), 64)
		if err != nil {
			logrus.WithFields(logrus.Fields{"func": "MwSub"}).Fatalf("Convert account string to float: %v", err)
			return err
		}
		idTransaction := ctx.Get("id_transaction").(uuid.UUID)

		msgCacheNew := model.Caches{}
		msgCacheNew.NumberTransaction = idTransaction.String()
		msgCacheNew.Status = 0
		msgCacheNew.ID = int64(ID)

		err = M.Cache.Add(idTransaction.String(), msgCacheNew, cache.DefaultExpiration)
		if err != nil {
			logrus.WithFields(logrus.Fields{"func": "MwSub"}).Fatalf("Sub to cache: %v", err)
			return err
		}

		// Отправляем в бд транзакций
		return nil

	}
}

//func (M *MW) MwAdd(next echo.HandlerFunc) echo.HandlerFunc {
//	return func(ctx echo.Context) error {
//		next(ctx)
//
//		ID, err := strconv.Atoi(ctx.Request().Header.Get("id"))
//		log.Print("прошло ")
//		if err != nil {
//			logrus.WithFields(logrus.Fields{"func": "MwAdd"}).Fatalf("Convert ID string to int: %v", err)
//			return err
//		}
//		account, err := strconv.ParseFloat(ctx.Request().Header.Get("account"), 64)
//		if err != nil {
//			logrus.WithFields(logrus.Fields{"func": "MwAdd"}).Fatalf("Convert account string to float: %v", err)
//			return err
//		}
//		idTransaction := ctx.Get("id_transaction").(uuid.UUID)
//		var message proto.Message
//		message.Numtrans = idTransaction.String()
//		message.Id = int64(ID)
//		message.Account = account
//
//		msgCacheNew := model.Caches{}
//		msgCacheNew.NumberTransaction = idTransaction.String()
//		msgCacheNew.Status = 0
//		msgCacheNew.ID = int64(ID)
//
//		err = M.Cache.Add(idTransaction.String(), msgCacheNew, cache.DefaultExpiration)
//		if err != nil {
//			logrus.WithFields(logrus.Fields{"func": "MwAdd"}).Fatalf("Add to cache: %v", err)
//			return err
//		}
//
//		go func() {
//
//			msg, err := proto2.Marshal(&message)
//			if err != nil {
//				logrus.WithFields(
//					logrus.Fields{
//						"package": "server",
//						"func":    "addHttpAnswer",
//						"method":  "Marshal",
//					}).Warningln(err)
//				return
//			}
//
//			topic := StateAdddb
//			err = M.KafkaProduce.Produce(topic, []byte(msg), []byte(idTransaction.String()))
//			if err != nil {
//				logrus.WithFields(logrus.Fields{"func": "MwAdd"}).Fatalf("Transport to produce: %v", err)
//				return
//			}
//
//			M.KafkaProduce.Flush(1 * 1000)
//			log.Print("send ")
//		}()
//
//		return nil
//	}
//}
//func (M *MW) MwSub(next echo.HandlerFunc) echo.HandlerFunc {
//	return func(ctx echo.Context) error {
//
//		err := next(ctx)
//		if err != nil {
//			return err
//		}
//
//		ID, err := strconv.Atoi(ctx.Request().Header.Get("id"))
//		log.Print("прошло ")
//		if err != nil {
//			logrus.WithFields(logrus.Fields{"func": "MwSub"}).Fatalf("Convert ID string to int: %v", err)
//			return err
//		}
//		account, err := strconv.ParseFloat(ctx.Request().Header.Get("account"), 64)
//		if err != nil {
//			logrus.WithFields(logrus.Fields{"func": "MwSub"}).Fatalf("Convert account string to float: %v", err)
//			return err
//		}
//		idTransaction := ctx.Get("id_transaction").(uuid.UUID)
//		var message proto.Message
//		message.Numtrans = idTransaction.String()
//		message.Id = int64(ID)
//		message.Account = account
//
//		msgCacheNew := model.Caches{}
//		msgCacheNew.NumberTransaction = idTransaction.String()
//		msgCacheNew.Status = 0
//		msgCacheNew.ID = int64(ID)
//
//		err = M.Cache.Add(idTransaction.String(), msgCacheNew, cache.DefaultExpiration)
//		if err != nil {
//			logrus.WithFields(logrus.Fields{"func": "MwSub"}).Fatalf("Sub to cache: %v", err)
//			return err
//		}
//
//		go func() {
//
//			msg, err := proto2.Marshal(&message)
//			if err != nil {
//				logrus.WithFields(
//					logrus.Fields{
//						"package": "server",
//						"func":    "SubHttpAnswer",
//						"method":  "Marshal",
//					}).Warningln(err)
//				return
//			}
//
//			topic := StateSubdb
//			err = M.KafkaProduce.Produce(topic, []byte(msg), []byte(idTransaction.String()))
//			if err != nil {
//				logrus.WithFields(logrus.Fields{"func": "MwSub"}).Fatalf("Transport to produce: %v", err)
//				return
//			}
//
//			M.KafkaProduce.Flush(1 * 1000)
//			log.Print("send ")
//		}()
//		return nil
//
//	}
//}
