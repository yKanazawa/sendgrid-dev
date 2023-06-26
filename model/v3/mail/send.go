package send

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jordan-wright/email"
	"gopkg.in/go-playground/validator.v9"
)

type PostRequest struct {
	Personalizations []struct {
		Subject string `json:"subject"`
		To      []struct {
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
	} `json:"personalizations" validate:"required"`
	From struct {
		Email string `json:"email" validate:"required"`
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
	} `json:"content" validate:"required"`
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
	return json.NewDecoder(requestBody).Decode(&postRequest)
}

func (postRequest *PostRequest) Validate() (int, ErrorResponse) {
	validate := validator.New()
	if err := validate.Struct(postRequest); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.ActualTag() {
			case "required":
				switch err.StructField() {
				case "Personalizations":
					return http.StatusBadRequest,
						GetErrorResponse(
							"The personalizations field is required and must have at least one personalization.",
							"personalizations",
							"http://sendgrid.com/docs/API_Reference/Web_API_v3/Mail/errors.html#-Personalizations-Errors",
						)
				case "Email":
					return http.StatusBadRequest,
						GetErrorResponse(
							"The from object must be provided for every email send. It is an object that requires the email parameter, but may also contain a name parameter.  e.g. {\"email\" : \"example@example.com\"}  or {\"email\" : \"example@example.com\", \"name\" : \"Example Recipient\"}.",
							"from.email",
							"http://sendgrid.com/docs/API_Reference/Web_API_v3/Mail/errors.html#message.from",
						)
				case "Content":
					return http.StatusBadRequest,
						GetErrorResponse(
							"Unless a valid template_id is provided, the content parameter is required. There must be at least one defined content block. We typically suggest both text/plain and text/html blocks are included, but only one block is required.",
							"content",
							"http://sendgrid.com/docs/API_Reference/Web_API_v3/Mail/errors.html#message.content",
						)
				}
			}
		}
	}

	return sendMailWithSMTP(*postRequest)
}

func GetErrorResponse(message string, field interface{}, help interface{}) ErrorResponse {
	errorJSON := ErrorResponse{}
	e := struct {
		Message string      `json:"message"`
		Field   interface{} `json:"field"`
		Help    interface{} `json:"help"`
	}{
		message,
		field,
		help,
	}
	errorJSON.Errors = append(errorJSON.Errors, e)

	return errorJSON
}

// Send mail with SMTP
func sendMailWithSMTP(postRequest PostRequest) (int, ErrorResponse) {
	for _, personalizations := range postRequest.Personalizations {
		e := email.NewEmail()

		e.From = postRequest.From.Name + " <" + postRequest.From.Email + ">"

		for _, to := range personalizations.To {
			e.To = append(e.To, getEmailwithName(to))
		}

		for _, cc := range personalizations.Cc {
			e.Cc = append(e.Cc, getEmailwithName(cc))
		}

		for _, bcc := range personalizations.Bcc {
			e.Bcc = append(e.Bcc, getEmailwithName(bcc))
		}

		e.Subject = personalizations.Subject

		for _, content := range postRequest.Content {
			if content.Type == "text/html" {
				e.HTML = []byte(content.Value)
			} else {
				e.Text = []byte(content.Value)
			}
		}

		i := 0
		for _, attachment := range postRequest.Attachments {
			dirName := createAttachment(attachment.Filename, attachment.Content, i)
			if dirName == "" {
				return http.StatusBadRequest,
					GetErrorResponse(
						"The attachment content must be base64 encoded.",
						"attachments."+strconv.Itoa(i)+".content",
						"http://sendgrid.com/docs/API_Reference/Web_API_v3/Mail/errors.html#message.attachments.content",
					)
			}
			e.AttachFile(filepath.Join(dirName, attachment.Filename))
			i++
		}

		if os.Getenv("SENDGRID_DEV_TEST") == "1" {
			continue
		}

		if len(os.Getenv("SENDGRID_DEV_SMTP_USERNAME")) > 0 {
			arr := strings.Split(os.Getenv("SENDGRID_DEV_SMTP_SERVER"), ":")
			e.Send(
				os.Getenv("SENDGRID_DEV_SMTP_SERVER"),
				smtp.PlainAuth(
					"",
					os.Getenv("SENDGRID_DEV_SMTP_USERNAME"),
					os.Getenv("SENDGRID_DEV_SMTP_PASSWORD"),
					arr[0],
				),
			)
		}

		e.Send(os.Getenv("SENDGRID_DEV_SMTP_SERVER"), nil)
	}
	return http.StatusAccepted, GetErrorResponse("", nil, nil)
}

// Get "Name <name@example.com>"
func getEmailwithName(t struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}) string {
	return t.Name + " <" + t.Email + ">"
}

// Create attachment from base64 string
func createAttachment(fileName string, base64Content string, i int) string {
	data, err := base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		return ""
	}

	dirName := filepath.Join(os.TempDir(), "attachment_"+strconv.FormatInt(time.Now().UnixNano(), 10))
	os.Mkdir(dirName, 0777)
	file, err := os.Create(filepath.Join(dirName, fileName))
	if err != nil {
		fmt.Println("Create file failed.", fileName)
		return ""
	}

	defer file.Close()
	file.Write(data)

	return dirName
}
