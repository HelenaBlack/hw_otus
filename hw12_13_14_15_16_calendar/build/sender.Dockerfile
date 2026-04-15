FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/sender ./cmd/calendar_sender/

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/sender /app/sender
COPY configs/config.docker.yaml /app/config.yaml

ENTRYPOINT ["/app/sender"]
CMD ["-config", "/app/config.yaml"]
