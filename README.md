[![Build Status][travis-badge]][travis-url]

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

[travis-badge]: https://travis-ci.org/apirator/apirator.svg?branch=master
[travis-url]: https://travis-ci.org/apirator/apirator

# APIrator

APIMocks for developers made easy!


## Development commands

### Build

We are using [operator-sdk](https://github.com/operator-framework/operator-sdk) to manage our lifecycle, you should have the operator-sdk installed.
On the root folder you can use the following command to build and create the docker image

````shell script
operator-sdk build apirator/apirator
````

### Deploy

These instructions should be user for **development** purpose, for **production** usage we 
recommend you to use [apirator-ops](https://github.com/apirator/apirator-ops) repository that contains our **helm** installation. 

#### CRD (Custom Resource Definition)

There is a simple custom resource definition that define our main object a Open API 
Specification object. To install it on kubernetes you should execute the following
command

````shell script
kubectl apply -f <YOUR_GOLANG_PROJECT_FOLDER>/apirator/deploy/crds/apirator.io_apimocks_crd.yaml -n oas
````

#### Deployment

You should deploy the files located on /deploy folder. There are files about
ServiceAccount, RoleBindings, Roles and deployment.yaml. The apirator is a simple
pod that implemented following the operator pattern.

To deploy the apirator you should execute the following command

````shell script
kubectl apply -f <YOUR_GOLANG_PROJECT_FOLDER>/apirator/deploy/operator.yaml -n oas
````


###References:

You can find more instructions about operator-sdk command line [here](https://docs.openshift.com/container-platform/4.1/applications/operator_sdk/osdk-getting-started.html)
