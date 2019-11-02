#!/bin/bash

IP=$(awk '/esportsdrafts.localhost/ {print $1}' /etc/hosts | head -n1)
MINIKUBE_IP=$(minikube ip)

if [ "${IP}" != "${MINIKUBE_IP}" ]; then
    echo "Minikube IP is not matching /etc/hosts entries. Updating..."
    sudo sed -i.bak '/esportsdrafts.localhost/d' /etc/hosts
    echo "${MINIKUBE_IP} api.esportsdrafts.localhost" | sudo tee -a /etc/hosts
    echo "${MINIKUBE_IP} esportsdrafts.localhost" | sudo tee -a /etc/hosts
fi

