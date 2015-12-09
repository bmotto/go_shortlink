FROM golang:1.3-onbuild
ADD . /code
WORKDIR /code
RUN go get github.com/go-redis/redis
CMD ./redis-3.0.5/src/redis-cli ./go_shortlink
