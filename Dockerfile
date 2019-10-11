FROM golang:1.11.13-alpine3.10
LABEL maintainer="Siddhartha Basu <siddhartha-basu@northwestern.edu>"
RUN apk add --no-cache git build-base
RUN mkdir -p /static-server
WORKDIR /static-server
COPY go.mod go.sum main.go ./
RUN go mod download
ADD commands commands
ADD logger logger
ADD validate validate
RUN go build -o app

FROM alpine:3.7
RUN apk --no-cache add ca-certificates
COPY --from=0 /static-server/app /usr/local/bin/