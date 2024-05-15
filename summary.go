package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	resource "github.com/Tchoupinax/k8s-labels-migrator/resources"
	authorizationPolicy "github.com/Tchoupinax/k8s-labels-migrator/resources/istio"
	keda "github.com/Tchoupinax/k8s-labels-migrator/resources/keda"
	"github.com/Tchoupinax/k8s-labels-migrator/summary/webapp"
	utils "github.com/Tchoupinax/k8s-labels-migrator/utils"
	table "github.com/jedib0t/go-pretty/v6/table"
	istio "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

func displaySummary(
	namespace string,
	deploymentName string,
	labelToChangeKey string,
	labelToChangeValue string,
	goalOfOperationIsToRemoveLabel bool,
) {
	fmt.Println()
	t := table.NewWriter()
	t.AppendHeader(table.Row{"Parameter", "Value"})
	t.AppendRows([]table.Row{{"Deployment name", deploymentName}})
	t.AppendRows([]table.Row{{"Namespace", namespace}})
	if goalOfOperationIsToRemoveLabel {
		t.AppendRows([]table.Row{{"Label", labelToChangeKey}})
	} else {
		t.AppendRows([]table.Row{{"Label", fmt.Sprintf("%s=%s", labelToChangeKey, labelToChangeValue)}})
	}
	t.AppendRows([]table.Row{{"Will the label be removed?", goalOfOperationIsToRemoveLabel}})
	t.SetStyle(table.StyleColoredBlackOnYellowWhite)
	t.SetOutputMirror(os.Stdout)
	t.Render()
	fmt.Println()
}

func resourcesAnalyze(
	clientset *kubernetes.Clientset,
	istioClient *istio.Clientset,
	crdClient *dynamic.DynamicClient,
	namespace string,
	deploymentName string,
	changingLabelKey string,
) {
	deployment, _ := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, v1.GetOptions{})
	service, _ := clientset.CoreV1().Services(namespace).Get(context.TODO(), deploymentName, v1.GetOptions{})

	var resources []resource.Resource
	resources = append(resources, resource.Resource{
		Kind:       "Deployment",
		ApiVersion: "apps/v1",
		Name:       deploymentName,
		Selectors:  deployment.ObjectMeta.Labels,
		Category:   "Native",
	})
	resources = append(resources, resource.Resource{
		Kind:       "Service",
		ApiVersion: "v1",
		Name:       service.Name,
		Selectors:  service.Spec.Selector,
		Category:   "Native",
	})

	matchingLabels := deployment.Spec.Template.ObjectMeta.Labels

	// Istio Authorization Policies
	authorizationPolicies := authorizationPolicy.IstioAuthorizationPolicyResourceAnalyze(
		istioClient,
		namespace,
		matchingLabels,
	)
	for _, a := range authorizationPolicies {
		resources = append(resources, a)
	}
	// Destination rules
	destinationRules := authorizationPolicy.IstiDestinationRuleResourceAnalyze(
		istioClient,
		namespace,
		matchingLabels,
	)
	for _, a := range destinationRules {
		resources = append(resources, a)
	}
	// Keda — ScaledObject
	scaledobjects := keda.KedaScaledObjectResourceAnalyze(
		crdClient,
		namespace,
		deploymentName,
	)
	for _, a := range scaledobjects {
		resources = append(resources, a)
	}

	fmt.Println()
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Type", "Name", "Detected", "labels count", "labels", "valid"})
	for _, resource := range resources {
		t.AppendRows([]table.Row{{
			resource.Kind,
			resource.Name,
			"✅",
			len(resource.Selectors),
			strings.Join(utils.MapToArray(resource.Selectors), "\n"),
			utils.If(len(resource.Selectors) == 1 && resource.Selectors[changingLabelKey] != "", "❌", "✅"),
		}})
	}

	t.SetStyle(table.StyleColoredBlackOnYellowWhite)
	t.Render()
	fmt.Println()

	if len(matchingLabels) == 1 && matchingLabels[changingLabelKey] != "" {
		utils.LogError(fmt.Sprintf("The label \"%s\" can not be edited because it's the only one in the matching set for the deployment", changingLabelKey))
		os.Exit(1)
	}

	if len(matchingLabels) == 1 && matchingLabels[changingLabelKey] != "" {
		utils.LogError(fmt.Sprintf("The label \"%s\" can not be edited because it's the only one in the matching set for the service", changingLabelKey))
		os.Exit(1)
	}

	webapp.StartWebServer(resources)
}
