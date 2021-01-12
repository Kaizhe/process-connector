FROM golang:1.15-alpine3.12 AS builder
RUN apk add --no-cache build-base git
WORKDIR $GOPATH/src/github.com/kaizhe/process-connector
COPY . $GOPATH/src/github.com/kaizhe/process-connector
RUN env GO111MODULE=on GOOS=linux GOARCH=amd64 go build -a -o process-connector

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=builder /go/src/github.com/kaizheh/process-connector/process-connector /process-connector

ENTRYPOINT ["/process-connector"]
