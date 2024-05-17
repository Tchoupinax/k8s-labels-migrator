package native

import (
	"context"

	resource "github.com/Tchoupinax/k8s-labels-migrator/resources"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func NativeHPAResourceAnalyze(
	clientset *kubernetes.Clientset,
	namespace string,
	matchingDeploymentName string,
) []resource.Resource {
	hpa, _ := clientset.AutoscalingV2().HorizontalPodAutoscalers(namespace).List(context.TODO(), v1.ListOptions{})

	var final []resource.Resource
	for _, item := range hpa.Items {
		if item.Spec.ScaleTargetRef.Name != "" {
			if item.Spec.ScaleTargetRef.Name == matchingDeploymentName {
				final = append(final, resource.Resource{
					ApiVersion: "autoscaling/v2",
					Category:   "Native",
					Kind:       "HorizontalPodAutoscaler",
					Name:       item.GetName(),
					Labels:     item.ObjectMeta.Labels,
					Selectors: map[string]string{
						"kind": "Deployment",
						"name": item.Spec.ScaleTargetRef.Name,
					},
				})
			}
		}
	}

	return final
}
