# Lab scripts

Repo that contain experiments scripts with contenful.

## Content type to github - (bash)

dev - on run will download a migration of conten types and push to a github repo. The idea is to have an independent versioning of contenful.

[content_type_into_github](content_type_into_github)

## Merge spaces - (bash/python)

bash script to copy/merge entries and assets from n spaces into 1.
[merge_spaces_script copy](merge_spaces_script copy)

## Migration script - (bash)

bash script to copy entries and assets between spaces

[migration_script](migration_script)

### How to use

#### Shell script

$./run.sh CMA_TOKEN ORIGIN_SPACE_ID DESTINATION_SPACE_ID

#### Shell script with Docker

The only difference with the shell script, is that will invoke the shell script from inside a ubuntu docker container.

$./run_docker.sh CMA_TOKEN ORIGIN_SPACE_ID DESTINATION_SPACE_ID

## Webhooks scripts

### Copy ContentTypes between spaces - (Golang or Ruby)

Script that listen to POST request (webhook) from Contenful, will copy a content types between spaces.

[Golang contentTypeBetweenSpaces](webhooks/golang/contentTypeBetweenSpaces)
[ruby contentTypeBetweenSpaces](webhooks/ruby)

### On request post the webhook payload into RabbitMQ  - (Golang)

Script that listen to POST request (webhook) from Contenful, will post the webhook paylod into RMQ.

[rmqConsumerProducer](webhooks/golang/rmqConsumerProducer)

# Comming on...

# error log anaylser
