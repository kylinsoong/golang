FROM golang:alpine3.19 AS builder

LABEL maintainer "kylinsoong.1214@gmail.com"

WORKDIR /app

COPY . .

RUN go mod init github.com/kylinsoong/golang/as3-mock
RUN go mod tidy
RUN go build -o as3-mock -ldflags "-X main.version=1.0.0" cmd/as3-mock/*.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/as3-mock .
COPY --from=builder /app/cert/certificate.crt .
COPY --from=builder /app/cert/private.key .

EXPOSE 443 

CMD ["./as3-mock"]
