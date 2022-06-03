FROM golang:alpine3.16

RUN mkdir /go/crypto_checker
WORKDIR /go/crypto_checker

RUN cd /go/crypto_checker
COPY ./ /go/crypto_checker
RUN go build -o app

EXPOSE 8080

ENTRYPOINT ["/go/crypto_checker/app"]