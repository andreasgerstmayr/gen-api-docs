# gen-api-docs
This tool reads a Kubernetes Custom Resource Definition (CRD) and outputs a commented, full example Custom Resource (CR).

# Usage
To generate a full CR of the [Tempo Operator CRD](https://raw.githubusercontent.com/grafana/tempo-operator/5a79e619f268dac0fefd6cc394555582b17de520/bundle/community/manifests/tempo.grafana.com_tempostacks.yaml):
```bash
go run main.go < tempo.grafana.com_tempostacks.yaml > tempostack.yaml
```

With docker:
```bash
docker run -i ghcr.io/andreasgerstmayr/gen-api-docs < tempo.grafana.com_tempostacks.yaml > tempostack.yaml
```

# Example Output
```yaml
apiVersion: tempo.grafana.com/v1alpha1   # APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
kind: TempoMonolithic                    # Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
metadata:
  name: example
spec:                                    # TempoMonolithicSpec defines the desired state of TempoMonolithic.
  extraConfig:                           # ExtraConfig defines any extra (overlay) configuration of components.
    tempo: {}                            # Tempo defines any extra Tempo configuration, which will be merged with the operator's generated Tempo configuration
  ingestion:                             # Ingestion defines the trace ingestion configuration.
    otlp:                                # OTLP defines the ingestion configuration for the OTLP protocol.
      grpc:                              # GRPC defines the OTLP over gRPC configuration.
        enabled: true                    # Enabled defines if OTLP over gRPC is enabled. Default: enabled.
        tls:                             # TLS defines the TLS configuration for OTLP/gRPC ingestion.
          enabled: false                 # Enabled defines if TLS is enabled.
          caName: ""                     # CA is the name of a ConfigMap containing a CA certificate (service-ca.crt). It needs to be in the same namespace as the Tempo custom resource.
          certName: ""                   # Cert is the name of a Secret containing a certificate (tls.crt) and private key (tls.key). It needs to be in the same namespace as the Tempo custom resource.
          minVersion: ""                 # MinVersion defines the minimum acceptable TLS version.
      http:                              # HTTP defines the OTLP over HTTP configuration.
        enabled: true                    # Enabled defines if OTLP over HTTP is enabled. Default: enabled.
        tls:                             # TLS defines the TLS configuration for OTLP/HTTP ingestion.
          enabled: false                 # Enabled defines if TLS is enabled.
          caName: ""                     # CA is the name of a ConfigMap containing a CA certificate (service-ca.crt). It needs to be in the same namespace as the Tempo custom resource.
          certName: ""                   # Cert is the name of a Secret containing a certificate (tls.crt) and private key (tls.key). It needs to be in the same namespace as the Tempo custom resource.
          minVersion: ""                 # MinVersion defines the minimum acceptable TLS version.
  jaegerui:                              # JaegerUI defines the Jaeger UI configuration.
    enabled: false                       # Enabled defines if the Jaeger UI component should be created.
    ingress:                             # Ingress defines the Ingress configuration for the Jaeger UI.
      enabled: false                     # Enabled defines if an Ingress object should be created for Jaeger UI.
      annotations: {}                    # Annotations defines the annotations of the Ingress object.
      host: ""                           # Host defines the hostname of the Ingress object.
      ingressClassName: ""               # IngressClassName defines the name of an IngressClass cluster resource. Defines which ingress controller serves this ingress resource.
    route:                               # Route defines the OpenShift route configuration for the Jaeger UI.
      enabled: false                     # Enabled defines if a Route object should be created for Jaeger UI.
      annotations: {}                    # Annotations defines the annotations of the Route object.
      host: ""                           # Host defines the hostname of the Route object.
      termination: "edge"                # Termination specifies the termination type. Default: edge.
    resources:                           # Resources defines the compute resource requirements of the Jaeger UI container.
      claims:                            # Claims lists the names of resources, defined in spec.resourceClaims, that are used by this container.   This is an alpha field and requires enabling the DynamicResourceAllocation feature gate.   This field is immutable. It can only be set for containers.
      - name: ""                         # Name must match the name of one entry in pod.spec.resourceClaims of the Pod where this field is used. It makes that resource available inside a container.
      limits:                            # Limits describes the maximum amount of compute resources allowed. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
        cpu: "750m"
        memory: "2Gi"
      requests:                          # Requests describes the minimum amount of compute resources required. If Requests is omitted for a container, it defaults to Limits if that is explicitly specified, otherwise to an implementation-defined value. Requests cannot exceed Limits. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
        cpu: "500m"
        memory: "1Gi"
  management: ""                         # ManagementState defines whether this instance is managed by the operator or self-managed. Default: Managed.
  observability:                         # Observability defines the observability configuration of the Tempo deployment.
    grafana:                             # Grafana defines the Grafana configuration of the Tempo deployment.
      dataSource:                        # DataSource defines the Grafana data source configuration.
        enabled: false                   # Enabled defines if a Grafana data source should be created for this Tempo deployment.
        instanceSelector:                # InstanceSelector defines the Grafana instance where the data source should be created.
          matchExpressions:              # matchExpressions is a list of label selector requirements. The requirements are ANDed.
          - key: ""                      # key is the label key that the selector applies to.
            operator: ""                 # operator represents a key's relationship to a set of values. Valid operators are In, NotIn, Exists and DoesNotExist.
            values:                      # values is an array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. This array is replaced during a strategic merge patch.
            - ""
          matchLabels: {}                # matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed.
    metrics:                             # Metrics defines the metric configuration of the Tempo deployment.
      prometheusRules:                   # ServiceMonitors defines the PrometheusRule configuration.
        enabled: false                   # Enabled defines if PrometheusRule objects should be created for this Tempo deployment.
      serviceMonitors:                   # ServiceMonitors defines the ServiceMonitor configuration.
        enabled: false                   # Enabled defines if ServiceMonitor objects should be created for this Tempo deployment.
  storage:                               # Storage defines the storage configuration.
    traces:                              # Traces defines the storage configuration for traces.
      azure:                             # Azure defines the configuration for Azure Storage.
        secret: ""                       # Secret is the name of a Secret containing credentials for accessing object storage. It needs to be in the same namespace as the TempoMonolithic custom resource.
      backend: "memory"                  # Backend defines the backend for storing traces. Default: memory.
      gcs:                               # GCP defines the configuration for Google Cloud Storage.
        secret: ""                       # Secret is the name of a Secret containing credentials for accessing object storage. It needs to be in the same namespace as the TempoMonolithic custom resource.
      s3:                                # S3 defines the configuration for Amazon S3.
        secret: ""                       # Secret is the name of a Secret containing credentials for accessing object storage. It needs to be in the same namespace as the TempoMonolithic custom resource.
        tls:                             # TLS defines the TLS configuration for Amazon S3.
          enabled: false                 # Enabled defines if TLS is enabled.
          caName: ""                     # CA is the name of a ConfigMap containing a CA certificate (service-ca.crt). It needs to be in the same namespace as the Tempo custom resource.
          certName: ""                   # Cert is the name of a Secret containing a certificate (tls.crt) and private key (tls.key). It needs to be in the same namespace as the Tempo custom resource.
          minVersion: ""                 # MinVersion defines the minimum acceptable TLS version.
      size: "10Gi"                       # Size defines the size of the volume where traces are stored. For in-memory storage, this defines the size of the tmpfs volume. For persistent volume storage, this defines the size of the persistent volume. For object storage, this defines the size of the persistent volume containing the Write-Ahead Log (WAL) of Tempo. Default: 10Gi.
  affinity:                              # Affinity defines the Affinity rules for scheduling pods.
    nodeAffinity: {}                     # Describes node affinity scheduling rules for the pod.
    podAffinity: {}                      # Describes pod affinity scheduling rules (e.g. co-locate this pod in the same node, zone, etc. as some other pod(s)).
    podAntiAffinity: {}                  # Describes pod anti-affinity scheduling rules (e.g. avoid putting this pod in the same node, zone, etc. as some other pod(s)).
  nodeSelector: {}                       # NodeSelector defines which labels are required by a node to schedule the pod onto it.
  resources:                             # Resources defines the compute resource requirements of the Tempo container.
    claims:                              # Claims lists the names of resources, defined in spec.resourceClaims, that are used by this container.   This is an alpha field and requires enabling the DynamicResourceAllocation feature gate.   This field is immutable. It can only be set for containers.
    - name: ""                           # Name must match the name of one entry in pod.spec.resourceClaims of the Pod where this field is used. It makes that resource available inside a container.
    limits:                              # Limits describes the maximum amount of compute resources allowed. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
      cpu: "750m"
      memory: "2Gi"
    requests:                            # Requests describes the minimum amount of compute resources required. If Requests is omitted for a container, it defaults to Limits if that is explicitly specified, otherwise to an implementation-defined value. Requests cannot exceed Limits. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
      cpu: "500m"
      memory: "1Gi"
  tolerations: {}                        # Tolerations defines the tolerations of a node to schedule the pod onto it.
status:                                  # TempoMonolithicStatus defines the observed state of TempoMonolithic.
  components:                            # Components provides summary of all Tempo pod status, grouped per component.
    tempo:                               # Tempo is a map of the pod status of the Tempo pods.
      "key":
      - ""
  conditions:                            # Conditions of the Tempo deployment health.
  - lastTransitionTime: ""               # lastTransitionTime is the last time the condition transitioned from one status to another. This should be when the underlying condition changed.  If that is not known, then using the time when the API field changed is acceptable.
    message: ""                          # message is a human readable message indicating details about the transition. This may be an empty string.
    observedGeneration: 0                # observedGeneration represents the .metadata.generation that the condition was set based upon. For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date with respect to the current state of the instance.
    reason: ""                           # reason contains a programmatic identifier indicating the reason for the condition's last transition. Producers of specific condition types may define expected values and meanings for this field, and whether the values are considered a guaranteed API. The value should be a CamelCase string. This field may not be empty.
    status: ""                           # status of the condition, one of True, False, Unknown.
    type: ""                             # type of condition in CamelCase or in foo.example.com/CamelCase. --- Many .condition.type values are consistent across resources like Available, but because arbitrary conditions can be useful (see .node.status.conditions), the ability to deconflict is important. The regex it matches is (dns1123SubdomainFmt/)?(qualifiedNameFmt)
```
