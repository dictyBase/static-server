FROM golang:1.13.15-alpine3.12
LABEL maintainer="Siddhartha Basu <siddhartha-basu@northwestern.edu>"
RUN apk add --no-cache git build-base
RUN mkdir -p /static-server
WORKDIR /static-server
COPY go.mod ./
COPY go.sum ./
COPY main.go ./
RUN go mod download
ADD commands commands
ADD logger logger
ADD validate validate
RUN go build -o app

FROM alpine:3.12
RUN apk --no-cache add ca-certificates
COPY --from=0 /static-server/app /usr/local/bin/