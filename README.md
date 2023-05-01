# Example 1

Build the counter app:

```bash
docker buildx build -t danijelradakovic/counter \
  --target release-debian \
  --build-arg SRC=cmd/counter/main.go .
```
Go to kube/example-1 directory (all commands should be executed from this directory)
```shell
cd kube/example-1
```
Create `demo` namespace:
```shell
kubectl create namespace demo
```
Deploy Pod:
```shell
kubectl -n demo apply -f pod.yaml
```
Deploy Deployment:
```shell
kubect -n demo apply -f deployment.yaml
```

Testing load balancing:
```shell
kubectl -n demo run -it --rm  --image curlimages/curl:8.00.1 curl -- sh
```
Inside the container execute several times `curl http://counter:8000`
```shell
/ $ curl http://counter:8000
Counter:  10
/ $ curl http://counter:8000
Counter:  11
/ $ curl http://counter:8000
Counter:  8
/ $ curl http://counter:8000
Counter:  9
/ $ curl http://counter:8000
Counter:  12
/ $ curl http://counter:8000
Counter:  13
/ $ curl http://counter:8000
Counter:  14
/ $ curl http://counter:8000
Counter:  15
/$ curl http://counter:8000
Counter:  10
```

Testing connection from another namespace:
```shell
kubectl create namespace tmp
kubectl -n demo run -it --rm  --image curlimages/curl:8.00.1 curl -- sh

/ $ curl http://counter.demo.svc.cluter.local:8000
Counter:  16
```

Ingress setup:

Create echo pod and service:
```shell
kubectl -n demo apply -f echo-server.yaml
```
Deploy ingress:
```shell
kubectl -n demo apply -f ingress.yaml
```

# Example 2
```shell
kubectl -n demo apply -R -f kube/example-2
# Destroy
kubectl -n demo delete -R -f kube/example-2
```

# Example 3

```shell
export INGRESS_HOST=$(minikube ip) 

# install echo server
helm upgrade --install \
    echo helm/app-1 \
    --namespace demo \
    --create-namespace \
    --values kube/example-3/echo-values.yaml \
    --set "ingress.hosts[0].host=$INGRESS_HOST.nip.io" \
    --wait

# install counter app
helm upgrade --install \
    counter helm/app-1 \
    --namespace demo \
    --create-namespace \
    --values kube/example-3/counter-values.yaml \
    --set "ingress.hosts[0].host=$INGRESS_HOST.nip.io" \
    --wait
```
# Example 4
Build the image:
```shell
docker buildx build -t danijelradakovic/dojo \
  --target release-alpine \
  --build-arg SRC=cmd/dojo/main.go .
```

Create postgresql database:
```shell
helm repo add bitnami https://charts.bitnami.com/bitnami

helm upgrade --install \
  postgresql bitnami/postgresql \
  --namespace postgresql \
  --create-namespace
```
Execute init script:
```shell
export POSTGRES_PASSWORD=$(kubectl get secret --namespace postgresql postgresql -o jsonpath="{.data.postgres-password}" | base64 -d)
kubectl run postgresql-client --rm --tty -i --restart='Never' \
  --namespace postgresql \
  --image docker.io/bitnami/postgresql:15.2.0-debian-11-r26 \
  --env="PGPASSWORD=$POSTGRES_PASSWORD" \
  --command -- \
  psql --host postgresql -U postgres -d postgres -p 5432 -c \
  'CREATE TABLE weapons(id text not null, name text not null, PRIMARY KEY(id));'
```
Deploy dojo app:
```shell
export INGRESS_HOST=$(minikube ip) 

# install dojo app
helm upgrade --install \
    dojo helm/app-2 \
    --namespace dojo \
    --create-namespace \
    --values kube/example-4/dojo-values.yaml \
    --set "ingress.hosts[0].host=$INGRESS_HOST.nip.io" \
    --wait
```
Test app:
```shell
curl -X POST "http://$INGRESS_HOST.nip.io/weapon?id=0&weapon=katana"
curl -X POST "http://$INGRESS_HOST.nip.io/weapon?id=1&weapon=ninjaStar"
curl -X POST "http://$INGRESS_HOST.nip.io/weapon?id=2&weapon=ninjaSword"

curl "http://$INGRESS_HOST.nip.io/weapon"
# [{"id":"0","name":"katana"},{"id":"1","name":"ninjaStar"},{"id":"2","name":"ninjaSword"}]
```

Delete the postgresql:
```shell
helm -n postgresql uninstall postgresql
kubectl -n postgresql get pvc
kubectl -n postgresql pvc data-postgresql-0
kubectl get pv -A
kubectl delete pv <name of pv for postgresql>
```

# Example 5
```shell
export INGRESS_HOST=$(minikube ip) # ip address of minikube node

helm repo add prometheus https://prometheus-community.github.io/helm-charts
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
```

