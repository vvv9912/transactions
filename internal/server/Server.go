package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"transaction/internal/handler"
	"transaction/internal/kafka"
	"transaction/internal/mw"
)

type Server struct {
	UsersStorage kafka.UsersStorager
	KafkaConsume mw.KafkaConsumer
	echo         *echo.Echo
	Cache        *cache.Cache
}

func NewServer(Dbusers kafka.UsersStorager, Cache *cache.Cache, KafkaProduce mw.KafkaProducer, KafkaConsumer mw.KafkaConsumer) *Server {

	s := &Server{KafkaConsume: KafkaConsumer, Cache: Cache}
	s.echo = echo.New()
	m := mw.MW{Dbusers: Dbusers, KafkaProduce: KafkaProduce, KafkaConsume: KafkaConsumer, Cache: Cache}

	//Создам вне сервера а передам сюда только интерфейс

	s.echo.POST("/add", handler.HandlerAdd, m.Mw, m.MwAdd) //Общий MW с SUB и внутри еще MW с кафкой и прочим
	s.echo.POST("/sub", handler.HandlerSub, m.Mw, m.MwSub)
	s.echo.GET("/id", handler.HandlerID)
	s.echo.POST("/status", handler.HandlerStatus)
	// Создам Produce и Consumer
	//ctx := context.TODO()
	//consum.ConsumerStart(ctx)
	//defer consum.C.Close() //?todo
	//p := kafka.NewProducer()
	//defer p.P.Close() //todo
	//m.P = p.P
	return s
}
func (s *Server) ServerStart(ctx context.Context, addr string) error {

	err := s.KafkaConsume.ConsumerStart(ctx)
	if err != nil {
		logrus.WithFields(logrus.Fields{"func": "ServerStart"}).Fatalf("Consumer start error: %v", err)
	}

	err = s.echo.Start(addr)
	if err != nil {
		logrus.WithFields(logrus.Fields{"func": "ServerStart"}).Fatalf("Server star error: %v", err)
	}

	return nil
}
