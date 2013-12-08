#!/bin/bash

echo "`date`: Restarting haproxy" >> /tmp/mitigate.log

if haproxy -D -f /etc/haproxy/haproxy.cfg -p /var/run/haproxy.pid -sf $(cat /var/run/haproxy.pid) ; then
	echo "Successfully restarted haproxy"
	exit 0
fi

echo "Failed to restart haproxy"

exit -1
