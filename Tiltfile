# -*- mode: Python -*-
default_registry(read_json('tilt_option.json', {})
                 .get('default_registry', 'gcr.io/windmill-public-containers/servantes'))

services_k8s_files = [
    # GENERAL/GLOBAL CONFIG
    'certs/k8s-secrets.yaml',

    # AUTH
    'services/auth/k8s/deployment.yaml',
    'services/auth/k8s/service.yaml',

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

    # FRONTEND
    'services/frontend/k8s/deployment.yaml',
    'services/frontend/k8s/service.yaml',
]

# Kubernetes YAML config files
k8s_yaml(services_k8s_files)

# Docker images
docker_build('efantasy-mysql', 'services/mysql',
             dockerfile='services/mysql/Dockerfile')

# Live updates in dev mode
docker_build('efantasy-auth', 'services/auth',
             dockerfile='services/auth/Dockerfile')

docker_build('efantasy-frontend', '../efantasy-frontend',
             dockerfile='../fantasy-frontend/Dockerfile')
