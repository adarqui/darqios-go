#!/bin/bash

echo "`date`: Restarting aspera" >> /tmp/mitigate.log

service asperalee restart
service asperanoded restart
/usr/sbin/sshd -p 33001

exit 0
