#! /bin/bash

set -e -x

# shellcheck source=devops/scripts/util.sh
source "$(dirname "$0")"/util.sh
require gojsonschema

TARGET_DIRECTORY="$(pwd)/internal/infra/http"

count_generated=$(find "$TARGET_DIRECTORY"/*_generated.go 2>/dev/null | wc -l)
if [ "$count_generated" != 0 ]
then 
  rm "$TARGET_DIRECTORY"/*_generated.go
fi

for req_file in internal/infra/http/schemas/*_response.json
do
  file="$(basename "$req_file")"
  generated="${file%.json}"
  gojsonschema  -p http -o "$TARGET_DIRECTORY"/"$generated"_generated.go "$req_file"
done

rm -rf "$TMP_DIR"
