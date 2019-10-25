# -*- mode: Python -*-

default_registry('docker.pkg.github.com/barreyo/efantasy')

services_k8s_files = [
    # GENERAL/GLOBAL CONFIG
    'certs/k8s-secrets.yaml',

    # AUTH
    'services/auth/k8s/deployment.yaml',
    'services/auth/k8s/service.yaml',

    # NOTIFICATIONS
    'services/notifications/k8s/deployment.yaml',
    'services/notifications/k8s/service.yaml',

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

    # FRONTEND
    'services/frontend/k8s/deployment.yaml',
    'services/frontend/k8s/service.yaml',
]

# Kubernetes YAML config files
k8s_yaml(services_k8s_files)

# Frontend is an edge-case since it lives in a seperate repo
docker_build('efantasy-frontend', '../efantasy-frontend/',
             dockerfile='../efantasy-frontend/Dockerfile')

# Docker images
docker_build('efantasy-base', './',
             dockerfile='Dockerfile')

docker_build('efantasy-mysql', 'services/mysql',
             dockerfile='services/mysql/Dockerfile')

docker_build('efantasy-auth', 'services/auth',
             dockerfile='services/auth/Dockerfile')

docker_build('efantasy-notifications', 'services/notifications',
             dockerfile='services/notifications/Dockerfile')

docker_build('efantasy-beanstalkd', 'services/beanstalkd',
             dockerfile='services/beanstalkd/Dockerfile')
