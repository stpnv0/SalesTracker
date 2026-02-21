package handler

import (
	"net/http"

	"github.com/wb-go/wbf/ginext"
)

type errorResponse struct {
	Error string `json:"error"`
}

func respondError(c *ginext.Context, code int, msg string) {
	c.JSON(code, errorResponse{Error: msg})
}

func respondJSON(c *ginext.Context, code int, data interface{}) {
	c.JSON(code, data)
}

func respondNoContent(c *ginext.Context) {
	c.Status(http.StatusNoContent)
}
