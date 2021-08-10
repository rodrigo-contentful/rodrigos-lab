#!/bin/bash

set -e

## IMPORTANT Before running:
# This script is offer as it is with no responsability.
# all spaces must have the same locale (the file try to migrate locaes also)

function show_usage(){
echo ""
echo "Migrate Contentful content typea and entries between space(s)."
echo ""
echo "The order of migrations are: 1.- Copy locales from origin spaces to destintion, 2.- Import contenty types from origin spaces to destintion, 3.- Import entries/assets from origin spaces to destintion"
echo "Usage: ./run.sh [args]"
echo ""
echo "Options:"
echo "  -t  The Contentful management token to use"
echo ""
echo "  -o  The Origin Contentful space"
echo ""
echo "  -d  The Destination Contentful space"
echo ""
echo "  -c  The list of types to migrate as csv"
echo ""
echo "  -l  Option to migrate locales: yes/no. "
echo ""
echo "  -h  Show help"
echo ""
echo "Examples:"
echo "  ./run.sh -t CMA_TOKEN -o SPACE_ORIGIN -d SPACE_DESTINATION -c CSV_OF_CONTENTTYPES_IDS -l no"
echo ""
    return 0
}

ENV_DEST=master
MIG_LOCALES=yes
while getopts e:t:o:d:c:l:h: option
do
case "${option}"
in
t) CMA_TOKEN=${OPTARG};;
o) SPACE_ORIG=${OPTARG};;
d) SPACE_DEST=${OPTARG};;
e) ENV_DEST=${OPTARG};;
c) CSV_TYPES=${OPTARG};;
l) MIG_LOCALES=${OPTARG};;
h)
    show_usage
    exit 1;;
\?) echo "Invalid option: -$OPTARG" >&2
    exit 1;;
esac
done

echo "Space orig: $SPACE_ORIG"
echo "Space dest: $SPACE_DEST"
echo "Space dest Env: $ENV_DEST"
echo "Types to migrate: $CSV_TYPES"
echo "Migrate locales: $MIG_LOCALES"

if [ "$MIG_LOCALES" = 'no' ]; then
    echo "Will NOT migrate locales"
else 
    
    # PART 1 - COPY LOCALES
    # python file and curl
    # Before start export languages from spaces, and create in destination space
    echo ""
    echo "******************************************************"
    echo "**** MIGRATING LOCALES "
    echo "******************************************************"

    curl --location --request GET "https://api.contentful.com/spaces/$SPACE_ORIG/environments/master/locales" \
                --header "Authorization: Bearer $CMA_TOKEN" > "locales_$SPACE_ORIG.json"

    python locales.py "$SPACE_ORIG"

    while IFS= read -r line; do
    curl --location --request POST "https://api.contentful.com/spaces/$SPACE_DEST/environments/$ENV_DEST/locales" \
            --header 'Content-Type: application/vnd.contentful.management.v1+json' \
            --header "Authorization: Bearer $CMA_TOKEN" \
            --data-raw "$line"
    done < all_spaces_locales.txt

fi

echo ""
echo "******************************************************"
echo "**** GETTING SPACE CONTENT TYPES "
echo "******************************************************"

mkdir content_types
cd content_types

IFS=, read -ra values <<< "$CSV_TYPES"
for v in "${values[@]}"
do
    # todo, add --filename SOMEFILENAME to better contentType export
    #echo $v
    contentful space --space-id "$SPACE_ORIG" --management-token "$CMA_TOKEN" generate migration --content-type-id "$v"
done

cd ..

echo ""
echo "******************************************************"
echo "**** GETTING SPACE CONTENT TYPES ENTRIES"
echo "******************************************************"

# queries of entries are done using a parameter based on CDA search examples:
# https://www.contentful.com/developers/docs/references/content-delivery-api/#/reference/search-parameters

contentful space --space-id "$SPACE_ORIG" --management-token "$CMA_TOKEN" export --skip-roles --skip-webhooks --include-drafts  --query-entries "'sys.contentType.sys.id[in]=$CSV_TYPES'" --content-file exported_content.json

echo ""
echo "******************************************************"
echo "**** CREATE NEW SPACE CONTENT TYPES"
echo "******************************************************"

for filename in content_types/*; do
    #echo $filename
    contentful space --space-id "$SPACE_DEST" --environment-id "$ENV_DEST" migration --yes "$filename"
done

echo ""
echo "******************************************************"
echo "**** CREATE NEW SPACE entries"
echo "******************************************************"

contentful space --space-id "$SPACE_DEST" --environment-id "$ENV_DEST" import --skip-content-model --skip-content-publishing --content-file exported_content.json
