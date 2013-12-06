#!/bin/bash

echo "`date`: Restarting PHP5-FPM" >> /tmp/mitigate.log

if service php5-fpm restart ; then
	echo "Successfully restarted php5-fpm"
	exit 0
fi

echo "Failed to restart php5-fpm"

exit -1
