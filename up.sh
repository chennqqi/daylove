#!/bin/bash
set -ex
DNSSERVERS=" --dns=208.67.222.222 --dns=208.67.220.220 --dns=8.8.8.8 --dns=8.8.4.4 "
if [ ! -f config.toml ]; then
    cp config.toml.dist config.toml
fi
if [ $(docker network ls | grep daylove-network | wc -l ) -eq 0 ]; then
    docker network create -d bridge daylove-network
fi
docker run ${DNSSERVERS} --rm --name go-build -v $HOME/go:/go -v $(pwd):/www golang sh -c "cd /www ;go get -v ; go build -o /www/daylove "
if [ $(docker ps -a | grep gs_db | wc -l) -le 0 ]; then
    docker run ${DNSSERVERS} --restart=always --net=daylove-network -d --name gs_db  netroby/docker-mysql
    while true; do
        if [ $(docker logs gs_db 2>&1 | grep "ready for connections" | wc -l)  -ge 2 ]; then
            break;
        else
            echo "not ready, waiting"
            sleep 2
        fi
    done
    docker cp sql/bak.sql gs_db:/root/
    docker exec gs_db sh -c "mysql < /root/bak.sql"
fi
if [ $(docker ps -a | grep daylove | wc -l) -ge 1 ]; then
    docker rm -vf daylove
fi
docker run ${DNSSERVERS} --restart=always --net=daylove-network -d -p 127.0.0.1:8702:8080  -v $(pwd):/www --name daylove debian  sh -c "cd /www && /www/daylove"

