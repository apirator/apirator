#!/usr/bin/env bash

set -e

function run {
    echo "$(tput setaf 6)$1$(tput sgr0)"
    eval $1
}

VERSION=$(grep -o '".*"' ./version/version.go | sed 's/"//g')
run "operator-sdk build apirator/apirator:$VERSION"
run "docker push apirator/apirator:$VERSION"

DEPLOYMENT=$(kubectl get deployment -l app=apirator -n oas -o=name)
if [[ -z "$DEPLOYMENT" ]]
then
    run "cat deploy/operator.yaml | sed \"s/image:\ apirator\/apirator/image:\ apirator\/apirator:$VERSION/\" | kubectl -n oas apply -f -"
else
    run "kubectl set image deployment.extensions/apirator apirator=apirator/apirator:$VERSION --record -n oas"
fi