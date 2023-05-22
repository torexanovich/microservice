FROM golang:1.17-alpine

# installing make
# RUN apk add --no-cache make git \
#     && go get -u -d github.com/golang-migrate/migrate/cmd/migrate \
#     && cd $GOPATH/src/github.com/golang-migrate/migrate/cmd/migrate \
#     && git checkout $(git describe --tags $(git rev-list --tags --max-count=1)) \
#     && go build -tags 'postgres' -o /usr/local/bin/migrate \
#     && apk del git

# ARG MAKE_TARGET

RUN mkdir /comment
COPY . /comment
WORKDIR /comment
# RUN make $MAKE_TARGET

# Set environment variables
ENV POSTGRES_USER=smn
ENV POSTGRES_PASSWORD=123
ENV POSTGRES_HOST=database
ENV POSTGRES_PORT=5434
ENV POSTGRES_DATABASE=smn
ENV ENVIRONMENT=develop
ENV LOG_LEVEL=debug
ENV USER_SERVICE_HOST=user_service
ENV USER_SERVICE_PORT=8000
ENV POST_SERVICE_HOST=post_service
ENV POST_SERVICE_PORT=8010
ENV COMMENT_SERVICE_HOST=comment_service
ENV COMMENT_SERVICE_PORT=:8020

RUN go build -o main cmd/main.go
CMD ./main
EXPOSE 8020
