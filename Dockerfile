FROM golang:1.11.10-alpine3.9
LABEL maintainer="Siddhartha Basu <siddhartha-basu@northwestern.edu>"
RUN apk add --no-cache git build-base
RUN mkdir -p /static-server
WORKDIR /static-server
ADD commands commands
ADD logger logger
ADD validate validate
COPY go.mod go.sum ./
RUN go get ./...
RUN go build -o app

FROM alpine:3.7
RUN apk --no-cache add ca-certificates
COPY --from=0 /static-server/app /usr/local/bin/