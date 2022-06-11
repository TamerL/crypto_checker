FROM golang:alpine3.16 as build 

RUN mkdir /go/crypto_checker
WORKDIR /go/crypto_checker

COPY ./main.go ./go.mod ./
RUN go build -o crypto_checker


FROM alpine as runtime

COPY --from=build /go/crypto_checker/crypto_checker /app/crypto_checker

EXPOSE 8080

ENTRYPOINT ["/app/crypto_checker"]