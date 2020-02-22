#!/usr/bin/env bash


set -e

function run {
    echo "$(tput setaf 4)$1$(tput sgr0)"
    eval $1
}

run "operator-sdk build apirator/apirator"
run "kubectl delete pod -l name=apirator -n oas"
run "docker push apirator/apirator:latest"