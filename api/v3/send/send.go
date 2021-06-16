package send

import (
	"net/http"

	"github.com/labstack/echo"
)

func GetSend() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		id := c.Param("id")
		return c.String(http.StatusAccepted, id)
	}
}

func PostSend() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		id := c.Param("id")
		return c.String(http.StatusAccepted, id)
	}
}
