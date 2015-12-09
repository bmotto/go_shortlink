# install go_shortlink
go install github.com/bmotto/go_shortlink

# copy required files in release folder
cd $GOPATH/src/github.com/bmotto/go_shortlink
mkdir -f ./release
cp config.yml ./release/config.yml
cp $GOPATH/bin/go_shortlink ./release/go_shortlink

# download and compile redis
wget http://download.redis.io/releases/redis-3.0.5.tar.gz /tmp/redis-3.0.5.tar.gzredis-3.0.5.tar.gz
tar -xvzf /tmp/redis-3.0.5.tar.gz /opt/redis
cd /opt/redis-3.0.5
make
# copy required files in release folder
cd $GOPATH/src/github.com/bmotto/go_shortlink
cp -r /opt/redis-3.0.5 ./release/redis-3.0.5

#docker
docker build -t docker_shortlink .
