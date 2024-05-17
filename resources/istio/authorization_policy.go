package authorizationPolicy

import (
	"context"

	resource "github.com/Tchoupinax/k8s-labels-migrator/resources"
	"github.com/Tchoupinax/k8s-labels-migrator/utils"
	istio "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func IstioAuthorizationPolicyResourceAnalyze(
	istioClient *istio.Clientset,
	namespace string,
	matchingLabels map[string]string,
) []resource.Resource {
	authorizationPolicies, _ := istioClient.SecurityV1beta1().AuthorizationPolicies(namespace).List(context.TODO(), v1.ListOptions{})

	var final []resource.Resource
	for _, item := range authorizationPolicies.Items {
		if item.Spec.Selector != nil {
			if utils.IsMatchSelectorsInclude(matchingLabels, item.Spec.Selector.MatchLabels) {
				final = append(final, resource.Resource{
					ApiVersion: "security.istio.io/v1beta1",
					Category:   "Istio",
					Kind:       "AuthorizationPolicy",
					Labels:     item.ObjectMeta.Labels,
					Selectors:  item.Spec.Selector.MatchLabels,
					Name:       item.GetName(),
				})
			}
		}
	}

	return final
}
