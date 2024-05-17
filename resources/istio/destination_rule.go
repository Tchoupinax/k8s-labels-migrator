package authorizationPolicy

import (
	"context"

	resource "github.com/Tchoupinax/k8s-labels-migrator/resources"
	"github.com/Tchoupinax/k8s-labels-migrator/utils"
	istio "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func IstioDestinationRuleResourceAnalyze(
	istioClient *istio.Clientset,
	namespace string,
	matchingLabels map[string]string,
) []resource.Resource {
	destinationRules, _ := istioClient.NetworkingV1alpha3().DestinationRules(namespace).List(context.TODO(), v1.ListOptions{})

	var final []resource.Resource
	for _, item := range destinationRules.Items {
		if item.Spec.Subsets[0].Labels != nil {
			if utils.IsMatchSelectorsInclude(matchingLabels, item.Spec.Subsets[0].Labels) {
				final = append(final, resource.Resource{
					ApiVersion: "networking.istio.io/v1beta1",
					Category:   "Istio",
					Kind:       "DestinationRule",
					Labels:     item.ObjectMeta.Labels,
					Selectors:  item.Spec.Subsets[0].Labels,
					Name:       item.GetName(),
				})
			}
		}
	}

	return final
}
