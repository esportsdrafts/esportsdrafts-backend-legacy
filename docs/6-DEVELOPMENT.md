# Development
Local development depens on a few things. Install the requirements 
below to get started. Rest of this document relies on them being
installed.

For `OS X` these can be installed via `brew`.

**Requirements**: 
* Go 
* Docker
* Minikube
* kubectl
* Tilt
* Python3.7

## Configuring Local Env
To be written

## Adding a New Service
Take a look at the folder structure for a service in the `services/`
directory and then copy it over to your new service. Then work on the
API specification in the `schemas/` directory. When you are ready you
can generate the stubs using `oapi-gen`. The Makefile should have a
target that can do this, just remember to rename any naming variables
to match your new service.

