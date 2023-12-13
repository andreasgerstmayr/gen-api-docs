# crd-to-cr
This tool reads a Kubernetes Custom Resource Definition (CRD) and outputs a commented, full example Custom Resource (CR).

# Usage
To generate a full CR of the [Tempo Operator CRD](https://raw.githubusercontent.com/grafana/tempo-operator/5a79e619f268dac0fefd6cc394555582b17de520/bundle/community/manifests/tempo.grafana.com_tempostacks.yaml):
```
./crd-to-cr.py < bundle/community/manifests/tempo.grafana.com_tempostacks.yaml > full_cr.yaml
```
