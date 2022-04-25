#! /bin/bash

set -e 

FILES=$(git ls-files | grep "\.sh$")
if [ "$1" == "fix" ]
then
    for FILE in $FILES; do 
        echo "checking $FILE file"
        shellcheck -f diff "$FILE" | patch "$FILE"
    done
else
    echo "$FILES" | xargs shellcheck -f tty
fi
