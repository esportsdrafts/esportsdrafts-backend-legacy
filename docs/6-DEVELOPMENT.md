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

## Committing Code
The project follows a trunk-based development process, in other words the 
`master` branch should always be in a functional state and a release should
be possible at any given time. To make changes, create a new branch off 
master, work on it and commit to the branch. When you are ready to get it
into the main code, create `Pull Request` using the branch you have been
working on and have someone review and approve it. When it has been 
approved an all tests have passed you will be able to merge it into
master.

The codebase uses `semantic-release` to auto-generate version numbers.
Which means there has to be a way for the tool to understand the impact
of a commit. Therefore, the 
[Angular Commit Message Conventions](https://github.com/angular/angular.js/blob/master/DEVELOPERS.md#-git-commit-guidelines)
are used when writing your commit message.

Some tools are available to help you remember the conventions:
* [commitizen](https://github.com/commitizen/cz-cli)
* [commitlint](https://github.com/conventional-changelog/commitlint)

Furthermore, our CI system is going to lint your commit message and
bark if you are doing it wrong. :kissing_closed_eyes:

## Adding a New Service
Take a look at the folder structure for a service in the `services/`
directory and then copy it over to your new service. Then work on the
API specification in the `schemas/` directory. When you are ready you
can generate the stubs using `oapi-gen`. The Makefile should have a
target that can do this, just remember to rename any naming variables
to match your new service.

TODO: Make a small utility that creates the boilerplate structure 
given a name of a new service.

