set -ex

kubectl -n monitoring port-forward svc/kube-prometheus-stack-grafana "8888:80"