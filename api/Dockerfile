FROM golang:latest

RUN adduser --disabled-password --gecos '' api
USER api

WORKDIR /go/src/src/github.com/alzaburetz/workflow-api/api

COPY . .

RUN go get ./...

EXPOSE 3000



