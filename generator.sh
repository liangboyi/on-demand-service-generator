#!/bin/bash -e
SERVICE_NAME=`echo $1 | tr A-Z a-z`

if [ ! $# -eq 1 ] ; then
	echo "you need provide a name to your service adapter. like:"
	echo "./generator.sh redis"
	exit 1
elif [ "$SERVICE_NAME" == "template" ]; then
	echo "key word [template] can not be used, please provide another one ."
	exit 1
fi

SERVICE_NAME=`echo $1 | tr A-Z a-z`
rm -rf $SERVICE_NAME
cp -r template $SERVICE_NAME
cd $SERVICE_NAME

#update file self-service-adapter to $SERVICE_NAME-service-adapter
grep -l -r self- |xargs -l sed -i 's/self-/'${SERVICE_NAME}'-/g'

#update folder
find . -type d -name 'self-*' -printf %p\\n | sort -r |sed 's,^\(.*\)self-\([^/]*\)$,mv & \1'$SERVICE_NAME'-\2,' | /bin/sh
