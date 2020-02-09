FROM golang:latest

RUN adduser --disabled-password --gecos '' api
USER api

WORKDIR /go/src/app

COPY . .

RUN go get github.com/pilu/fresh
RUN go get gopkg.in/mgo.v2
RUN go get -u github.com/gorilla/mux

EXPOSE 3000
CMD [ "fresh" ]
