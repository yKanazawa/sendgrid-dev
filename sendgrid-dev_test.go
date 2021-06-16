package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func compare(t *testing.T, actual interface{}, expected interface{}) {
	if actual != expected {
		t.Errorf("got |%v| want |%v|", actual, expected)
	}
}

func requestTest(t *testing.T, method string, path string, headers map[string]string, body io.Reader, expectedStatus int, expectedBody interface{}) {
	request, _ := http.NewRequest(method, path, body)

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	recorder := httptest.NewRecorder()

	handler := http.HandlerFunc(send)
	handler.ServeHTTP(recorder, request)

	compare(t, recorder.Code, expectedStatus)
	compare(t, strings.TrimSpace(recorder.Body.String()), expectedBody)
}

func TestSend(t *testing.T) {
	os.Setenv("SENDGRID_DEV_TEST", "1")
	os.Setenv("SENDGRID_DEV_APIKEY", "SG.xxxxx")
	headers := make(map[string]string, 2)

	// NG (Not POST)
	requestTest(
		t,
		"GET",
		"/v3/send",
		headers,
		bytes.NewBuffer([]byte(``)),
		http.StatusMethodNotAllowed,
		`{"errors":[{"message":"POST method allowed only","field":null,"help":null}]}`,
	)

	// NG (Missing Authorization)
	requestTest(
		t,
		"POST",
		"/v3/send",
		headers,
		bytes.NewBuffer([]byte(``)),
		http.StatusUnauthorized,
		`{"errors":[{"message":"The provided authorization grant is invalid, expired, or revoked","field":null,"help":null}]}`,
	)

	headers["Authorization"] = "Bearer SG.xxxxx"

	// NG (Missing Content-Type)
	requestTest(
		t,
		"POST",
		"/v3/send",
		headers,
		bytes.NewBuffer([]byte(``)),
		http.StatusUnsupportedMediaType,
		`{"errors":[{"message":"Content-Type should be application/json.","field":null,"help":null}]}`,
	)

	headers["Content-Type"] = "text/plain"

	// NG (Content-Type is not application/json)
	requestTest(
		t,
		"POST",
		"/v3/send",
		headers,
		bytes.NewBuffer([]byte(``)),
		http.StatusUnsupportedMediaType,
		`{"errors":[{"message":"Content-Type should be application/json.","field":null,"help":null}]}`,
	)

	headers["Content-Type"] = "application/json"

	// NG (Missing PostData)
	requestTest(
		t,
		"POST",
		"/v3/send",
		headers,
		bytes.NewBuffer([]byte(``)),
		http.StatusBadRequest,
		`{"errors":[{"message":"Bad Request","field":null,"help":null}]}`,
	)

	// NG (Missing personalizations)
	requestTest(
		t,
		"POST",
		"/v3/send",
		headers,
		bytes.NewBuffer([]byte(`{
			"from": {
				"email": "from@example.com"
			},
			"subject": "Subject", 
			"content": [{
				"type": "text/plain", 
				"value": "Content"
			}]
		}`)),
		http.StatusBadRequest,
		`{"errors":[{"message":"The personalizations field is required and must have at least one personalization.","field":"personalizations","help":"http://sendgrid.com/docs/API_Reference/Web_API_v3/Mail/errors.html#-Personalizations-Errors"}]}`,
	)

	// NG (Missing from.Email)
	requestTest(
		t,
		"POST",
		"/v3/send",
		headers,
		bytes.NewBuffer([]byte(`{
			"personalizations": [{
				"to": [{
					"email": "to@example.com"
				}]
			}], 
			"subject": "Subject", 
			"content": [{
				"type": "text/plain", 
				"value": "Content"
			}]
		}`)),
		http.StatusBadRequest,
		`{"errors":[{"message":"The from object must be provided for every email send. It is an object that requires the email parameter, but may also contain a name parameter.  e.g. {\"email\" : \"example@example.com\"}  or {\"email\" : \"example@example.com\", \"name\" : \"Example Recipient\"}.","field":"from.email","help":"http://sendgrid.com/docs/API_Reference/Web_API_v3/Mail/errors.html#message.from"}]}`,
	)

	// NG (Missing subject)
	requestTest(
		t,
		"POST",
		"/v3/send",
		headers,
		bytes.NewBuffer([]byte(`{
			"personalizations": [{
				"to": [{
					"email": "to@example.com"
				}]
			}],
			"from": {
				"email": "from@example.com"
			}, 
			"content": [{
				"type": "text/plain", 
				"value": "Content"
			}]
		}`)),
		http.StatusBadRequest,
		`{"errors":[{"message":"The subject is required. You can get around this requirement if you use a template with a subject defined or if every personalization has a subject defined.","field":"subject","help":"http://sendgrid.com/docs/API_Reference/Web_API_v3/Mail/errors.html#message.subject"}]}`,
	)

	// NG (Missing content)
	requestTest(
		t,
		"POST",
		"/v3/send",
		headers,
		bytes.NewBuffer([]byte(`{
			"personalizations": [{
				"to": [{
					"email": "to@example.com"
				}]
			}], 
			"from": {
				"email": "from@example.com"
			}, 
			"subject": "Subject"
		}`)),
		http.StatusBadRequest,
		`{"errors":[{"message":"Unless a valid template_id is provided, the content parameter is required. There must be at least one defined content block. We typically suggest both text/plain and text/html blocks are included, but only one block is required.","field":"content","help":"http://sendgrid.com/docs/API_Reference/Web_API_v3/Mail/errors.html#message.content"}]}`,
	)

	// OK
	requestTest(
		t,
		"POST",
		"/v3/send",
		headers,
		bytes.NewBuffer([]byte(`{
			"personalizations": [{
				"to": [{
					"email": "to@example.com"
				}]
			}], 
			"from": {
				"email": "from@example.com"
			}, 
			"subject": "Subject", 
			"content": [{
				"type": "text/plain", 
				"value": "Content"
			}]
		}`)),
		http.StatusAccepted,
		``,
	)

	// OK (multiple to, cc, bcc and reply-to with name)
	requestTest(
		t,
		"POST",
		"/v3/send",
		headers,
		bytes.NewBuffer([]byte(`{
			"personalizations": [{
				"to": [{
					"email": "to1@example.com", 
					"name": "ToName1"
				}, {
					"email": "to2@example.com", 
					"name": "ToName2"
				}], 
				"cc": [{
					"email": "cc1@example.com", 
					"name": "CcName1"
				}, {
					"email": "cc2@example.com", 
					"name": "CcName2"
				}],
				"bcc": [{
					"email": "bcc1@example.com", 
					"name": "BccName1"
				}, {
					"email": "bcc2@example.com", 
					"name": "BccName2"
				}]
			}], 
			"from": {
				"email": "from@example.com", 
				"name": "FromName"
			}, 
			"reply_to": {
				"email": "reply_to@example.com", 
				"name": "ReplyToName"
			}, 
			"subject": "Subject", 
			"content": [{
				"type": "text/plain", 
				"value": "Content"
			}]
		}`)),
		http.StatusAccepted,
		``,
	)

	// OK (multiple personalizations)
	requestTest(
		t,
		"POST",
		"/v3/send",
		headers,
		bytes.NewBuffer([]byte(`{
			"personalizations": [
				{
					"to": [{
						"email": "to1@example.com", 
						"name": "ToName1"
					}, {
						"email": "to2@example.com", 
						"name": "ToName2"
					}]
				},
				{
					"to": [{
						"email": "to3@example.com", 
						"name": "ToName3"
					}, {
						"email": "to4@example.com", 
						"name": "ToName4"
					}]
				}
			], 
			"from": {
				"email": "from@example.com", 
				"name": "FromName"
			}, 
			"subject": "Subject", 
			"content": [{
				"type": "text/plain", 
				"value": "Content"
			}]
		}`)),
		http.StatusAccepted,
		``,
	)

	// OK (text/html)
	requestTest(
		t,
		"POST",
		"/v3/send",
		headers,
		bytes.NewBuffer([]byte(`{
			"personalizations": [{
				"to": [{
					"email": "to@example.com"
				}]
			}], 
			"from": {
				"email": "from@example.com"
			}, 
			"subject": "Subject", 
			"content": [{
				"type": "text/html", 
				"value": "<h1>Content</h1>"
			}]
		}`)),
		http.StatusAccepted,
		``,
	)

	// NG (multiple content text/plain, text/plain)
	requestTest(
		t,
		"POST",
		"/v3/send",
		headers,
		bytes.NewBuffer([]byte(`{
			"personalizations": [{
				"to": [{
					"email": "to@example.com"
				}]
			}], 
			"from": {
				"email": "from@example.com"
			}, 
			"subject": "Subject", 
			"content": [{
				"type": "text/plain", 
				"value": "Content1"
			}, {
				"type": "text/plain", 
				"value": "Content2"
			}]
		}`)),
		http.StatusBadRequest,
		`{"errors":[{"message":"If present, text/plain and text/html may only be provided once.","field":"content","help":null}]}`,
	)

	// NG (multiple content text/html, text/html)
	requestTest(
		t,
		"POST",
		"/v3/send",
		headers,
		bytes.NewBuffer([]byte(`{
			"personalizations": [{
				"to": [{
					"email": "to@example.com"
				}]
			}], 
			"from": {
				"email": "from@example.com"
			}, 
			"subject": "Subject", 
			"content": [{
				"type": "text/html", 
				"value": "<h1>Content1</h1>"
			}, {
				"type": "text/html", 
				"value": "<h1>Content2</h1>"
			}]
		}`)),
		http.StatusBadRequest,
		`{"errors":[{"message":"If present, text/plain and text/html may only be provided once.","field":"content","help":null}]}`,
	)

	// OK (multiple content text/plain, text/html)
	requestTest(
		t,
		"POST",
		"/v3/send",
		headers,
		bytes.NewBuffer([]byte(`{
			"personalizations": [{
				"to": [{
					"email": "to@example.com"
				}]
			}], 
			"from": {
				"email": "from@example.com"
			}, 
			"subject": "Subject", 
			"content": [{
				"type": "text/plain", 
				"value": "Content1"
			}, {
				"type": "text/html", 
				"value": "<h1>Content2</h1>"
			}]
		}`)),
		http.StatusAccepted,
		``,
	)

	// OK (attachements)
	requestTest(
		t,
		"POST",
		"/v3/send",
		headers,
		bytes.NewBuffer([]byte(`{
			"personalizations": [{
				"to": [{
					"email": "to@example.com"
				}]
			}], 
			"from": {
				"email": "from@example.com"
			}, 
			"subject": "Subject", 
			"content": [{
				"type": "text/plain", 
				"value": "Content"
			}],
			"attachments": [{
				"content": "dGVzdA==", 
				"type": "text/plain", 
				"filename": "attachment.txt"
			}]
		}`)),
		http.StatusAccepted,
		``,
	)

	// OK (multiple attachements)
	requestTest(
		t,
		"POST",
		"/v3/send",
		headers,
		bytes.NewBuffer([]byte(`{
			"personalizations": [{
				"to": [{
					"email": "to@example.com"
				}]
			}], 
			"from": {
				"email": "from@example.com"
			}, 
			"subject": "Subject", 
			"content": [{
				"type": "text/plain", 
				"value": "Content"
			}],
			"attachments": [{
				"content": "dGVzdA==", 
				"type": "text/plain", 
				"filename": "attachment.txt"
			}]
		}`)),
		http.StatusAccepted,
		``,
	)

	// NG (attachements content is not BASE64)
	requestTest(
		t,
		"POST",
		"/v3/send",
		headers,
		bytes.NewBuffer([]byte(`{
			"personalizations": [{
				"to": [{
					"email": "to@example.com"
				}]
			}], 
			"from": {
				"email": "from@example.com"
			}, 
			"subject": "Subject", 
			"content": [{
				"type": "text/plain", 
				"value": "Content"
			}],
			"attachments": [{
				"content": "dGVzdA==", 
				"type": "text/plain", 
				"filename": "attachment1.txt"
			}, {
				"content": "dGVzdA==", 
				"type": "text/plain", 
				"filename": "attachment2.txt"
			}]
		}`)),
		http.StatusAccepted,
		``,
	)

	// NG (attachements content is not BASE64 in multiple)
	requestTest(
		t,
		"POST",
		"/v3/send",
		headers,
		bytes.NewBuffer([]byte(`{
			"personalizations": [{
				"to": [{
					"email": "to@example.com"
				}]
			}], 
			"from": {
				"email": "from@example.com"
			}, 
			"subject": "Subject", 
			"content": [{
				"type": "text/plain", 
				"value": "Content"
			}],
			"attachments": [{
				"content": "dGVzdA==", 
				"type": "text/plain", 
				"filename": "attachment1.txt"
			}, {
				"content": "NOT_BASE64", 
				"type": "text/plain", 
				"filename": "attachment2.txt"
			}]
		}`)),
		http.StatusBadRequest,
		`{"errors":[{"message":"The attachment content must be base64 encoded.","field":"attachments.1.content","help":"http://sendgrid.com/docs/API_Reference/Web_API_v3/Mail/errors.html#message.attachments.content"}]}`,
	)
}
