# -*- mode: Python -*-
default_registry(read_json('tilt_option.json', {})
                 .get('default_registry', 'gcr.io/windmill-public-containers/servantes'))

services_k8s_files = [
    'certs/k8s-secrets.yaml',
    'services/auth/k8s/deployment.yaml',
    'services/auth/k8s/service.yaml',
    'services/ingress/roles.yaml',
    'services/ingress/default-backend.yaml',
    'services/ingress/ingress-deployment.yaml',
    'services/ingress/ingress.yaml',
    'services/ingress/service.yaml',
    'services/mysql/k8s/pv.yaml',
    'services/mysql/k8s/deployment.yaml',
    'services/mysql/k8s/service.yaml',
    'services/frontend/k8s/deployment.yaml',
    'services/frontend/k8s/service.yaml',
]

# Kubernetes YAML config files
k8s_yaml(services_k8s_files)

docker_build('efantasy-frontend', '../efantasy-frontend/',
             dockerfile='../efantasy-frontend/Dockerfile')

# Docker images
docker_build('efantasy-mysql', 'services/mysql',
             dockerfile='services/mysql/Dockerfile')

# Live updates in dev mode
docker_build('efantasy-auth', 'services/auth',
             live_update=[
                 sync('services/auth', '/workspace'),
                 run("cd /workspace/cmd/ && CGO_ENABLED=0 go build -installsuffix 'static' -o /app"),
                 restart_container(),
             ]
             )
