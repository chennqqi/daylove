#!/bin/bash
set -ex
if [ ! -f config.toml ]; then
    cp config.toml.dist config.toml
fi
if [ ! -z $1 ]; then
    git pull --rebase
    docker run --rm --name go-build -v $HOME/go:/go -v $(pwd):/www golang sh -c "cd /www ;go get -v ; go build -o /www/gosense "
fi
PIDGS=$(docker exec gosense pidof gosense)
echo "Pid of gosense $PIDGS"
docker exec gosense kill -s HUP $PIDGS
docker ps -a
docker logs -f -t --tail 30  gosense

