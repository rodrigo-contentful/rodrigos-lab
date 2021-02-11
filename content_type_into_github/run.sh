#!/bin/bash

## Before running:
# all spaces must have the same locale (the file try to migrate locaes also)
# run this script with the -e flag to spacifiy a sandbox environment, DO NOT RUN ON MASTER ENVIRONMENT

# sleep for 10 seconds (to giv chance to snapshot to be created)
# get a list of content types
# for each contentType get a snapshot
# push snapshot file into github

#CMA_TOKEN=$1
#SPACE=$2
#EMVIRONMENT=$3

function show_usage(){
echo ""
echo "Runs a migration between Contentful space(s)."
echo ""
echo "The order of migrations are: 1.- Copy locales from origin spaces to destintion, 2.- Import contenty types from origin spaces to destintion, 3.- Import entries/assets from origin spaces to destintion"
echo "Usage: ./run.sh [args]"
echo ""
echo "Options:"
echo "  -t  The Contentful management token to use"
echo ""
echo "  -o  The Contentful space"
echo ""
echo "  -e  The Contentful space environmnet"
echo ""
echo "  -help  Show help"
echo ""
echo "Examples:"
echo "  ./run.sh -t CMA_TOKEN -d SPACE -e SPACE_ENVIRONMENT"
echo ""
    return 0
}

while getopts t:o:e:d:h: option
do
case "${option}"
in
t) CMA_TOKEN=${OPTARG};;
o) SPACE=${OPTARG};;
e) SPACE_ENV=${OPTARG};;
h)
    show_usage
    exit 1;;
\?) echo "Invalid option: -"$OPTARG"" >&2
    exit 1;;
esac
done

echo "Space: "$SPACE
echo "Env: "$SPACE_ENV

echo ""
echo "******************************************************"
echo "**** GETTING SPACE CONTENTTYPES "
echo "******************************************************"


echo ""
echo "******************************************************"
echo "**** DOWNLOADING SNAPSHOTS "
echo "******************************************************"

# 
echo ""
echo "**** DOWNLOADING SNAPSHOTS FOR: aa "

echo ""
echo "******************************************************"
echo "**** PUSHING SNAPSHOTS TO GITHUB "
echo "******************************************************"
