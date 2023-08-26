package handler

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

// type Handler struct {
// }
//
//	func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
//		_, err := io.WriteString(w, "okay")
//		if err != nil {
//			log.Printf("Error add")
//			return
//		}
//	}
func HandlerStatus(ctx echo.Context) error {
	err := ctx.String(http.StatusOK, "test 200")
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
	_ = idTransaction
	_ = id
	_ = arr
	s := fmt.Sprintf("id: %s\naccount: %s\nnum. transaction: %s", id, arr, idTransaction.String())
	err := ctx.String(http.StatusOK, s)
	if err != nil {
		return err
	}
	return nil
}
func HandlerSub(ctx echo.Context) error {
	//Вывод средств , возврат 200 и ссылка на Get запрос
	//idTransaction := ctx.Get("id_transaction").(uuid.UUID)
	idTransaction := uuid.New()
	s := fmt.Sprintf("%d\n%s/id?%s", http.StatusOK, ctx.Request().Host, idTransaction.String())
	err := ctx.String(http.StatusOK, s)
	if err != nil {
		return err
	}
	return nil
}
func HandlerID(ctx echo.Context) error {
	//
	err := ctx.String(http.StatusOK, "test sub")
	if err != nil {
		return err
	}
	return nil
}
