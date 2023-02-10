package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/habibitcoin/habibtaro/configs"
	"github.com/habibitcoin/habibtaro/taro"
	"github.com/labstack/echo"
)

func (h *Handler) Address(c echo.Context) (err error) {
	var (
		ctx = context.Background()
		a   = &taro.TaroAddressRequest{}
	)

	ctx, _ = configs.LoadConfig(ctx)
	taroClient := taro.NewClient(ctx)

	c.Bind(a)

	resp, err := taroClient.CreateAddress(a.BootstrapInfo, a.GroupKey, a.Amount)
	if err != nil {
		log.Println("Error creating address")
		log.Println(err)
	}

	log.Println(resp.Address)

	return c.Render(http.StatusOK, "address.html", map[string]interface{}{
		"address": resp,
	})
}

func (h *Handler) ViewAddress(c echo.Context) (err error) {
	return c.Render(http.StatusOK, "address.html", map[string]interface{}{})
}
