# Testing
The project utilizes two mechanisms for testing: unit tests and integration
tests. Most go functions should have unit tests that test them in isolation and
run very fast. Meanwhile integation tests are written as an external process
that calls the API endpoints and verifies the behavior of a user is correct.

## Unit Testing
Standard GO unittests.

Run all of them:
```bash
$ make tests
```

## Integration Tests
Integration tests are written in `Python` and stored in the `tests/` directory.
Any code that interfaces with the API and can be shared across tests, add it to
the modules in the `common` directory. Eventually this would become a Python
SDK for the project. An area of research is if this part can be auto-generated
from the openAPI specs.

Since these tests are written in Python you need to install a few python
dependencies.
```bash
$ python3 -m pip install -r requirements.txt
```

Furthermore, in development please use Python linters and type checkers, which
can be installed through the `requirements-dev.txt` file.
```bash
$ python3 -m pip install -r requirements-dev.txt
```

Running all integration tests against local environment(Need to run Tilt
(`make watch`) in a different shell before this step):
```bash
$ make integration-tests
```

If you wanna run integration tests against any other environment (`staging`,
`production`), you can configure that with the `ENVIRONMENT` variable, for
example running against `staging`:

```bash
$ ENVIRONMENT=staging make integration-tests
```

Manual invocation(please see `pytest` docs) would be something like:
```bash
$ python3 -m pytest -vx -s --env local tests/
```

### Adding/updating Python dependencies
Managed by `pip-tools`, which means you should never directly edit the
`requirements.txt` file. To add/remove a dependency edit the `requirements.in`
file and then 'compile' it into a `requirements.txt` file by running:
```bash
$ pip-compile -o requirements.txt requirements.in
```

To update dependencies, throw a `-U` flag on above command.

`pip-tools` can be installed by running:

```bash
$ python3 -m pip install -r requirements-dev.txt
```

