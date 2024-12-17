package native

import (
	"context"

	resource "github.com/Tchoupinax/k8s-labels-migrator/resources"
	"github.com/Tchoupinax/k8s-labels-migrator/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func NativeServiceResourceAnalyze(
	clientset *kubernetes.Clientset,
	namespace string,
	matchingLabels map[string]string,
) []resource.Resource {
	destinationRules, _ := clientset.CoreV1().Services(namespace).List(context.TODO(), v1.ListOptions{})

	var final []resource.Resource
	for _, item := range destinationRules.Items {
		if item.Spec.Selector != nil {
			if utils.IsMatchSelectorsInclude(matchingLabels, item.Spec.Selector) {
				final = append(final, resource.Resource{
					ApiVersion: "v1",
					Category:   "Native",
					Kind:       "Service",
					Name:       item.GetName(),
					Labels:     item.ObjectMeta.Labels,
					Selectors:  item.Spec.Selector,
				})
			}
		}
	}

	return final
}
