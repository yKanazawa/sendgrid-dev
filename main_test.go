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
