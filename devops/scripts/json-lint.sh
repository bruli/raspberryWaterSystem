#!/usr/bin/env bash

for i in $(find . -type f -name '*.json')
do
  echo "checking file $i"
  go tool jv $i
  if [ $? == 1 ]; then
    echo "invalid json format in $i"
    exit 1
  fi
done
