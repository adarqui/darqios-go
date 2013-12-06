#!/bin/bash

echo "$@" >> /tmp/mitigate.log

if redis-server /etc/redis/redis.conf ; then
	exit 0
fi

exit -1
