package mw

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	proto2 "github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"strconv"
	"transaction/internal/model"
	"transaction/internal/proto"
)

const (
	StateAdddb string = "Add"
	StateSubdb        = "Sub"
)

type UsersStorage interface {
	AddUsers(ctx context.Context, Users model.Users) (int64, error)
}

//	type Proder interface {
//		kafka.Producer
//	}
type MW struct {
	Dbusers UsersStorage
	//Produs  Proder
	P     *kafka.Producer
	Cache *cache.Cache
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

		/*
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
			//
		*/
		next(ctx)
		//Тут реализация передачи в кафку
		log.Print("Общий мв в адд повтор ")
		//

		ID, err := strconv.Atoi(ctx.Request().Header.Get("id"))
		log.Print("прошло ")
		if err != nil {
			logrus.WithFields(logrus.Fields{"func": "MwAdd"}).Fatalf("Convert ID string to int: %v", err)
			return err
		}
		account, err := strconv.ParseFloat(ctx.Request().Header.Get("account"), 64)
		if err != nil {
			logrus.WithFields(logrus.Fields{"func": "MwAdd"}).Fatalf("Convert account string to float: %v", err)
			return err
		}
		id_transaction := ctx.Get("id_transaction").(uuid.UUID)
		var message proto.Message
		message.Numtrans = id_transaction.String()
		message.Id = int64(ID)
		message.Account = account

		msgCacheNew := make([]model.Caches, 1)
		msgCacheNew[0].NumberTransaction = id_transaction.String()
		msgCacheNew[0].Status = 0
		msgCacheNew[0].ID = int64(ID)

		//map ->транзакция -> структура
		//Добавляем кэш
		M.Cache.Add(id_transaction.String(), msgCacheNew, cache.DefaultExpiration)
		//начало консюмера
		go func() {

			msg, err := proto2.Marshal(&message)
			if err != nil {
				logrus.WithFields(
					logrus.Fields{
						"package": "server",
						"func":    "addHttpAnswer",
						"method":  "Marshal",
					}).Warningln(err)
				return
			}

			topic := StateAdddb

			M.P.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          []byte(msg),
				Key:            []byte(id_transaction.String()),
			}, nil)
			M.P.Flush(1 * 1000)
			log.Print("send ")
		}()
		//M.Cache.Add()
		return nil
	}
}
func (M *MW) MwSub(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		log.Print("Executing middlewareOne Sub")
		err := next(ctx)
		if err != nil {
			return err
		}
		//Тут реализация передачи в кафку
		log.Print("Executing middlewareOne again Sub")
		return nil

	}
}
