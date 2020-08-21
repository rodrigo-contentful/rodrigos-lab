#!/bin/bash

## Before running:
# spaces must have the same locale


CMA_TOKEN=$1
SPACEA=$2
SPACEB=$3

MIGRATION_FOLDER="migrations"


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

mkdir $MIGRATION_FOLDER

echo "Starting exports, saving into '$MIGRATION_FOLDER' folder "


#contentful login --mt $CMA_TOKEN
echo "******************************************"
echo "**** EXPORTING Content Model Schema ******"
echo "******************************************"
contentful space export --mt $CMA_TOKEN --space-id $SPACEA --skip-roles --skip-webhooks --skip-content --include-drafts --content-file $MIGRATION_FOLDER/OnlyCMTest.json

echo "******************************************"
echo "**** EXPORTING Content Entries ******"
echo "******************************************"
# result limit hardcoded to 500 entries per page
contentful space export --mt $CMA_TOKEN --space-id $SPACEA --skip-roles --skip-webhooks --skip-content-model --include-drafts --max-allowed-limit 500 --query-assets 'sys.id=none' --content-file $MIGRATION_FOLDER/OnlyEntriesTest.json

echo "******************************************"
echo "**** EXPORTING Content Assets ******"
echo "******************************************"
contentful space export --mt $CMA_TOKEN --space-id $SPACEA --skip-roles --skip-webhooks --skip-content-model --include-drafts --max-allowed-limit 500  --query-entries 'sys.id=none' --content-file $MIGRATION_FOLDER/OnlyAssetsTest.json
#contentful logout


echo "******************************************"
echo "**** COPYING Content Model Schema ******"
echo "******************************************"
contentful space import --mt $CMA_TOKEN --space-id $SPACEB --skip-content-publishing --content-model-only  --content-file $MIGRATION_FOLDER/OnlyCMTest.json

echo "******************************************"
echo "**** COPYING Content Entries ******"
echo "******************************************"
contentful space import --mt $CMA_TOKEN --space-id $SPACEB --skip-content-publishing --skip-content-model  --content-file $MIGRATION_FOLDER/OnlyEntriesTest.json

echo "******************************************"
echo "**** COPYING Content Assets ******"
echo "******************************************"
contentful space import --mt $CMA_TOKEN --space-id $SPACEB --skip-content-publishing --skip-content-model  --content-file $MIGRATION_FOLDER/OnlyAssetsTest.json

# cleanup
rm -r $MIGRATION_FOLDER