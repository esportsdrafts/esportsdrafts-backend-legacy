
Deploy steps:

```
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls-key.key -out tls-cert.crt

openssl dhparam -out dhparam.pem 2048

echo "Created SSL cert"

echo "start minikube"
minikube start

kubectl create secret tls tls-certificate --key tls-key.key --cert tls-cert.crt --namespace default

kubectl create secret generic tls-dhparam --from-file=dhparam.pem --namespace default

kubectl create -f .

echo "Appending minikube service IP to /etc/hosts"

echo "$(minikube ip) api.sample.com" | sudo tee -a /etc/hosts
```
