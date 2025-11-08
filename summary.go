package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	resource "github.com/Tchoupinax/k8s-labels-migrator/resources"
	authorizationPolicy "github.com/Tchoupinax/k8s-labels-migrator/resources/istio"
	keda "github.com/Tchoupinax/k8s-labels-migrator/resources/keda"
	"github.com/Tchoupinax/k8s-labels-migrator/resources/native"
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
	startTime := time.Now()

	deployment, _ := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, v1.GetOptions{})
	if deployment.Name == "" {
		utils.LogError(fmt.Sprintf("The deployment \"%s\" was not found in the namespace \"%s\"", deploymentName, namespace))
		os.Exit(1)
	}

	var resources []resource.Resource
	resources = append(resources, resource.Resource{
		Kind:       "Deployment",
		ApiVersion: "apps/v1",
		Name:       deploymentName,
		Labels:     deployment.Labels,
		Selectors:  map[string]string{},
		Category:   "Native",
	})

	podLabels := deployment.Spec.Template.Labels

	// Native Services
	services := native.NativeServiceResourceAnalyze(
		clientset,
		namespace,
		podLabels,
	)

	resources = append(resources, services...)

	// Native HPA
	hpa := native.NativeHPAResourceAnalyze(
		clientset,
		namespace,
		deploymentName,
	)

	resources = append(resources, hpa...)

	// Istio Authorization Policies
	authorizationPolicies := authorizationPolicy.IstioAuthorizationPolicyResourceAnalyze(
		istioClient,
		namespace,
		podLabels,
	)

	resources = append(resources, authorizationPolicies...)

	// Destination rules
	destinationRules := authorizationPolicy.IstioDestinationRuleResourceAnalyze(
		istioClient,
		namespace,
		podLabels,
	)

	resources = append(resources, destinationRules...)

	// Keda — ScaledObject
	scaledobjects := keda.KedaScaledObjectResourceAnalyze(
		crdClient,
		namespace,
		deploymentName,
	)

	resources = append(resources, scaledobjects...)

	if len(podLabels) == 1 && podLabels[changingLabelKey] != "" {
		utils.LogError(fmt.Sprintf("The label \"%s\" can not be edited because it's the only one in the matching set for the deployment", changingLabelKey))
		os.Exit(1)
	}

	if len(podLabels) == 1 && podLabels[changingLabelKey] != "" {
		utils.LogError(fmt.Sprintf("The label \"%s\" can not be edited because it's the only one in the matching set for the service", changingLabelKey))
		os.Exit(1)
	}

	utils.LogSuccess("Resources analyzed (" + strconv.FormatFloat(time.Since(startTime).Seconds(), 'f', 2, 64) + "s)")

	webapp.StartWebServer(
		deploymentName,
		resources,
		podLabels,
	)
}
