# syntax=docker/dockerfile:1

FROM golang:1.18.1-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY *.json ./
RUN go mod download

COPY *.go ./

RUN go build -o /leader-election

EXPOSE 8080

CMD [ "/leader-election" ]