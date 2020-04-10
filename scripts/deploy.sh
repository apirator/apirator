#!/usr/bin/env bash

set -e

function run {
    echo "$(tput setaf 6)$1$(tput sgr0)"
    eval $1
}

COMMIT_HASH=$(git log --oneline -1 --pretty="format:%h")
run "operator-sdk build apirator/apirator:$COMMIT_HASH"
run "docker push apirator/apirator:$COMMIT_HASH"

DEPLOYMENT=$(kubectl get deployment -l app=apirator -n oas -o=name)
if [[ -z "$DEPLOYMENT" ]]
then
    run "cat deploy/operator.yaml | sed \"s/image:\ apirator\/apirator/image:\ apirator\/apirator:$COMMIT_HASH/\" | kubectl -n oas apply -f -"
else
    run "kubectl set image deployment.extensions/apirator apirator=apirator/apirator:$COMMIT_HASH --record -n oas"
fi