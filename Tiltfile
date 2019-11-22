# -*- mode: Python -*-

default_registry('docker.pkg.github.com/esportsdrafts/esportsdrafts')

services_k8s_files = [
    # GENERAL/GLOBAL CONFIG
    'config/certs/k8s-secrets.yaml',
    'config/configmaps/dev.yaml',
    'config/secrets/dev-auth.yaml',

    # INGRESS
    'services/ingress/roles.yaml',
    'services/ingress/default-backend.yaml',
    'services/ingress/ingress-deployment.yaml',
    'services/ingress/ingress.yaml',
    'services/ingress/service.yaml',

    # MYSQL (MARIADB)
    'services/mysql/k8s/pv.yaml',
    'services/mysql/k8s/deployment.yaml',
    'services/mysql/k8s/service.yaml',

    # BEANSTALKD
    'services/beanstalkd/k8s/pdb.yaml',
    'services/beanstalkd/k8s/pvc.yaml',
    'services/beanstalkd/k8s/statefulset.yaml',
    'services/beanstalkd/k8s/service.yaml',

    # AUTH
    'services/auth/k8s/deployment.yaml',
    'services/auth/k8s/service.yaml',

    # NOTIFICATIONS
    'services/notifications/k8s/deployment.yaml',
    'services/notifications/k8s/service.yaml',

    # FRONTEND
    'services/frontend/k8s/deployment.yaml',
    'services/frontend/k8s/service.yaml',
]

# Kubernetes YAML config files
k8s_yaml(services_k8s_files)

go_ignores = ['tests', 'certs', 'docs', 'Dockerfile.testing', 'README.md',
              'requirements-dev.txt', 'requirements.in', 'requirements.txt',
              'config', 'scripts', 'Makefile']

# Core services
docker_build('esportsdrafts-mysql', 'services/mysql',
             dockerfile='services/mysql/Dockerfile')

docker_build('esportsdrafts-beanstalkd-metrics', 'services/beanstalkd',
             dockerfile='services/beanstalkd/Dockerfile.metrics')

docker_build('esportsdrafts-beanstalkd', 'services/beanstalkd',
             dockerfile='services/beanstalkd/Dockerfile')

# App Services
docker_build('esportsdrafts-base', './',
             dockerfile='Dockerfile', ignore=go_ignores)

docker_build('esportsdrafts-auth', 'services/auth',
             dockerfile='services/auth/Dockerfile',
             ignore=go_ignores + ['services/notifications', 'services/frontend'])

docker_build('esportsdrafts-notifications', 'services/notifications',
             dockerfile='services/notifications/Dockerfile',
             ignore=go_ignores + ['services/auth', 'services/frontend'])

# Frontend is an edge-case since it lives in a separate repo
docker_build('esportsdrafts-frontend', '../esportsdrafts-frontend/',
             dockerfile='../esportsdrafts-frontend/Dockerfile')
