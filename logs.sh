#!/usr/bin/env bash

set -e

SINCE=$1
if [[ -z "$SINCE" ]]
then
    SINCE="10m"
fi

function run {
    echo "$(tput setaf 2)$1$(tput sgr0)"
    eval $1
}

run "kubectl logs --since=$SINCE -f -l app=apirator -n oas | jq -Rc 'fromjson? | select(type == \"object\")'"