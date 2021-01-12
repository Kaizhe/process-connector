FROM golang:1.15 AS builder
RUN apt-get update
WORKDIR $GOPATH/src/github.com/kaizhe/process-connector
COPY . $GOPATH/src/github.com/kaizhe/process-connector
RUN echo $GOPATH
RUN env GO111MODULE=on CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -o process-connector

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=builder /go/src/github.com/kaizhe/process-connector/process-connector /process-connector

ENTRYPOINT ["/process-connector"]
