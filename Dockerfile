FROM golang:1.12 as builder

LABEL maintainer="Andrew Davidson <andrew@amdavidson.com>"

WORKDIR /go/src/github.com/amdavidson/shorten-go

COPY . .

RUN go get -d -v ./...

RUN go build -o shorten-go .

CMD ["./shorten-go"]
