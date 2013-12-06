#!/bin/bash

echo "`date`: Restarting redis" >> /tmp/mitigate.log

if redis-server /etc/redis/redis.conf ; then
	echo "Successfully restarted redis"
	exit 0
fi

echo "Failed to restart redis"

exit -1
