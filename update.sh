#!/bin/sh
echo "pulling latest changes from master"
git pull origin master
echo "building new docker image"
docker build -t co.plydge.property .
echo "stopping old container"
docker stop co.plydge.property
echo "removing old container"
docker rm co.plydge.property
echo "starting new container"
docker run --name co.plydge.property --publish 80:80 -d co.plydge.property

