package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	istio "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Parameter", "Value"})
	t.AppendRows([]table.Row{{"Deployment name", deploymentName}})
	t.AppendRows([]table.Row{{"Namespace", namespace}})
	if goalOfOperationIsToRemoveLabel {
		t.AppendRows([]table.Row{{"Label", labelToChangeKey}})
	} else {
		t.AppendRows([]table.Row{{"Label", fmt.Sprintf("%s=%s", labelToChangeKey, labelToChangeValue)}})
	}
	t.AppendRows([]table.Row{{"Will the label be removed?", goalOfOperationIsToRemoveLabel}})
	t.SetStyle(table.StyleColoredBright)
	t.Render()
	fmt.Println()
}

func resourcesAnalyze(
	clientset *kubernetes.Clientset,
	istioClient *istio.Clientset,
	namespace string,
	deploymentName string,
	changingLabelKey string,
) {
	deployment, _ := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, v1.GetOptions{})
	service, _ := clientset.CoreV1().Services(namespace).Get(context.TODO(), deploymentName, v1.GetOptions{})
	destinationRule, _ := istioClient.NetworkingV1alpha3().DestinationRules(namespace).Get(context.TODO(), deploymentName, v1.GetOptions{})

	deploymentSelectorLabels := deployment.Spec.Template.ObjectMeta.Labels
	serviceSelectorLabels := service.Spec.Selector
	destinationRuleSelectorLabels := destinationRule.Spec.Subsets[0].Labels

	fmt.Println()
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Type", "Name", "Detected", "labels count", "valid"})
	t.AppendRows([]table.Row{{
		"Deployment",
		If(deployment.Name != "", deployment.Name, "—"),
		If(deployment != nil, "✅", "❌"),
		len(deploymentSelectorLabels),
		If(len(deploymentSelectorLabels) == 1 && deploymentSelectorLabels[changingLabelKey] != "", "❌", "✅"),
	}})
	t.AppendRows([]table.Row{{
		"Service",
		If(service.Name != "", service.Name, "—"),
		If(service.Name != "", "✅", "❌"),
		len(serviceSelectorLabels),
		If(len(serviceSelectorLabels) == 1 && serviceSelectorLabels[changingLabelKey] != "", "❌", "✅"),
	}})
	t.AppendRows([]table.Row{{
		"<Istio> DestinationRule",
		If(service.Name != "", destinationRule.Name, "—"),
		If(service.Name != "", "✅", "❌"),
		len(destinationRuleSelectorLabels),
		If(len(destinationRuleSelectorLabels) == 1 && destinationRuleSelectorLabels[changingLabelKey] != "", "❌", "✅"),
	}})
	t.SetStyle(table.StyleColoredBright)
	t.Render()
	fmt.Println()

	if len(deploymentSelectorLabels) == 1 && deploymentSelectorLabels[changingLabelKey] != "" {
		logError(fmt.Sprintf("The label \"%s\" can not be edited because it's the only one in the matching set for the deployment", changingLabelKey))
		os.Exit(1)
	}

	if len(serviceSelectorLabels) == 1 && serviceSelectorLabels[changingLabelKey] != "" {
		logError(fmt.Sprintf("The label \"%s\" can not be edited because it's the only one in the matching set for the service", changingLabelKey))
		os.Exit(1)
	}
}
