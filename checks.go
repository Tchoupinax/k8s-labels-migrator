package main

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func isTheEditedLabelTheOnlyOne(
	namespace string,
	clientset *kubernetes.Clientset,
	deploymentName string,
	changingLabelKey string,
) bool {
	deployment, _ := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, v1.GetOptions{})
	labels := deployment.Spec.Template.ObjectMeta.Labels
	return len(labels) == 1 && labels[changingLabelKey] != ""
}
