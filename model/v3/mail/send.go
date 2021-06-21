package send

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io"

	"gopkg.in/go-playground/validator.v9"
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
	Subject string `json:"subject" validate:"required"`
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

func (postRequest *PostRequest) SetPostRequest(requestBody io.ReadCloser) error {
	return json.NewDecoder(requestBody).Decode(&postRequest);
}

func (postRequest *PostRequest) Validate() (int, ErrorResponse) {
	validate := validator.New()
	if err := validate.Struct(postRequest); err != nil {
		for _, err := range err.(validator.ValidationErrors) {

			// Debug
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
				return http.StatusBadRequest, GetErrorResponse("The subject is required. You can get around this requirement if you use a template with a subject defined or if every personalization has a subject defined.", "subject", "http://sendgrid.com/docs/API_Reference/Web_API_v3/Mail/errors.html#message.subject")
			}
		}
	}
	return http.StatusAccepted, GetErrorResponse("", nil, nil)
}

func GetErrorResponse (message string, field interface{}, help interface{}) ErrorResponse {
	errorJSON := ErrorResponse{}
    e := struct {
		Message string      `json:"message"`
		Field   interface{} `json:"field"`
		Help    interface{} `json:"help"`
    } {
        message,
		field,
		help,
    }
	errorJSON.Errors = append(errorJSON.Errors, e)

	return errorJSON
}