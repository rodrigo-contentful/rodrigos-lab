#!/bin/bash

CMA_TOKEN=$1
SPACEA=$2
SPACEB=$3

if [ -z $CMA_TOKEN ]
    then
        echo "I need a CMA token"
        exit 0
fi

if [ -z $SPACEA ]
    then
        echo "Origin space needed"
        exit 0
fi

if [ -z $SPACEB ]
    then
        echo "Destination space needed"
        exit 0
fi

#if [ -z $(docker images sometag) ]
 #   then
    echo "crating image"
    docker build --tag migration-docker .
#fi

docker run --rm -e CMA_TOKEN="$CMA_TOKEN" -e SPACEA="$SPACEA" -e SPACEB="$SPACEB" migration-docker
docker rmi migration-docker
