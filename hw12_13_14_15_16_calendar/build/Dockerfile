# Собираем в гошке
FROM golang:1.23 as build

ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o /opt/calendar/calendar cmd/calendar/*

RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o /opt/calendar/calendar_scheduler cmd/calendar_scheduler/*

RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o /opt/calendar/calendar_sender cmd/calendar_sender/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="student@otus.ru"

COPY --from=build /opt/calendar/calendar /opt/calendar/calendar
COPY --from=build /opt/calendar/calendar_scheduler /opt/calendar/calendar_scheduler
COPY --from=build /opt/calendar/calendar_sender /opt/calendar/calendar_sender

ENV CONFIG_FILE /etc/calendar/config.yaml
COPY ./configs/config.yaml ${CONFIG_FILE}
COPY ./migrations /migrations
WORKDIR /opt/calendar

CMD ["/opt/calendar/calendar", "-config", "/etc/calendar/config.yaml"]
