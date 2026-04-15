FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/calendar ./cmd/calendar/

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/calendar /app/calendar
COPY configs/config.docker.yaml /app/config.yaml
COPY migrations /app/migrations

EXPOSE 8080 50051

ENTRYPOINT ["/app/calendar"]
CMD ["-config", "/app/config.yaml"]
