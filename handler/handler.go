package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/habibitcoin/habibtaro/configs"
	"github.com/habibitcoin/habibtaro/model"
	"github.com/habibitcoin/habibtaro/taro"
	"github.com/labstack/echo"
)

type (
	Handler struct {
	}
)

const (
	// Key (Should come from somewhere else).
	Key = "secret"
)

func (h *Handler) Index(c echo.Context) (err error) {
	var (
		ctx = context.Background()
	)

	ctx, _ = configs.LoadConfig(ctx)
	taroClient := taro.NewClient(ctx)

	resp, err := taroClient.ListAssets() // leave 500000 cushion
	if err != nil {
		log.Println("Error opening channel")
		log.Println(err)
	}

	// Retrieve featured (first 4) users from database
	users := []*model.PublicUser{}

	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"users":  users,
		"assets": resp.TaroAssets,
	})
}
