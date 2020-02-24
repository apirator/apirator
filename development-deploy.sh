#!/usr/bin/env bash


set -e

function run {
    echo "$(tput setaf 4)$1$(tput sgr0)"
    eval $1
}

run "operator-sdk build apirator/apirator "
run "kubectl delete deployment -l app=apirator -n oas "
run "docker push apirator/apirator:latest "
run "kubectl apply -f deploy/operator.yaml -n oas "