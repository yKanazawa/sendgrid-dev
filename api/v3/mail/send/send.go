package send

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type PostRequest struct {
	Personalizations []struct {
		To []struct {
			Email string `json:"email"`
			Name  string `json:"name"`
		} `json:"to"`
		Cc []struct {
			Email string `json:"email"`
			Name  string `json:"name"`
		} `json:"cc"`
		Bcc []struct {
			Email string `json:"email"`
			Name  string `json:"name"`
		} `json:"bcc"`
	} `json:"personalizations"`
	From struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"from"`
	ReplyTo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"reply_to"`
	Subject string `json:"subject"`
	Content []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"content"`
	Attachments []struct {
		Content     string `json:"content"`
		Type        string `json:"type"`
		Filename    string `json:"filename"`
		Disposition string `json:"disposition"`
		ContentId   string `json:"content_id"`
	} `json:"attachments"`
}

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

		var postRequest PostRequest

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