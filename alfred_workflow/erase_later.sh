#!/bin/bash

content_to_erase=`cat`

function conditionally_erase_clipboard_later() {
  sleep 10
  current_content=`pbpaste`

  if [ "$current_content" == "$content_to_erase" ]; then
    echo "" | pbcopy
  fi
}

conditionally_erase_clipboard_later </dev/null >/dev/null 2>&1 &
disown

