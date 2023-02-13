package handler

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/habibitcoin/habibtaro/configs"
	"github.com/habibitcoin/habibtaro/taro"
	"github.com/labstack/echo"
)

func (h *Handler) FetchAsset(c echo.Context) (err error) {
	var (
		ctx = context.Background()
		a   = &taro.TaroAddressRequest{}
	)

	ctx, _ = configs.LoadConfig(ctx)
	taroClient := taro.NewClient(ctx)

	c.Bind(a)

	resp, err := taroClient.GetAsset(c.Param("name"))
	if err != nil {
		log.Println("Error creating address")
		log.Println(err)
	}

	str, _ := json.MarshalIndent(resp, "", "\t")

	str1 := strings.Replace(string(str), "\n", "<br>", -1)
	str1 = strings.Replace(string(str1), "\t", "&emsp;", -1)

	return c.Render(http.StatusOK, "asset.html", map[string]interface{}{
		"asset": template.HTML(str1),
	})
}
