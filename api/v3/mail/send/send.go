package send

import (
	"net/http"

	"github.com/labstack/echo"
	model "github.com/yKanazawa/sendgrid-dev/model/v3/mail"
)

type ErrorResponse struct {
	Errors []struct {
		Message string      `json:"message"`
		Field   interface{} `json:"field"`
		Help    interface{} `json:"help"`
	} `json:"errors"`
}

func GetSend() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		return c.JSON(http.StatusMethodNotAllowed, model.GetErrorResponse("POST method allowed only", nil, nil))
	}
}

func PostSend() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		var postRequest model.PostRequest;
		if err := postRequest.SetPostRequest(c.Request().Body); err != nil {
			return c.JSON(http.StatusBadRequest, model.GetErrorResponse("Bad Request", nil, nil))
		}

		statusCode, errorResponse := postRequest.Validate();
		if (statusCode != http.StatusAccepted) {
			return c.JSON(statusCode, errorResponse);
		}

		return c.String(http.StatusAccepted, "")
	}
}
