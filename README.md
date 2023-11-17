# SendGrid Mock API

SendGrid Dev is SengGrid mock API for test your sendgrid emails during development.

[SengGrid MailDev](https://hub.docker.com/r/ykanazawa/sendgrid-maildev) is Docker container with SendGrid Mock API + [MailDev](https://maildev.github.io/maildev/).

## Requirements

- Go 1.21+

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
