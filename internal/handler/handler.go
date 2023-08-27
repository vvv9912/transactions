package handler

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

func HandlerStatus(ctx echo.Context) error {
	//err := ctx.String(http.StatusOK, "test 200")
	status := ctx.Get("status").(int)
	err := ctx.String(http.StatusOK, fmt.Sprintf("Status %d", status))
	if err != nil {
		return err
	}
	return nil
}
func HandlerAdd(ctx echo.Context) error {
	id := ctx.Request().Header.Get("id")
	arr := ctx.Request().Header.Get("account")
	idTransaction := ctx.Get("id_transaction").(uuid.UUID)
	fmt.Println(idTransaction.String())

	s := fmt.Sprintf("id: %s\naccount: %s\nnum. transaction: %s", id, arr, idTransaction.String())
	err := ctx.String(http.StatusOK, s)
	if err != nil {
		return err
	}
	return nil
}
func HandlerSub(ctx echo.Context) error {
	idTransaction := ctx.Get("id_transaction").(uuid.UUID)
	s := fmt.Sprintf("%d\n%s/status?id=%s", http.StatusOK, ctx.Request().Host, idTransaction.String())
	err := ctx.String(http.StatusOK, s)
	if err != nil {
		return err
	}
	return nil
}
