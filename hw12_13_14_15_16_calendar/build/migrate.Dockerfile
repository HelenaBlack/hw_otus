FROM golang:1.23-alpine as builder

RUN go install github.com/pressly/goose/v3/cmd/goose@v3.24.1

FROM alpine:latest

COPY --from=builder /go/bin/goose /usr/local/bin/goose

WORKDIR /app

ENTRYPOINT ["goose"]
