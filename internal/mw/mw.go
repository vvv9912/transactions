package mw

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
	"log"
	"net/http"
	"transaction/internal/model"
)

type UsersStorage interface {
	AddUsers(ctx context.Context, Users model.Users) (int64, error)
}
type MW struct {
	dbusers UsersStorage
	P       *kafka.Producer
	Cache   *cache.Cache
}

func (M *MW) MwAdd(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		//return echo.NewHTTPError(300, "failde")
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
		log.Print("Executing middlewareOne again Add")
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
		log.Print("Executing middlewareOne again Sub")
		return nil

	}
}
