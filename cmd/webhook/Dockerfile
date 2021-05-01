FROM golang:1.16

WORKDIR /gobuild

ADD . ./

RUN go build -a -o app ./cmd/webhook/

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=0 /gobuild/app .

CMD ["./app"]