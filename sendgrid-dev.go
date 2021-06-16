package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jordan-wright/email"
)

type Messages struct {
	Errors []Error `json:"errors"`
}

type Error struct {
	Message interface{} `json:"message"`
	Field   interface{} `json:"field"`
	Help    interface{} `json:"help"`
}

type SendPostData struct {
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

//
// main
//
func main() {
	fmt.Println("SENDGRID_DEV_PORT", os.Getenv("SENDGRID_DEV_PORT"))
	fmt.Println("SENDGRID_DEV_APIKEY", os.Getenv("SENDGRID_DEV_APIKEY"))
	fmt.Println("SENDGRID_SERVER_PORT", os.Getenv("SENDGRID_SMTP_SERVER"))
	fmt.Println("SENDGRID_SMTP_PORT", os.Getenv("SENDGRID_SMTP_PORT"))

	// Initiate Router
	router := mux.NewRouter()

	// Route Hnadlers / Endpoints
	router.HandleFunc("/v3/mail/send", send).Methods("GET", "POST")

	port, err := strconv.Atoi(os.Getenv("SENDGRID_DEV_PORT"))
	if err != nil || port < 0 || port > 65535 {
		port = 8000
	}
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), router))
}

//
// /v3/mail/send
//
func send(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	if request.Method != "POST" {
		errorMessage(response, http.StatusMethodNotAllowed, "POST method allowed only", nil, nil)
		return
	}

	if request.Header.Get("Authorization") != "Bearer "+os.Getenv("SENDGRID_DEV_APIKEY") {
		errorMessage(response, http.StatusUnauthorized, "The provided authorization grant is invalid, expired, or revoked", nil, nil)
		return
	}

	if request.Header.Get("Content-Type") != "application/json" {
		errorMessage(response, http.StatusUnsupportedMediaType, "Content-Type should be application/json.", nil, nil)
		return
	}

	// Debug
	fmt.Println("request.Method", request.Method)
	fmt.Println("request.Header", request.Header)
	fmt.Println("request.Body", request.Body)
	fmt.Println()

	var sendPostData SendPostData

	if err := json.NewDecoder(request.Body).Decode(&sendPostData); err != nil {
		errorMessage(response, http.StatusBadRequest, "Bad Request", nil, nil)
		return
	}

	if validateSendPostData(response, sendPostData) == false {
		return
	}

	// Debug
	fmt.Printf("sendPostData %#v\n", sendPostData)
	fmt.Printf("sendPostData.Personalizations %#v\n", sendPostData.Personalizations)
	fmt.Printf("sendPostData.Personalizations[0].To[0].Email %#v\n", sendPostData.Personalizations[0].To[0].Email)
	fmt.Printf("sendPostData.From.Email %#v\n", sendPostData.From.Email)
	fmt.Printf("sendPostData.Subject %#v\n", sendPostData.Subject)
	fmt.Printf("sendPostData.Content[0].Type %#v\n", sendPostData.Content[0].Type)
	fmt.Printf("sendPostData.Content[0].Value %#v\n", sendPostData.Content[0].Value)
	fmt.Println()

	sendMailWithSMTP(sendPostData, response)

	response.WriteHeader(http.StatusAccepted)
}

//
// Return error message as JSON
//
func errorMessage(response http.ResponseWriter, statusCode int, message string, field interface{}, help interface{}) {
	response.WriteHeader(statusCode)

	messages := Messages{}
	error := Error{}
	error.Message = message
	error.Field = field
	error.Help = help

	messages.Errors = append(messages.Errors, error)
	json.NewEncoder(response).Encode(messages)
}

