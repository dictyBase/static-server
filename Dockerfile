FROM golang:1.17.2-alpine3.14 as builder
LABEL maintainer="Siddhartha Basu <siddhartha-basu@northwestern.edu>"
ENV GOPROXY https://proxy.golang.org
ENV GO111MODULE=auto \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
RUN apk add --no-cache git build-base upx binutils
RUN mkdir -p /static-server
WORKDIR /static-server
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY main.go ./
ADD commands commands
ADD logger logger
ADD handlers handlers
RUN go build -o /bin/app \
    -a \
    -ldflags "-s -w -extldflags '-static'" \
    -installsuffix cgo \
    -tags netgo \
    -o /bin/app 
RUN strip /bin/app \
    && upx -q -9 /bin/app

FROM alpine:3.14
RUN apk --no-cache add ca-certificates
COPY --from=builder /bin/app /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/app"]
