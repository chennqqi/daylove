#!/bin/sh
if [ ! -z $1 ]; then
	docker logs -f --tail 200 gosense
else
	docker logs --tail 200 gosense
fi