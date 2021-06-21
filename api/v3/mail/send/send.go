package send

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	model "github.com/yKanazawa/sendgrid-dev/model/v3/mail"
	"gopkg.in/go-playground/validator.v9"
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

		validate := validator.New()
		if err := validate.Struct(postRequest); err != nil {
			for _, err := range err.(validator.ValidationErrors) {

				fmt.Println("err.Namespace()", err.Namespace())
				fmt.Println("err.Field()", err.Field())
				fmt.Println("err.StructNamespace()", err.StructNamespace())
				fmt.Println("err.StructField()", err.StructField())
				fmt.Println("err.Tag()", err.Tag())
				fmt.Println("err.ActualTag()", err.ActualTag())
				fmt.Println("err.Kind()", err.Kind())
				fmt.Println("err.Type()", err.Type())
				fmt.Println("err.Value()", err.Value())
				fmt.Println("err.Param()", err.Param())
				fmt.Println()
				switch err.ActualTag() {
				case "required":
					return c.JSON(http.StatusBadRequest, getErrorResponse("Validate Failed", nil, nil))
				}
			}
			return c.JSON(http.StatusBadRequest, getErrorResponse("Validate Failed", nil, nil))
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
