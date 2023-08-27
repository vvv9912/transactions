package main

import (
	"context"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"transaction/internal/kafka"
	"transaction/internal/server"
	"transaction/internal/storage"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
	"os/signal"
	"syscall"
	"transaction/internal/config"
)

func main() {

	db, err := sqlx.Connect("postgres", config.Get().DatabaseDSN)
	if err != nil {
		logrus.WithFields(logrus.Fields{"func": "main"}).Fatalf("faild to connetct to database: %v", err)
		return
	}
	defer db.Close()

	var (
		usersStorage       = storage.NewUsersStorage(db) //подкл бд
		transactionStorang = storage.NewTransStorage(db)
		cache              = cache.New(cache.DefaultExpiration, 0)
		consumer           = kafka.NewConsumer(cache, usersStorage, transactionStorang)
		producer           = kafka.NewProducer()
	)

	defer consumer.C.Close()
	defer producer.P.Close()

	logrus.Infof("database_dsn: %v", config.Get().DatabaseDSN)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	s := server.NewServer(usersStorage, transactionStorang, cache, producer, consumer)
	s.ServerStart(ctx, config.Get().HTTPServer)

}
