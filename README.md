# Example 1

Build the counter app:

```bash
docker buildx build -t danijelradakovic/counter --target release-debian --build-arg SRC=cmd/counter.go .
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