//
// Validate POST data (/v3/mail/send)
//
func validateSendPostData(response http.ResponseWriter, sendPostData SendPostData) bool {
	if len(sendPostData.Personalizations) == 0 {
		errorMessage(
			response,
			http.StatusBadRequest,
			"The personalizations field is required and must have at least one personalization.",
			"personalizations",
			"http://sendgrid.com/docs/API_Reference/Web_API_v3/Mail/errors.html#-Personalizations-Errors",
		)
		return false
	}

	if sendPostData.From.Email == "" {
		errorMessage(
			response,
			http.StatusBadRequest,
			"The from object must be provided for every email send. It is an object that requires the email parameter, but may also contain a name parameter.  e.g. {\"email\" : \"example@example.com\"}  or {\"email\" : \"example@example.com\", \"name\" : \"Example Recipient\"}.",
			"from.email",
			"http://sendgrid.com/docs/API_Reference/Web_API_v3/Mail/errors.html#message.from",
		)
		return false
	}

	if sendPostData.Subject == "" {
		errorMessage(
			response,
			http.StatusBadRequest,
			"The subject is required. You can get around this requirement if you use a template with a subject defined or if every personalization has a subject defined.",
			"subject",
			"http://sendgrid.com/docs/API_Reference/Web_API_v3/Mail/errors.html#message.subject",
		)
		return false
	}

	if len(sendPostData.Content) == 0 {
		errorMessage(
			response,
			http.StatusBadRequest,
			"Unless a valid template_id is provided, the content parameter is required. There must be at least one defined content block. We typically suggest both text/plain and text/html blocks are included, but only one block is required.",
			"content",
			"http://sendgrid.com/docs/API_Reference/Web_API_v3/Mail/errors.html#message.content",
		)
		return false
	}

	typeCount := map[string]int{}
	for _, content := range sendPostData.Content {
		typeCount[content.Type]++
		if typeCount[content.Type] > 1 {
			errorMessage(
				response,
				http.StatusBadRequest,
				"If present, text/plain and text/html may only be provided once.",
				"content",
				nil,
			)
			return false
		}
	}

	return true
}

//
// Send mail with SMTP
//
func sendMailWithSMTP(sendPostData SendPostData, response http.ResponseWriter) {
	smtpServer := os.Getenv("SENDGRID_DEV_SMTP_SERVER")
	if smtpServer == "" {
		smtpServer = "localhost"
	}

	smtpPort, err := strconv.Atoi(os.Getenv("SENDGRID_DEV_SMTP_PORT"))
	if err != nil || smtpPort < 0 || smtpPort > 65535 {
		smtpPort = 1025
	}

	for _, personalizations := range sendPostData.Personalizations {
		e := email.NewEmail()

		e.From = sendPostData.From.Email

		for _, to := range personalizations.To {
			e.To = append(e.To, getEmailwithName(to))
		}

		for _, cc := range personalizations.Cc {
			e.Cc = append(e.Cc, getEmailwithName(cc))
		}

		for _, bcc := range personalizations.Bcc {
			e.Bcc = append(e.Bcc, getEmailwithName(bcc))
		}

		e.Subject = sendPostData.Subject

		for _, content := range sendPostData.Content {
			if content.Type == "text/html" {
				e.HTML = []byte(content.Value)
			} else {
				e.Text = []byte(content.Value)
			}
		}

		i := 0
		for _, attachment := range sendPostData.Attachments {
			// Debug
			fmt.Println("attachment.Filename", attachment.Filename)
			dirName := createAttachment(attachment.Filename, attachment.Content, i, response)
			e.AttachFile(filepath.Join(dirName, attachment.Filename))
			i++
		}

		if os.Getenv("SENDGRID_DEV_TEST") == "1" {
			continue
		}

		e.Send(smtpServer+":"+strconv.Itoa(smtpPort), nil)
	}
}

//
// Get "Name <name@example.com>"
//
func getEmailwithName(t struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}) string {
	return t.Name + " <" + t.Email + ">"
}

//
// Create attachment from base64 string
//
func createAttachment(fileName string, base64Content string, i int, response http.ResponseWriter) string {
	data, err := base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		errorMessage(
			response,
			http.StatusBadRequest,
			"The attachment content must be base64 encoded.",
			"attachments."+strconv.Itoa(i)+".content",
			"http://sendgrid.com/docs/API_Reference/Web_API_v3/Mail/errors.html#message.attachments.content",
		)
		return ""
	}

	dirName := filepath.Join(os.TempDir(), "attachment_"+strconv.FormatInt(time.Now().UnixNano(), 10))
	// Debug
	fmt.Println("dirName", dirName)
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
