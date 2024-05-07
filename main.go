package main

import (
	"flag"
	"os"

	istio "istio.io/client-go/pkg/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Check if the helper is asked by flag
	cliCommandDisplayHelp(os.Args)

	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	istioClient, err := istio.NewForConfig(config)
	//crdClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	var deploymentName = ""
	var namespace = ""
	var labelToChangeKey = ""
	var labelToChangeValue = ""
	var goalOfOperationIsToRemoveLabel = false

	flag.StringVar(&deploymentName, "deployment", "", "Name of the deployment to edit label")
	flag.StringVar(&namespace, "namespace", "", "Namespace of the deployment to edit label")
	flag.BoolVar(&goalOfOperationIsToRemoveLabel, "remove-label", false, "If true, the label will be removed instead of be added/edited")
	flag.StringVar(&labelToChangeKey, "label", "app.kubernetes.io/name", "Name of the label")
	flag.StringVar(&labelToChangeValue, "value", "", "Value of the label")
	flag.Parse()

	if deploymentName == "" {
		logError("Deployment name is mandatory")
		os.Exit(1)
	}
	if namespace == "" {
		logError("Namespace is mandatory")
		os.Exit(1)
	}
	if labelToChangeValue == "" && !goalOfOperationIsToRemoveLabel {
		logError("label value is mandatory")
		os.Exit(1)
	}

	logInfo("Analyzing your cluster...")
	resourcesAnalyze(clientset, istioClient, namespace, deploymentName, labelToChangeKey)
	logSuccess("Cluster ready")
	displaySummary(
		namespace,
		deploymentName,
		labelToChangeKey,
		labelToChangeValue,
		goalOfOperationIsToRemoveLabel,
	)

	c := askForConfirmation("Do you validate these parameters?")
	if !c {
		logInfo("Operation aborted by the user")
		os.Exit(0)
	}
	c2 := askForConfirmation("I confirm that I have no gitops tool overriding my config (e.g. ArgoCD auto-sync)")
	if !c2 {
		logInfo("Operation aborted by the user")
		os.Exit(0)
	}

	MigrationWorkflow(
		namespace,
		clientset,
		istioClient,
		deploymentName,
		labelToChangeKey,
		labelToChangeValue,
		goalOfOperationIsToRemoveLabel,
	)

	if labelToChangeKey == "app.kubernetes.io/name" {
		AddLabelToServiceSelector(
			namespace,
			clientset,
			deploymentName,
			labelToChangeKey,
			labelToChangeValue,
			goalOfOperationIsToRemoveLabel,
		)

		AddLabelToIstioDestinatonRulesSelector(
			namespace,
			clientset,
			istioClient,
			deploymentName,
			labelToChangeKey,
			labelToChangeValue,
			goalOfOperationIsToRemoveLabel,
		)
	}
}
