FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/scheduler ./cmd/calendar_scheduler/

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/scheduler /app/scheduler
COPY configs/config.docker.yaml /app/config.yaml

ENTRYPOINT ["/app/scheduler"]
CMD ["-config", "/app/config.yaml"]
