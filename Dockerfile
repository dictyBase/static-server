FROM golang:1.13.15-buster AS builder
LABEL maintainer="Siddhartha Basu <siddhartha-basu@northwestern.edu>"
ENV GOPROXY https://proxy.golang.org
ENV CGO_ENABLED=0
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

FROM gcr.io/distroless/static
COPY --from=builder /static-server/app /usr/local/bin/
