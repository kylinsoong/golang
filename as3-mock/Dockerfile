FROM golang:1.21.6-alpine3.18 

LABEL codebase   "https://github.com/kylinsoong/golang/tree/main/as3-mock"
LABEL maintainer "kylinsoong.1214@gmail.com"

WORKDIR /app

ADD as3-mock /app

EXPOSE 443 

CMD ["/app/as3-mock"]
