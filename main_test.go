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

	// NG (Not POST)
	apitest.New().
		Handler(route.Init()).
		Get("/v3/mail/send").
		Expect(t).
		Body(`{"errors":[{"message":"POST method allowed only","field":null,"help":null}]}`).
		Status(http.StatusMethodNotAllowed).
		End()

	// NG (Missing Authorization)
	apitest.New().
		Handler(route.Init()).
		Post("/v3/mail/send").
		Expect(t).
		Body(`{"errors":[{"message":"The provided authorization grant is invalid, expired, or revoked","field":null,"help":null}]}`).
		Status(http.StatusUnsupportedMediaType).
		End()

	// NG (Missing Content-Type)
	apitest.New().
		Handler(route.Init()).
		Post("/v3/mail/send").
		Headers(map[string]string{"Authorization": "Bearer " + os.Getenv("SENDGRID_DEV_API_KEY")}).
		Expect(t).
		Body(`{"errors":[{"message":"Content-Type should be application/json","field":null,"help":null}]}`).
		Status(http.StatusUnsupportedMediaType).
		End()

	// NG (Content-Type is not application/json)
	apitest.New().
		Handler(route.Init()).
		Post("/v3/mail/send").
		Headers(map[string]string{"Authorization": "Bearer " + os.Getenv("SENDGRID_DEV_API_KEY")}).
		Headers(map[string]string{"Content-Type": "text/plain"}).
		Expect(t).
		Body(`{"errors":[{"message":"Content-Type should be application/json","field":null,"help":null}]}`).
		Status(http.StatusUnsupportedMediaType).
		End()

	// NG (Missing PostData)
	apitest.New().
		Handler(route.Init()).
		Post("/v3/mail/send").
		Headers(map[string]string{"Authorization": "Bearer " + os.Getenv("SENDGRID_DEV_API_KEY")}).
		JSON(``).
		Expect(t).
		Body(`{"errors":[{"message":"Bad Request","field":null,"help":null}]}`).
		Status(http.StatusBadRequest).
		End()

	// NG (Missing personalizations)
	apitest.New().
		Handler(route.Init()).
		Post("/v3/mail/send").
		Headers(map[string]string{"Authorization": "Bearer " + os.Getenv("SENDGRID_DEV_API_KEY")}).
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
		Headers(map[string]string{"Authorization": "Bearer " + os.Getenv("SENDGRID_DEV_API_KEY")}).
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
		Headers(map[string]string{"Authorization": "Bearer " + os.Getenv("SENDGRID_DEV_API_KEY")}).
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

	// NG (Missing content)
	apitest.New().
		Handler(route.Init()).
		Post("/v3/mail/send").
		Headers(map[string]string{"Authorization": "Bearer " + os.Getenv("SENDGRID_DEV_API_KEY")}).
		JSON(`{
			"personalizations": [{
				"to": [{
					"email": "to@example.com"
				}]
			}], 
			"from": {
				"email": "from@example.com"
			}, 
			"subject": "Subject"
		}`).
		Expect(t).
		Body(`{"errors":[{"message":"Unless a valid template_id is provided, the content parameter is required. There must be at least one defined content block. We typically suggest both text/plain and text/html blocks are included, but only one block is required.","field":"content","help":"http://sendgrid.com/docs/API_Reference/Web_API_v3/Mail/errors.html#message.content"}]}`).
		Status(http.StatusBadRequest).
		End()

	// OK
	apitest.New().
		Handler(route.Init()).
		Post("/v3/mail/send").
		Headers(map[string]string{"Authorization": "Bearer " + os.Getenv("SENDGRID_DEV_API_KEY")}).
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

	// OK (multiple to, cc, bcc and reply-to with name)
	apitest.New().
		Handler(route.Init()).
		Post("/v3/mail/send").
		Headers(map[string]string{"Authorization": "Bearer " + os.Getenv("SENDGRID_DEV_API_KEY")}).
		JSON(`{
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
		}`).
		Expect(t).
		Body(``).
		Status(http.StatusAccepted).
		End()

	// OK (multiple personalizations)
	apitest.New().
		Handler(route.Init()).
		Post("/v3/mail/send").
		Headers(map[string]string{"Authorization": "Bearer " + os.Getenv("SENDGRID_DEV_API_KEY")}).
		JSON(`{
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
		}`).
		Expect(t).
		Body(``).
		Status(http.StatusAccepted).
		End()

	// OK (text/html)
	apitest.New().
		Handler(route.Init()).
		Post("/v3/mail/send").
		Headers(map[string]string{"Authorization": "Bearer " + os.Getenv("SENDGRID_DEV_API_KEY")}).
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
				"type": "text/html", 
				"value": "<h1>Content</h1>"
			}]
		}`).
		Expect(t).
		Body(``).
		Status(http.StatusAccepted).
		End()

	// OK (multiple content text/plain, text/html)
	apitest.New().
		Handler(route.Init()).
		Post("/v3/mail/send").
		Headers(map[string]string{"Authorization": "Bearer " + os.Getenv("SENDGRID_DEV_API_KEY")}).
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
				"value": "Content1"
			}, {
				"type": "text/html", 
				"value": "<h1>Content2</h1>"
			}]
		}`).
		Expect(t).
		Body(``).
		Status(http.StatusAccepted).
		End()

	// OK (attachements)
	apitest.New().
		Handler(route.Init()).
		Post("/v3/mail/send").
		Headers(map[string]string{"Authorization": "Bearer " + os.Getenv("SENDGRID_DEV_API_KEY")}).
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
			}],
			"attachments": [{
				"content": "dGVzdA==", 
				"type": "text/plain", 
				"filename": "attachment.txt"
			}]
		}`).
		Expect(t).
		Body(``).
		Status(http.StatusAccepted).
		End()

	// OK (multiple attachements)
	apitest.New().
		Handler(route.Init()).
		Post("/v3/mail/send").
		Headers(map[string]string{"Authorization": "Bearer " + os.Getenv("SENDGRID_DEV_API_KEY")}).
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
		}`).
		Expect(t).
		Body(``).
		Status(http.StatusAccepted).
		End()

	// NG (attachements content is not BASE64)
	apitest.New().
		Handler(route.Init()).
		Post("/v3/mail/send").
		Headers(map[string]string{"Authorization": "Bearer " + os.Getenv("SENDGRID_DEV_API_KEY")}).
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
			}],
			"attachments": [{
				"content": "NOT BASE64", 
				"type": "text/plain", 
				"filename": "attachment.txt"
			}]
		}`).
		Expect(t).
		Body(`{"errors":[{"message":"The attachment content must be base64 encoded.","field":"attachments.0.content","help":"http://sendgrid.com/docs/API_Reference/Web_API_v3/Mail/errors.html#message.attachments.content"}]}`).
		Status(http.StatusBadRequest).
		End()

	// NG (attachements content is not BASE64 in multiple)
	apitest.New().
		Handler(route.Init()).
		Post("/v3/mail/send").
		Headers(map[string]string{"Authorization": "Bearer " + os.Getenv("SENDGRID_DEV_API_KEY")}).
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
		}`).
		Expect(t).
		Body(`{"errors":[{"message":"The attachment content must be base64 encoded.","field":"attachments.1.content","help":"http://sendgrid.com/docs/API_Reference/Web_API_v3/Mail/errors.html#message.attachments.content"}]}`).
		Status(http.StatusBadRequest).
		End()

	// OK (with SMTP Auth)
	os.Setenv("SENDGRID_DEV_SMTP_USERNAME", "username@example.com")
	apitest.New().
		Handler(route.Init()).
		Post("/v3/mail/send").
		Headers(map[string]string{"Authorization": "Bearer " + os.Getenv("SENDGRID_DEV_API_KEY")}).
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
