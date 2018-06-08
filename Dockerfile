FROM golang:1.10.2-alpine3.7
LABEL maintainer="Siddhartha Basu <siddhartha-basu@northwestern.edu>"
RUN apk add --no-cache git build-base \
    && go get github.com/golang/dep/cmd/dep
RUN mkdir -p /go/src/github.com/static-server
WORKDIR /go/src/github.com/static-server
COPY Gopkg.* main.go ./
ADD commands commands
ADD logger logger
ADD validate validate
RUN dep ensure \
    && go build -o app

FROM alpine:3.7
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/src/github.com/static-server/app /usr/local/bin/
