# esportsdrafts
Daily fantasy esport leagues. This is the main repo for backend services. The
frontend and infrastructure code lives in separate repos within the
organization.

## Developing
The whole cluster can be run locally with auto-reload using Minikube. Below
are steps to get up and running. Install the requirements below to get started.

For `OS X` most of these can be installed via `brew`.

**Requirements**:
* Go
* Docker
* Minikube
* Virtualbox
* kubectl
* Tilt
* Python3.7 (Running integration tests)

Below are one-time setup for local development using `Minikube` and `Tilt`.
After following the guide you should be able to just run `make watch`. Each
service will be rebuilt and deployed to your local deployment automatically
when a change is detected.

1. Set up 'email inbox' for local email testing. More details on shared
   folders for Minikube: `https://minikube.sigs.k8s.io/docs/tasks/mount/`.

   **OS X:**
   ```bash
   $ sudo mkdir -p /Users/inbox
   $ sudo chmod 777 /Users/inbox
   ```

   **Linux:**
   ```bash
   $ sudo mkdir -p /home/inbox
   $ sudo chmod 777 /home/inbox
   ```

2. Start minikube:
   ```bash
   $ minikube start --memory=4096 --cpus=4 --vm-driver=virtualbox
   $ minikube addons enable ingress
   ```
   **Note:** Tweak the number of CPUs and memory based on your machine.

3. Add `/etc/hosts` entry:
   ```bash
   $ echo "$(minikube ip) api.esportsdrafts.localhost" | sudo tee -a /etc/hosts
   $ echo "$(minikube ip) esportsdrafts.localhost" | sudo tee -a /etc/hosts
   ```

4. Now everything should be ready to go. To get a local cluster up and running
   on your newly created Minikube setup, run this command from the root of the
   project:
   ```bash
   $ make watch
   ```

   Which will launch Tilt and setup your local cluster. A browser window will
   open where you can view the state of the cluster and logs for each service.

When you reboot your computer minikube could be disabled. You do not have to
go through the whole setup to get it back up.

To check the status of the cluster run:
```bash
$ minikube status
host: Stopped
kubelet:
apiserver:
kubeconfig:
```

To get it back up:
```bash
$ minikube start
```

Furthermore, you might have to configure the docker daemon by executing:
```bash
$ eval $(minikube docker-env)
```

If nothing works try recreating everything from scratch by running
`minikube delete`, and then follow the guide above again. The IP of minikube
will most likely change so make sure you delete the entries `/etc/hosts` and
add them back with the new IP.

**NOTE:** Steps have not been tested on Linux so setup process might need some
love work properly

## Committing Code
The project follows a trunk-based development process, in other words the
`master` branch should always be in a functional state and a release should
be possible at any given time. To make changes, create a new branch off master,
work on it and commit to the branch. When you are ready to get it into the main
code, create a `Pull Request` using the branch you have been working on and
have someone review and approve it. When it has been approved and all tests
have passed you will be able to merge it into master.

The codebase uses `semantic-release` to auto-generate version numbers. Which
means there has to be a way for the tool to understand the impact of a commit.
Therefore, the [Angular Commit Message Conventions](https://github.com/angular/angular.js/blob/master/DEVELOPERS.md#-git-commit-guidelines) are used when writing your commit message.

Some tools are available to help you remember the conventions:
* [commitizen](https://github.com/commitizen/cz-cli)
* [commitlint](https://github.com/conventional-changelog/commitlint)

Furthermore, our CI system is going to lint your commit message and bark if you
are doing it wrong. :kissing_closed_eyes:

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

Furthermore, more detailed documentation is available in the `docs/` directory.
Feel free to give it a good portion of editing love every time you read the
docs.
