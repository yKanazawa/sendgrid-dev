package send

import (
	"net/http"
	"os"

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
		Authorization := c.Request().Header["Authorization"]
		if len(Authorization) == 0 || Authorization[0] != "Bearer "+os.Getenv("SENDGRID_DEV_APIKEY") {
			return c.JSON(http.StatusUnsupportedMediaType, model.GetErrorResponse("The provided authorization grant is invalid, expired, or revoked", nil, nil))
		}

		contentType := c.Request().Header["Content-Type"]
		if len(contentType) == 0 || contentType[0] != "application/json" {
			return c.JSON(http.StatusUnsupportedMediaType, model.GetErrorResponse("Content-Type should be application/json", nil, nil))
		}

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
