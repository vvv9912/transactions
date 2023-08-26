package server

import (
	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"transaction/internal/handler"
	"transaction/internal/kafka"
	"transaction/internal/mw"
)

type Server struct {
	echo  *echo.Echo
	Cache *cache.Cache
}

func NewServer() *Server {
	s := &Server{}
	s.echo = echo.New()
	m := mw.MW{}

	c := kafka.NewConsumer()

	c.Cache = s.Cache
	m.Cache = s.Cache

	s.echo.POST("/add", handler.HandlerAdd, m.MwAdd)
	s.echo.GET("/sub", handler.HandlerSub, m.MwSub)
	s.echo.POST("/status", handler.HandlerStatus)
	// Создам Produce и Consumer
	c.ConsumerStart()
	defer c.C.Close() //?todo
	p := kafka.NewProducer()
	defer p.P.Close() //todo
	m.P = p.P
	return s
}
func (s *Server) ServerStart(addr string) error {
	var (
	//h Handler
	)

	err := s.echo.Start(addr)
	if err != nil {
		logrus.WithFields(logrus.Fields{"func": "ServerStart"}).Fatalf("Server star error: %v", err)
	}

	return nil
}
