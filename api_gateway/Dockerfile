FROM golang:1.17-alpine

# install redis
RUN apk add --no-cache redis

RUN mkdir /api
COPY . /api
WORKDIR /api

RUN go build -o main cmd/main.go
CMD ./main
EXPOSE 5050
