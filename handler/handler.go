package handler

import (
	"net/http"

	"github.com/habibitcoin/habibtaro/model"
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
	var ()

	// Retrieve featured (first 4) users from database
	users := []*model.PublicUser{}
	posts := []*model.Post{}

	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"users": users,
		"posts": posts,
	})
}
