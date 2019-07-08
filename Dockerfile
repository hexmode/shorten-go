FROM golang:1.12 as builder

LABEL maintainer="Andrew Davidson <andrew@amdavidson.com>"

WORKDIR /go/src/github.com/amdavidson/shorten-go

COPY app/ .

RUN go get -d -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o shorten-go .



FROM scratch

COPY --from=builder /go/src/github.com/amdavidson/shorten-go/shorten-go /bin/shorten-go

ADD templates ./templates

VOLUME /data

ENV dbpath /data/bolt.db

CMD ["/bin/shorten-go"]
