## k8s deployment overviews

haproxy-ingress pods are in the namespace `ingress` and can be inspected with

```sh
kubectl --context et-ma2-prod -n ingress get pods
```

There are 4 Daemon Sets: haproxy-bh1, haproxy-bh2, haproxy-tr1, haproxy-tr2

## To bash into a running ingress

kubectl --context et-ma2-prod -n ingress exec -it haproxy-bh2-vpkql -- bash

## To find out what haproxy-ingress produces partial results without 'syn' part

From a node running bid:

```sh
sudo nsenter -t `ps auxw | grep java | grep header-bidder | grep -v grep | head -n 1 | awk '{print $2}'` -n tcpdump -nnn -c 100000 -A dst port 8072 | grep -a1 'ppfp: {.*cap_drv_id":0.*' | grep 'x-pphn' | sort | uniq -c | sort -nr
```

## Applying changes to k8s maps without commiting into kube-manifests

From local kube-manifests working copy: 

```sh
kubectl --context et-ma2-prod apply --dry-run=server -f bh/ingresses/bh-ingress-haproxy-ing.yaml -f tr/ingresses/tr-ingress-haproxy-ing.yaml
```

# Simulate reloading

To simulate reloading that happens in haproxy-ingress, you can run

```sh
haproxy -f /etc/haproxy -p /var/run/haproxy/haproxy.pid -D -sf CURRENT_HAPROXY_PID -x /var/run/haproxy/admin.sock
```

from a haproxy-ingress pod. Where CURRENT_HAPROXY_PID needs to be substituted with a pid of the currently running haproxy master.

# To check what kubernetes think it should be running vs what it's actually running

```sh
kubectl --context et-ma2-prod -n ingress get ds haproxy-bh1 -o json | jq '.spec.template.spec.containers[0].image'
```

will return the image from the applied config map, aka "what it should be running".

```sh
for pod in $(kubectl --context et-ma2-prod -n ingress get pods | grep haproxy | awk '{print $1}'); do echo "$pod:";  kubectl --context et-ma2-prod -n ingress get pod $pod -o json | jq -r '.spec.containers[0].image'; done
```

will show the actual images pods are running

# Forcefully delete/restart all haproxy pods

```sh
for pod in $(kubectl --context et-ma2-prod -n ingress get pods | grep haproxy- | awk '{print $1}'); do kubectl --context et-ma2-prod -n ingress delete pod $pod &; done; wait
```