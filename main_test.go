package main

import (
	"net/http"
	"os"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/yKanazawa/sendgrid-dev/route"
)

func TestSend(t *testing.T) {
	os.Setenv("SENDGRID_DEV_TEST", "1")
	os.Setenv("SENDGRID_DEV_APIKEY", "SG.xxxxx")
	// NG (Not POST)
	apitest.New().
		Handler(route.Init()).
		Get("/v3/mail/send").
		Expect(t).
		Body(`{"errors":[{"message":"POST method allowed only","field":null,"help":null}]}`).
		Status(http.StatusMethodNotAllowed).
		End()

	// NG (Missing Content-Type)
	apitest.New().
		Handler(route.Init()).
		Post("/v3/mail/send").
		Expect(t).
		Body(`{"errors":[{"message":"Content-Type should be application/json","field":null,"help":null}]}`).
		Status(http.StatusUnsupportedMediaType).
		End()

	// NG (Content-Type is not application/json)
	apitest.New().
		Handler(route.Init()).
		Post("/v3/mail/send").
		Headers(map[string]string{"Content-Type": "text/plain"}).
		Expect(t).
		Body(`{"errors":[{"message":"Content-Type should be application/json","field":null,"help":null}]}`).
		Status(http.StatusUnsupportedMediaType).
		End()

	// NG (Missing PostData)
	apitest.New().
		Handler(route.Init()).
		Post("/v3/mail/send").
		JSON(``).
		Expect(t).
		Body(`{"errors":[{"message":"Bad Request","field":null,"help":null}]}`).
		Status(http.StatusBadRequest).
		End()

	// NG (Missing personalizations)
	apitest.New().
		Handler(route.Init()).
		Post("/v3/mail/send").
		JSON(`{
			"from": {
				"email": "from@example.com"
			},
			"subject": "Subject", 
			"content": [{
				"type": "text/plain", 
				"value": "Content"
			}]
		}`).
		Expect(t).
		Body(`{"errors":[{"message":"The personalizations field is required and must have at least one personalization.","field":"personalizations","help":"http://sendgrid.com/docs/API_Reference/Web_API_v3/Mail/errors.html#-Personalizations-Errors"}]}`).
		Status(http.StatusBadRequest).
		End()

	// NG (Missing from.Email)
	apitest.New().
		Handler(route.Init()).
		Post("/v3/mail/send").
		JSON(`{
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
		}`).
		Expect(t).
		Body(`{"errors":[{"message":"The from object must be provided for every email send. It is an object that requires the email parameter, but may also contain a name parameter.  e.g. {\"email\" : \"example@example.com\"}  or {\"email\" : \"example@example.com\", \"name\" : \"Example Recipient\"}.","field":"from.email","help":"http://sendgrid.com/docs/API_Reference/Web_API_v3/Mail/errors.html#message.from"}]}`).
		Status(http.StatusBadRequest).
		End()

	// NG (Missing subject)
	apitest.New().
		Handler(route.Init()).
		Post("/v3/mail/send").
		JSON(`{
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
		}`).
		Expect(t).
		Body(`{"errors":[{"message":"The subject is required. You can get around this requirement if you use a template with a subject defined or if every personalization has a subject defined.","field":"subject","help":"http://sendgrid.com/docs/API_Reference/Web_API_v3/Mail/errors.html#message.subject"}]}`).
		Status(http.StatusBadRequest).
		End()

	// OK
	apitest.New().
		Handler(route.Init()).
		Post("/v3/mail/send").
		JSON(`{
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
		}`).
		Expect(t).
		Body(``).
		Status(http.StatusAccepted).
		End()
}
