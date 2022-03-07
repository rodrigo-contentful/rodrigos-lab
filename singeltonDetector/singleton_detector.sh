#!/usr/bin/env bash

access_token=$CF_CDA_TOKEN
space_id=$CF_SPACE_ID
environment_id=$CF_ENV_ID

BASE_URL="https://cdn.contentful.com/spaces/${space_id}/environments/${environment_id}"
AUTH="Authorization: Bearer ${access_token}"

content_types=$(curl -X GET -H "${AUTH}" "${BASE_URL}/content_types?limit=1000" 2> /dev/null | jq '.items[].sys.id' | sort | sed "s/\"//g")
for ct in $content_types; do
  total=$(curl -X GET -H "${AUTH}" "${BASE_URL}/entries?limit=0&content_type=${ct}" 2> /dev/null | jq '.total')
  if [[ $total -le 1 ]]; then
    echo "Content Type: ${ct} - has ${total} entries. It's either a singleton or unused."
  fi
done
