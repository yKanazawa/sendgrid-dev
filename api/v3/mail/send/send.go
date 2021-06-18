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
		return c.JSON(http.StatusMethodNotAllowed, getErrorResponse("POST method allowed only", nil, nil))
	}
}

func PostSend() echo.HandlerFunc {
	return func(c echo.Context) (err error) {

		var postRequest model.PostRequest
		if err := postRequest.SetPostRequest(c.Request().Body); err != nil {
			return c.JSON(http.StatusBadRequest, getErrorResponse("Bad Request", nil, nil))
		}

		return c.String(http.StatusAccepted, "")
	}
}

func getErrorResponse (message string, field interface{}, help interface{}) ErrorResponse {
	errorJSON := ErrorResponse{}
    e := struct {
		Message string      `json:"message"`
		Field   interface{} `json:"field"`
		Help    interface{} `json:"help"`
    } {
        message,
		nil,
		nil,
    }
	errorJSON.Errors = append(errorJSON.Errors, e)

	return errorJSON
}
