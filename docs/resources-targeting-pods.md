# Service is not the only resource that target pods

We did the base migration considering pods are matched by service. But for a complexe app, there is a lot of resources that can match pods.

## Identify resources

- Kubernetes
  - PodDisruptionBudget
      
- Keda
  - ❌ ScaledObject: — `.spec.scaleTargetRef.name` (equal to the deployment name)
    - ➡️ It is acceptable to not manage this as the migration will be fast.
- Istio
  - AuthorizationPolicy
  - RequestAuthentication
  - DestinationRule
  - Virtual service (host match the DNS name, means it match the deployment name [docs](https://istio.io/latest/docs/reference/config/networking/virtual-service/#VirtualService))
- Monitoring
  - PrometheusRule (rules could match pod or deployment in the query)
  - PodMonitor — `.spec.selector.matchLabels`

## Tips

```bash
# Display existing resources in the cluster
kubectl api-resources --verbs=list --namespaced -o name
```

export KUBERNETES_RESOURCE=ScaledObject
export NAME=
export NAMESPACE=

kubectl get $KUBERNETES_RESOURCE $NAME -n $NAMESPACE -o yaml | yq '.spec.selector.matchLabels'