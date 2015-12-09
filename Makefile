

all: build

build:
	# install go_shortlink
	go install github.com/bmotto/go_shortlink

	# copy required files in release folder
	cd $$GOPATH/src/github.com/bmotto/go_shortlink && mkdir -p ./release
	cp config.yml ./release/config.yml
	cp $$GOPATH/bin/go_shortlink ./release/go_shortlink

	# download and compile redis
	wget http://download.redis.io/releases/redis-3.0.5.tar.gz
	tar -xvzf redis-3.0.5.tar.gz -C $$GOPATH/
	cd $$GOPATH/redis-3.0.5 && make
	# copy required files in release folder
	cd $$GOPATH/ && cp -r ./redis-3.0.5 ./src/github.com/bmotto/go_shortlink/release/redis-3.0.5

	#docker
	cd $$GOPATH/src/github.com/bmotto/go_shortlinkcd && sudo docker build -t docker_shortlink .
