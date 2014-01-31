#!/bin/bash
#
# launches a mail script & sends the alert to hipchat
#

echo "`date`: $@" >> /var/logs/all.log

HOST="$1"
IP="$2"
TYPE="$3"
NAME="$4"
IDX="$5"
SUBJECT="$6"
BODY="$7"
TIME="$8"
ACTUAL="$9"
OPERATOR="${10}"
THRESHOLD="${11}"

if [ "${NAME}" = "Memory" ] ; then
    exit 0
fi

export SUBJECT="$TYPE ALERT RECEIVED FOR [$NAME:$IDX]: host=$HOST, ip=$IP"
export BODY="$TYPE ALERT RECEIVED FOR [$NAME:$IDX]: host=$HOST, ip=$IP, actual=$ACTUAL, threshold=$THRESHOLD, subject=$SUBJECT"

FROM="chef" ROOM=XYZ /home/lxc/bin/hipchat-notify "${SUBJECT}: ${BODY}"
node /etc/darqios/alerts/nodemailer.js

exit 0
