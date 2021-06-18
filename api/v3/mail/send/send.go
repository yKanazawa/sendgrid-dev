package send

import (
	"encoding/json"
	"fmt"
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

		if err := json.NewDecoder(c.Request().Body).Decode(&postRequest); err != nil {
			return c.JSON(http.StatusBadRequest, getErrorResponse("Bad Request", nil, nil))
		}

		// Debug
		fmt.Printf("postRequest %#v\n", postRequest)
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