FROM golang:latest

RUN adduser --disabled-password --gecos '' api
USER api

WORKDIR /go/src/src/github.com/alzaburetz/workflow-api/api

COPY . .

RUN export PORT=3000
RUN export DBNAME=heroku_gwrf0w5w

RUN go get ./...

EXPOSE 3000



