#!/bin/sh

MACHINE=e1z1r${1}p${2}
PORT=${3:-4242}

while true
do
	screencapture -x /tmp/screen.png
	echo Sending to $MACHINE...
	nc $MACHINE $PORT < /tmp/screen.png
	echo Sleeping...
	sleep 2
done
