FROM golang:alpine3.19 AS builder

LABEL maintainer "kylinsoong.1214@gmail.com"

WORKDIR /app

COPY . .

RUN go mod init github.com/kylinsoong/golang/nim-mock
RUN go mod tidy
RUN go build -o nim-mock -ldflags "-X main.version=1.0.0" cmd/nim-mock/*.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/nim-mock .
COPY --from=builder /app/cert/certificate.crt .
COPY --from=builder /app/cert/private.key .

EXPOSE 443 

CMD ["./nim-mock"]
