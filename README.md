# SendGrid Mock API

SendGrid Dev is SengGrid mock API for test your sendgrid emails during development.

[SengGrid MailDev](https://hub.docker.com/r/ykanazawa/sendgrid-maildev) is Docker container with SendGrid Mock API + [MailDev](https://maildev.github.io/maildev/).

## Docker (GitHub Container Registry)

Prebuilt multi-arch images (`linux/amd64`, `linux/arm64`) are published from
`main` on every push. Point the SendGrid SDK at this container's
`/v3/mail/send` endpoint; it forwards to whatever SMTP host you configure
(MailDev, Mailpit, MailTrap, real SMTP, …).

```bash
docker pull ghcr.io/paxa/sendgrid-dev:latest

# Example: forward to a local Mailpit on the same docker network.
docker run --rm -p 3030:3030 \
  -e SENDGRID_DEV_API_KEY=SG.local-dev-key \
  -e SENDGRID_DEV_SMTP_SERVER=mailpit:1025 \
  ghcr.io/paxa/sendgrid-dev:latest
```

Available tags: `latest`, `main`, `sha-<short>`, and `v*` for release tags.

Build locally:

```bash
docker build -t sendgrid-dev:local .
```

## Requirements

- Go 1.26+

## Debug

### Sample with MailDev (Can work by default)

Run maildev
```
docker pull maildev/maildev
docker run -p 1080:1080 -p 1025:1025 maildev/maildev
```

Run SendGrid Mock API
```
go run main.go
```

Send mail by curl
```
curl --request POST \
  --url http://localhost:3030/v3/mail/send \
  --header 'Authorization: Bearer SG.xxxxx' \
  --header 'Content-Type: application/json' \
  --data '{"personalizations": [{ 
    "to": [{"email": "to@example.com"}]}], 
    "from": {"email": "from@example.com"}, 
    "subject": "Test Subject", 
    "content": [{"type": "text/plain", "value": "Test Content"}] 
  }'
```

Check with maildev

http://localhost:1080/

### Sample with MailTrap (with SMTP Auth)

Run SendGrid Mock API
```
export SENDGRID_DEV_API_SERVER=:3030
export SENDGRID_DEV_API_KEY=SG.xxxxx
export SENDGRID_DEV_SMTP_SERVER=smtp.mailtrap.io:25
export SENDGRID_DEV_SMTP_USERNAME=mailtrap_username
export SENDGRID_DEV_SMTP_PASSWORD=mailtrap_password
go run main.go
```

Send mail by curl
```
curl --request POST \
  --url http://localhost:3030/v3/mail/send \
  --header 'Authorization: Bearer SG.xxxxx' \
  --header 'Content-Type: application/json' \
  --data '{"personalizations": [{ 
    "to": [{"email": "to@example.com"}]}], 
    "from": {"email": "from@example.com"}, 
    "subject": "Test Subject", 
    "content": [{"type": "text/plain", "value": "Test Content"}] 
  }'
```

Check with mailtrap Inbox

https://mailtrap.io/inboxes

## Test

```
go test
```

## Build

### x86_64

```
env GOOS=linux GOARCH=amd64 go build -o sendgrid-dev_x86_64 main.go
```

### arm64

```
env GOOS=linux GOARCH=arm64 go build -o sendgrid-dev_aarch64 main.go
```
