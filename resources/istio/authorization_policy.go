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
		if item.Spec.GetSelector() != nil {
			if utils.IsMatchSelectorsInclude(matchingLabels, item.Spec.GetSelector().GetMatchLabels()) {
				final = append(final, resource.Resource{
					ApiVersion: "security.istio.io/v1beta1",
					Category:   "Istio",
					Kind:       "AuthorizationPolicy",
					Labels:     item.Labels,
					Selectors:  item.Spec.GetSelector().GetMatchLabels(),
					Name:       item.GetName(),
				})
			}
		}
	}

	return final
}
