# base go image
FROM golang:1.18-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o brokerAPP ./cmd/api

RUN chmod +x /app/brokerAPP


# build the tiny broker image
FROM alpine:lastest

RUN mkdir /app

COPY --from=builder /app/brokerAPP /app

CMD [ "/app/brokerApp"]