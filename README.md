# go-server-test
Incldes script that can be triggered by a webhook to copy contentypes and entries between n.. spaces.

* Includes a shell scriot that simplified the copy of spaces using contenful cli

* Includes a docker that simplified the copy of spaces using contenful cli

## Golang 

Scrip that can be triggered with a webhook and copy a content type and entry between spaces

### Min requirements

go 1.14

### How to run

$ go run server.go

### How to use

POST localhost:80


## Ruby

Scrip that can be triggered with a webhook and copy a content type and entry between spaces

## Bash script (and docker)

Shell scriot that simplified the copy of spaces using contenful cli, can be invoked using docker.

### How to use

#### Shell script

$./run.sh CMA_TOKEN ORIGIN_SPACE_ID DESTINATION_SPACE_ID


#### Shell script with Docker

The only difference with the shell script, is that will invoke the shell script from inside a ubuntu docker container.

$./run_docker.sh CMA_TOKEN ORIGIN_SPACE_ID DESTINATION_SPACE_ID
