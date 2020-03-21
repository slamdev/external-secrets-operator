# External Secrets Operator

External Secrets Operator that watches changes in external key-value backends to sync them with Kubernetes Secrets and ConfigMaps.

## Supported Backends

- [Vault](https://www.vaultproject.io/)
- [Consul](https://www.consul.io/)

## Installing to a cluster

- via helm chart: https://github.com/slamdev/helm-charts/tree/master/charts/external-secrets-operator
- via kube manifests:
```sh
$ curl https://raw.githubusercontent.com/slamdev/external-secrets-operator/master/k8s/k8s.yaml | kubectl apply -f -
```

Check [samples](config/samples) for usage examples.

## Developing

```sh
# install all CRDs to a cluster
$ make install
# run operator in a cluster
$ make run
# deploy samples to a cluster
$ make samples
```
