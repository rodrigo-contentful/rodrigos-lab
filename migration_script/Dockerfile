#Get base image
FROM ubuntu:16.04

# copy files and folder
COPY run.sh /

#If we need extra packages
RUN apt-get update
RUN apt-get -y install curl gnupg
RUN curl -sL https://deb.nodesource.com/setup_12.x  | bash -
RUN apt-get -y install nodejs
RUN npm install
RUN npm install -g contentful-cli

RUN chmod 777 run.sh
CMD ./run.sh $CMA_TOKEN $SPACEA $SPACEB 

## .run_docker-sh CMA_TOKEN SPACE_ORIGIN SPACE_DESTINATION