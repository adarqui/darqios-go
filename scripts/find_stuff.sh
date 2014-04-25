#!/bin/bash

# ./find_stuff.sh <host> loadavg.last1min 0 40 2 "2014-04-10" "2014-04-20"

usage() {
    echo "./find_stuff <HOST> <STATE_OBJ> <VAL_START> <VAL_END> <LIMIT> [DATE_START] [DATE_END]"
    exit 1
}

# "Thu Apr 24 12:33:34 EDT 2011"

mul_date() {
    echo $(( $1 * 1000 ))
}

fix_date() {
    echo "$1T00:00:00.000Z"
}

grab_date_now() {
#    mul_date `date +%s`
    fix_date "`date +%Y-%m-%d`"
}

grab_date() {
#    mul_date `date -d "$1" +%s`
    fix_date "$1"
}

if [ $# -lt 5 ] ; then
    usage
fi

HOST=$1
STATE_OBJ=$2
VAL_START=$3
VAL_END=$4
LIMIT=$5

# 2010-04-30T00:00:00.000Z
DATE_START=`grab_date "1975-01-01"`
DATE_END=`grab_date_now`

if [ $# -eq 7 ] ; then
    DATE_START=`grab_date "$6"`
    DATE_END=`grab_date "$7"`
fi

CMD="use darqios\ndb.state.find({'host': '$HOST', '$STATE_OBJ': {\$gte:$VAL_START, \$lte:$VAL_END}, 'ts': {\$gte:new ISODate('$DATE_START'), \$lte: new ISODate('$DATE_END')}}).limit($LIMIT).pretty()\n"

echo $CMD

printf "$CMD" | mongo
