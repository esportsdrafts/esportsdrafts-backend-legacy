#!/bin/bash

# Check if minikube is alive and configured
MINIKUBE_STATUS=$(minikube status | grep host: | awk '{print $2}')
if [ "${MINIKUBE_STATUS}" != "Running" ]; then
    echo "Minikube is not running, run 'minikube start' and try again..."
    echo "Please refer to /docs/6-DEVELOPMENT.md for more details on how to configure your dev environment"
    exit 1
fi

IP=$(awk '/esportsdrafts.localhost/ {print $1}' /etc/hosts | head -n1)
MINIKUBE_IP=$(minikube ip)

# Fix /etc/hosts with minikube IP
# TODO: Figure out better solution that does not require editing hosts file
if [ "${IP}" != "${MINIKUBE_IP}" ]; then
    echo "Minikube IP is not matching /etc/hosts entries. Updating..."
    sudo sed -i.bak '/esportsdrafts.localhost/d' /etc/hosts
    echo "${MINIKUBE_IP} api.esportsdrafts.localhost" | sudo tee -a /etc/hosts
    echo "${MINIKUBE_IP} esportsdrafts.localhost" | sudo tee -a /etc/hosts
fi
