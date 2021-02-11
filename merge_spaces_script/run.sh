#!/bin/bash

## Before running:
# all spaces must have the same locale (the file try to migrate locaes also)
# run this script with the -e flag to spacifiy a sandbox environment, DO NOT RUN ON MASTER ENVIRONMENT

# Download locales for each space
# copy locales into destination space
# Loop in all spaces and for each genereate migration files 
# Import migrations files of each space into DESTINATION
# export entries and assets for each space
# Import entries and assets for each space  into DESTINATION

#CMA_TOKEN=$1
#SPACE_DEST=$2
#ENV_DEST=$2
#SPACE_ORIG=$3

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
echo "  -o  The Contentful space(s) of origin"
echo ""
echo "  -d  The Contentful space of destination"
echo ""
echo "  -e  The Contentful space of destination environment"
echo ""
echo "  -help  Show help"
echo ""
echo "Examples:"
echo "  ./run.sh -t CMA_TOKEN -d SPACE_DESTINATION -e SPACE_DESTINATION_ENVIRONEMTN -o SPACE_ORIGIN1,SPACE_ORIGIN2"
echo ""
    return 0
}

while getopts t:o:e:d:h: option
do
case "${option}"
in
t) CMA_TOKEN=${OPTARG};;
o) SPACE_ORIG=${OPTARG};;
e) ENV_DEST=${OPTARG};;
d) SPACE_DEST=${OPTARG};;
h)
    show_usage
    exit 1;;
\?) echo "Invalid option: -"$OPTARG"" >&2
    exit 1;;
esac
done


echo "Space dest: "$SPACE_DEST
echo "Env dest: "$ENV_DEST
echo "Space orig: "$SPACE_ORIG
#echo "token: "$CMA_TOKEN

MIGRATION_FOLDER="migrations"

# Here comma is our delimiter value
IFS="," read -ra space_csv <<< "$SPACE_ORIG"

echo "My array: ${space_csv[@]}"
echo "Number of elements in the array: ${#space_csv[@]}"

# Create migration folder
mkdir $MIGRATION_FOLDER

# PART 1 - COPY LOCALES
# python file and curl
# Before start export languages from spaces, and create in destination space
echo ""
echo "******************************************************"
echo "**** MIGRATING LOCALES "
echo "******************************************************"
for space_i in "${space_csv[@]}"
do
    curl --location --request GET "https://api.contentful.com/spaces/$space_i/environments/master/locales" \
            --header "Authorization: Bearer $CMA_TOKEN" > "migrations/locales_$space_i.json"
done

python locales.py $SPACE_ORIG

while IFS= read -r line; do
  curl --location --request POST "https://api.contentful.com/spaces/$SPACE_DEST/environments/$ENV_DEST/locales" \
        --header 'Content-Type: application/vnd.contentful.management.v1+json' \
        --header "Authorization: Bearer $CMA_TOKEN" \
        --data-raw "$line"
done < migrations/all_spaces_locales.txt

echo ""
echo "******************************************************"
echo "**** MIGRATING LOCALES DONE"
echo "******************************************************"
echo ""

# PART 2 - IMPORT EXPORT ContentTypes
# contenful-cli
for space_i in "${space_csv[@]}"
do
    echo "******************************************************"
    echo "**** EXPORTING Migration file - Space: $space_i"
    echo "******************************************************"
    
    contentful space --mt $CMA_TOKEN -s $space_i generate migration -f "$MIGRATION_FOLDER/migration_$space_i.js"

done
IFS=' '     # reset to default value after usage

# Import migration files into dest space
for space_i in "${space_csv[@]}"
do
    echo "******************************************************"
    echo "**** IMPORTING Migration file 'migration_$space_i.js' IN Space: $SPACE_DEST"
    echo "******************************************************"
    
    contentful space --mt $CMA_TOKEN -s $SPACE_DEST -e $ENV_DEST migration "$MIGRATION_FOLDER/migration_$space_i.js"

done

# PART 3 - IMPORT EXPORT Entries and Assets
# contenful-cli
for space_i in "${space_csv[@]}"
do
    echo ""
    echo "******************************************"
    echo "**** EXPORTING Entries and Assets - Space: $space_i"
    echo "******************************************"
    # result limit hardcoded to 500 entries per page
    contentful space export --mt $CMA_TOKEN --space-id $space_i --skip-roles --skip-webhooks --skip-content-model --include-drafts --include-archived --max-allowed-limit 500 --content-file $MIGRATION_FOLDER/migration_entries_assets_$space_i.json

done


# Import entries and assets into dest space
for space_i in "${space_csv[@]}"
do
    echo "******************************************************"
    echo "**** IMPORTING Entries 'migration_entries_assets_$space_i.json' IN Space: $SPACE_DEST"
    echo "******************************************************"
    
    contentful space import --mt $CMA_TOKEN --space-id $SPACE_DEST --environment-id $ENV_DEST --skip-content-publishing --skip-content-model  --content-file $MIGRATION_FOLDER/migration_entries_assets_$space_i.json
done
