# FROM golang:latest

# RUN adduser --disabled-password --gecos '' api
# USER api

# WORKDIR /go/src/api

# COPY . .

# RUN go get gopkg.in/mgo.v2
# RUN go get -u github.com/gorilla/mux
# RUN go get -u github.com/gorilla/handlers
# RUN go get github.com/pilu/fresh
# RUN go get github.com/gomodule/redigo/redis
# RUN go get github.com/satori/go.uuid
# RUN go get golang.org/x/crypto/bcrypt
# RUN go get ./...

# EXPOSE 3000
# CMD [ "fresh" ]

FROM heroku/heroku:18-build as build

COPY . /api
WORKDIR /api

# Setup buildpack
RUN mkdir -p /tmp/buildpack/heroku/go /tmp/build_cache /tmp/env
RUN curl https://codon-buildpacks.s3.amazonaws.com/buildpacks/heroku/go.tgz | tar xz -C /tmp/buildpack/heroku/go

#Execute Buildpack
RUN STACK=heroku-18 /tmp/buildpack/heroku/go/bin/compile /api /tmp/build_cache /tmp/env

# Prepare final, minimal image
FROM heroku/heroku:18

COPY --from=build /api /api
ENV HOME /api
WORKDIR /api
RUN useradd -m heroku
USER heroku
CMD /api/bin/api