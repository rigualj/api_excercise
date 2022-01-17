FROM golang:alpine

WORKDIR /go/src/app
COPY *.go . 

RUN go mod init
RUN go get github.com/golang/gddo/httputil/header
RUN go build -o webserver

CMD ./webserver