FROM golang:1.17-alpine

RUN mkdir /comment
COPY . /comment
WORKDIR /comment

RUN go build -o main cmd/main.go
CMD ./main
EXPOSE 8020
