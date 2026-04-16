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

ENV BIN_FILE_CALENDAR /opt/calendar/calendar
ENV BIN_FILE_SCHEDULER /opt/calendar/calendar_scheduler
ENV BIN_FILE_SENDER /opt/calendar/calendar_sender

COPY --from=build /opt/calendar/calendar ${BIN_FILE_CALENDAR}
COPY --from=build /opt/calendar/calendar_scheduler ${BIN_FILE_SCHEDULER}
COPY --from=build /opt/calendar/calendar_sender ${BIN_FILE_SENDER}

ENV CONFIG_FILE /etc/calendar/config.yaml
COPY ./configs/config.yaml ${CONFIG_FILE}
COPY ./migrations /migrations
WORKDIR /opt/calendar

CMD ["${BIN_FILE_CALENDAR}", "-config", "/etc/calendar/config.yaml"]