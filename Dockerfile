FROM golang:alpine AS build-env

RUN apk update && apk add curl git
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN mkdir -p /go/src/github.com/ahume/delete-stalled-concourse-workers
WORKDIR /go/src/github.com/ahume/delete-stalled-concourse-workers

COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only

COPY . .
RUN go build -o app
RUN chmod +x app

FROM alpine
COPY --from=build-env /go/src/github.com/ahume/delete-stalled-concourse-workers/app /app

