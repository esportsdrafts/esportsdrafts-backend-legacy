# Deployment
Deployments to `staging` happens automatically by Jenkins each time a Pull
Request is merged into `master`. If the integration tests pass against
`staging` a deployment can be manually approved in Jenkins.

Section is TBD.

## Releases
The CI system will automatically generate a new tag, change list and create a
release whenever anything is merged to master. This is based off the tool
`semantic-release`, which parses commit messages to understand what kind of
changes have been committed to the codebase. Please see the `6-DEVELOPMENT` doc
for more details on `semantic-releases` and committing code correctly.
