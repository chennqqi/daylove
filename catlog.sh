#!/bin/sh
if [ ! -z $1 ]; then
	docker logs -f --tail 200 daylove
else
	docker logs --tail 200 daylove
fi