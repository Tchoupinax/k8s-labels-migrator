# Managed resourced

We identify a set of resources that are matching or are linked to the pods we are trying to migrate.

- #### Keda
  - ✅ `ScaledObject`: This resources target the deployment (`.spec.scaleTargetRef.name`) and not pods. However, we don't want Keda to perfom operations during the procedure. That's why we pause the keda scaled object during the operation. This is done with an annotation meant for that purpose by keda.sh. (see `keda.go`).

- #### Kubernetes
  - 🛑 `PodDisruptionBudget`

- #### Istio
  - 🛑 `AuthorizationPolicy`
  - ✅ `DestinationRule`: like Kubernetes services, destination rules match pods with matching labels. We repeat the same operation we did for service for the replacement.
  - 🛑  `VirtualService` (host match the DNS name, means it match the deployment name [docs](https://istio.io/latest/docs/reference/config/networking/virtual-service/#VirtualService))

- #### Monitoring
  - 🛑 `PrometheusRule`: (rules could match pod or deployment in the query)
  - 🛑 `PodMonitor`: matches pods by labels (`.spec.selector.matchLabels`)


```bash
# Display existing resources in the cluster
kubectl api-resources --verbs=list --namespaced -o name
```

export KUBERNETES_RESOURCE=ScaledObject
export NAME=
export NAMESPACE=

kubectl get $KUBERNETES_RESOURCE $NAME -n $NAMESPACE -o yaml | yq '.spec.selector.matchLabels'