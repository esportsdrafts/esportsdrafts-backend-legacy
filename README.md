# e-Fantasy
Daily fantasy esport leagues.

## Developing

**Requirements:**
* Docker
* Minikube
* Tilt (`brew tap windmilleng/tap && brew install windmilleng/tap/tilt`)
* kubectl (`brew install kubectl`)
* Python 3.6 (**Only** for integration tests)

Running a local dev environment:
```bash
$ tilt up
```

### Setting up your dev machine
TODO

## Building

Before building any you need to build the base docker image. This image
contains global dependencies from the `libs` directory as well as vendored
Go dependencies. Note that you only need to rebuild this image if you change
anything in the `libs` directory or add another external dependency
(if you do, remember to run `go mod vendor` from the root of the project).

To build the base image:

```bash
$ make docker-base
```

Docker is used for all building in the repo. To build a service run:

```bash
$ make SERVICE_NAME
```
Where `SERVICE_NAME` is replaced with the name of a directory in the `services`
directory. For example building the `auth` services is a matter of running:

```bash
$ make auth
```

You can also build all services by running:

```bash
$ make services
```

## Tests

Unit tests are stored side-by-side with Go source. Running them locally can be
done through:

```bash
$ make tests
```

### Integration tests

Integration tests are written in `Python` and stored in the `tests/` directory.
Any code that interfaces with the API and can be shared across tests, add it to the modules in
the `common` directory. Eventually this would become a Python SDK for the
project. An area of research is if this part can be auto-generated from the
openAPI specs.

Since these tests are written in Python you need to install a few python
dependencies. `python3.6 -m pip install -r requirements.txt`

Running all integrations tests against local environment(Need to run `Tilt`
before this step):

```bash
$ make integration-tests
```

If you wanna run integration tests against any other environment (`dev`,
`stage`, `prod`), you can configure that with the `ENVIRONMENT` variable, for example running against `stage`:

```bash
$ ENVIRONMENT=stage make integration-tests
```

Manual invocation(please see `pytest` docs) would be something like:

```bash
$ python3.6 -m pytest -vx -s --env local tests/
```

### Adding/updating Python dependencies
Managed by `pip-tools`, which means you should never directly edit the
`requirements.txt` file. To add/remove a dependency edit the
`requirements.in` file and then 'compile' it into a `requirements.txt` file by
running:

```bash
$ pip-compile -o requirements.txt requirements.in
```

To update dependencies, throw a `-U` flag on above command.

`pip-tools` can be installed by running:

```bash
$ python3.6 -m pip install -r requirements-dev.txt
```

## Deployment
TBD

## Help
An overview of all the make targets can be found by running `make help`.

Furthermore, more detailed documentation is available in the `docs/` directory. Feel free to it a good portion of editing love every time you read the docs.
