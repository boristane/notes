FROM golang:1.13.4-alpine3.10 as builder
ADD . /go/src/github.com/boristane/notes
WORKDIR /go/src/github.com/boristane/notes
RUN apk add --no-cache --update git build-base dep python3 jq
RUN pip3 install awscli
RUN dep ensure -v
RUN go build -o bin/notes .
# RUN go test -c -v -timeout 120s
CMD ./bin/notes