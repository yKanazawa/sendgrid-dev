# SendGrid Mock API

SendGrid Dev is SengGrid mock API for test your sendgrid emails during development.

## Requirements

- Go 1.16+

## Debug

```
export SENDGRID_DEV_API_KEY=SG.xxxxx
export SENDGRID_DEV_SMTP_SERVER=smtp.example.com:25
go run main.go
```

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
