# Base go image

FROM golang:1.22.2-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o loggerApp ./cmd/api

RUN chmod +x /app/loggerApp


#Build the tiny docker image

FROM alpine:latest

RUN mkdir /app

COPY --from=builder  /app/loggerApp /app

CMD [ "/app/loggerApp" ]