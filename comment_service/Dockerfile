FROM golang:1.17-alpine
RUN mkdir micro
COPY . /micro
WORKDIR /micro
RUN go build -o main cmd/main.go
CMD ./main
EXPOSE 8020