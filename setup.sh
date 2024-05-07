k3d cluster create k8s-labels-migrator \
  --kubeconfig-update-default

helm repo add istio https://istio-release.storage.googleapis.com/charts
helm repo update

kubectl create namespace istio-system
helm install istio-base istio/base -n istio-system --set defaultRevision=default
helm install istio-istiod istio/istiod -n istio-system --set defaultRevision=default
