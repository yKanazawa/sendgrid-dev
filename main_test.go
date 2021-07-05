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

	// NG (Missing PostData)
	apitest.New().
		Handler(route.Init()).
		Post("/v3/mail/send").
		Expect(t).
		Body(`{"errors":[{"message":"Bad Request","field":null,"help":null}]}`).
		Status(http.StatusBadRequest).
		End()

}
