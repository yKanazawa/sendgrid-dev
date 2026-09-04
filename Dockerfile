# Multi-stage build for yKanazawa/sendgrid-dev.
# Produces a ~15 MB static image on a scratch base.

FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/sendgrid-dev main.go

FROM scratch
COPY --from=builder /out/sendgrid-dev /sendgrid-dev
# ca-certificates is not needed — the app is a server, not an outbound HTTPS client.
EXPOSE 3030
ENTRYPOINT ["/sendgrid-dev"]
