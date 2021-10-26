#!/bin/bash

# docker run -d -p 80:80 docker/getting-started
echo "crating image"
docker build --tag ctf-docker .
#fi

docker run --rm -v $PWD/task:/ctfanalyser/task ctf-docker
docker rmi ctf-docker