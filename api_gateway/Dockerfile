FROM golang:1.17-alpine
# install redis
RUN apk add --no-cache redis

RUN mkdir /api
COPY . /api
WORKDIR /api

# Set environment variables
# ENV ENVIRONMENT=develop
# ENV POSTGRES_USER=smn
# ENV POSTGRES_PASSWORD=123
# ENV POSTGRES_HOST=database
# ENV POSTGRES_PORT=5434
# ENV POSTGRES_DATABASE=smn
# ENV LOG_LEVEL=debug
# ENV HTTP_PORT=:5050
# ENV USER_SERVICE_HOST=user_service
# ENV USER_SERVICE_PORT=8000
# ENV POST_SERVICE_HOST=post_service
# ENV POST_SERVICE_PORT=8010
# ENV COMMENT_SERVICE_HOST=comment_service
# ENV COMMENT_SERVICE_PORT=8020
# ENV REDIS_HOST=redisdb
# ENV REDIS_PORT=6380
# ENV SIGNING_KEY=123
# ENV CTX_TIMEOUT=7
# ENV CSV_FILE_PATH=./config/casbin.csv
# ENV AUTH_PATH=./config/auth.conf

RUN go build -o main cmd/main.go
CMD ./main
EXPOSE 5050
