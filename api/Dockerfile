FROM golang:latest

RUN adduser --disabled-password --gecos '' api
USER api

WORKDIR /go/src/api

COPY . .

# RUN go get gopkg.in/mgo.v2
# RUN go get -u github.com/gorilla/mux
# RUN go get -u github.com/gorilla/handlers

# RUN go get github.com/gomodule/redigo/redis
# RUN go get github.com/satori/go.uuid
# RUN go get golang.org/x/crypto/bcrypt
RUN go get ./...

EXPOSE 3000



