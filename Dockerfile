FROM --platform=linux/arm64 golang:1.16-alpine AS builder

WORKDIR /app
COPY . ./
RUN go mod download
RUN GOOS=linux GOARCH=arm GOARM=5 CGO_ENABLED=0 go build -o ./app

## Deploy
FROM --platform=linux/arm64 alpine:3.16
COPY --from=builder /app/app /app
COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY config.toml ./config.toml
ENTRYPOINT [ "/app", "binance" ]