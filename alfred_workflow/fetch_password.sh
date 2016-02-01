#!/bin/bash

query=`cat`
password=`security find-generic-password -l "$query" -w`
code=$?
if [ $code -eq 44 ]; then
  exit 1
elif [ $code -ne 0 ]; then
  exit 2
else
  echo "$password" | tee >(bash ./erase_later.sh) >(pbcopy) > /dev/null
fi

