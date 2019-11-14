# Development
Local development depends on a few things. Install the requirements below to
get started. Rest of this document assumes them already being installed.

For `OS X` most of these can be installed via `brew`.

**Requirements**:
* Go
* Docker
* Minikube
* Virtualbox
* kubectl
* Tilt
* Python3.7 (Running integration tests)

## Configuring Local Development Environment
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

## Adding a New Service
Take a look at the folder structure for a service in the `services/` directory
and then copy it over to your new service. Then work on the API specification in
the `schemas/` directory. When you are ready you can generate the stubs using
`oapi-gen`. The Makefile should have a target that can do this, just remember to
rename any naming variables to match your new service.

TODO: Make a small utility that creates the boilerplate structure given a name
of a new service.