Install Prometheus, Alertmanager, Pushgateway, Node Exoprter, Kube State Metrics:
```shell
helm upgrade --install \
    prometheus prometheus/prometheus \
    --namespace monitoring \
    --create-namespace \
    --values kube/example-5/prometheus-values.yaml \
    --set "server.ingress.hosts[0]=prometheus.$INGRESS_HOST.nip.io" \
    --set "alertmanager.ingress.hosts[0]=alertmanager.$INGRESS_HOST.nip.io" \
    --wait

echo "Prometheus URL: http://prometheus.$INGRESS_HOST.nip.io"
echo "Alertmanager URL: http://alertmanager.$INGRESS_HOST.nip.io"
```
Install Promtail and Loki:
```shell
helm upgrade --install \
    loki grafana/loki-stack \
    --namespace monitoring \
    --create-namespace \
    --wait
```
Install Grafana:
```shell
helm upgrade --install \
    grafana grafana/grafana \
    --namespace monitoring \
    --create-namespace \
    --values kube/example-5/grafana-values.yml \
    --set "ingress.hosts[0]=grafana.$INGRESS_HOST.nip.io" \
    --wait

password=$(kubectl --namespace monitoring \
    get secret grafana \
    --output jsonpath="{.data.admin-password}" \
    | base64 --decode)

echo "Grafana URL: http://grafana.$INGRESS_HOST.nip.io"
echo "Username: admin"
echo "Password: $password"
```
Install Certificate Manager:
```shell
helm repo add jetstack https://charts.jetstack.io
helm repo update

helm upgrade --install \
    cert-manager jetstack/cert-manager \
    --namespace cert-manager \
    --create-namespace \
    --set installCRDs=true \
    --wait
```
Install Jaeger operator and Jaeger:
```shell
kubectl create namespace observability 
kubectl --namespace observability apply \
    --filename https://github.com/jaegertracing/jaeger-operator/releases/download/v1.39.0/jaeger-operator.yaml

sed  "s|\$INGRESS_HOST|$INGRESS_HOST|g" kube/example-5/jaeger.yaml | kubectl --namespace observability apply --filename -
echo "Jaeger URL: http://jaeger.$INGRESS_HOST.nip.io"
```

# Example 6
Preparation:
```shell
export INGRESS_HOST=$(minikube ip) # ip address of minikube node

helm repo add prometheus https://prometheus-community.github.io/helm-charts
helm repo add grafana https://grafana.github.io/helm-charts
helm repo add jetstack https://charts.jetstack.io
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
```

Install Promtail and Loki:
```shell
helm upgrade --install \
    loki grafana/loki-stack \
    --namespace monitoring \
    --create-namespace \
    --wait
```

Install Certificate Manager:
```shell
helm repo add jetstack https://charts.jetstack.io
helm repo update

helm upgrade --install \
    cert-manager jetstack/cert-manager \
    --namespace cert-manager \
    --create-namespace \
    --set installCRDs=true \
    --wait
```
Install Jaeger operator and Jaeger:
```shell
kubectl create namespace observability 
kubectl --namespace observability apply \
    --filename https://github.com/jaegertracing/jaeger-operator/releases/download/v1.39.0/jaeger-operator.yaml

sed  "s|\$INGRESS_HOST|$INGRESS_HOST|g" kube/example-6/jaeger.yaml | kubectl --namespace observability apply --filename -
echo "Jaeger URL: http://jaeger.$INGRESS_HOST.nip.io"
```

Install `kube-prometheus-stack`:
```shell
helm upgrade --install \
    monitoring prometheus/kube-prometheus-stack \
    --namespace monitoring \
    --create-namespace \
    --values kube/example-6/monitoring-values.yaml \
    --set "alertmanager.ingress.hosts[0]=alertmanager.$INGRESS_HOST.nip.io" \
    --set "prometheus.ingress.hosts[0]=prometheus.$INGRESS_HOST.nip.io" \
    --set "grafana.ingress.hosts[0]=grafana.$INGRESS_HOST.nip.io" \
    --wait
```

```shell
password=$(kubectl --namespace monitoring \
    get secret monitoring-grafana \
    --output jsonpath="{.data.admin-password}" \
    | base64 --decode)

echo "Prometheus URL: http://prometheus.$INGRESS_HOST.nip.io"
echo "Alertmanager URL: http://alertmanager.$INGRESS_HOST.nip.io"
echo "Grafana URL: http://grafana.$INGRESS_HOST.nip.io"
echo "Username: admin"
echo "Password: $password"
```

Create postgresql database:
```shell
helm upgrade --install \
  postgresql bitnami/postgresql \
  --namespace postgresql \
  --create-namespace
```
Execute init script:
```shell
export POSTGRES_PASSWORD=$(kubectl get secret --namespace postgresql postgresql -o jsonpath="{.data.postgres-password}" | base64 -d)
kubectl run postgresql-client --rm --tty -i --restart='Never' \
  --namespace postgresql \
  --image docker.io/bitnami/postgresql:15.2.0-debian-11-r26 \
  --env="PGPASSWORD=$POSTGRES_PASSWORD" \
  --command -- \
  psql --host postgresql -U postgres -d postgres -p 5432 -c \
  'CREATE TABLE weapons(id text not null, name text not null, PRIMARY KEY(id));'
```

Deploy dojo app:
```shell
helm upgrade --install \
    dojo helm/app-3 \
    --namespace dojo \
    --create-namespace \
    --values kube/example-6/dojo-values.yaml \
    --set "ingress.hosts[0].host=$INGRESS_HOST.nip.io" \
    --wait
```

Test app:
```shell
curl -X POST "http://$INGRESS_HOST.nip.io/weapon?id=0&weapon=katana"
curl -X POST "http://$INGRESS_HOST.nip.io/weapon?id=1&weapon=ninjaStar"
curl -X POST "http://$INGRESS_HOST.nip.io/weapon?id=2&weapon=ninjaSword"

curl "http://$INGRESS_HOST.nip.io/weapon"
# [{"id":"0","name":"katana"},{"id":"1","name":"ninjaStar"},{"id":"2","name":"ninjaSword"}]
```
