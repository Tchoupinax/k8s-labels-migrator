package main

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func waitUntilAllPodAreReady(
	clientset *kubernetes.Clientset,
	namespace string,
	deploymentName string,
) bool {
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logError(err.Error())
		return false
	}
	currentDeploymentHasAllPodReady := deployment.Status.Replicas-deployment.Status.ReadyReplicas == 0 && deployment.Status.Replicas > 1
	return currentDeploymentHasAllPodReady
}
