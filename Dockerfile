FROM golang:1.3-onbuild
ADD ./release /docker
RUN go get github.com/go-redis/redis
CMD ./release/go_shortlink ./release/redis-3.0.5/src/redis-cli
