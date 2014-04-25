#!/bin/bash

HOST=$1
shift
REST=$@

./find_stuff.sh $HOST loadavg.last1min $@
