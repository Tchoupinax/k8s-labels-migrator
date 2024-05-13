package main

import (
	"flag"
	"fmt"
	"os"

	utils "github.com/Tchoupinax/k8s-labels-migrator/utils"
	"github.com/mbndr/figlet4go"
	istio "istio.io/client-go/pkg/clientset/versioned"
	"k8s.io/client-go/dynamic"
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
	if err != nil {
		panic(err.Error())
	}
	istioClient, err := istio.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	crdClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	var deploymentName = ""
	var namespace = ""
	var labelToChangeKey = ""
	var labelToChangeValue = ""
	var goalOfOperationIsToRemoveLabel = false
	var matcherLabels = []string{
		"app.kubernetes.io/instance",
		"app.kubernetes.io/name",
	}

	flag.StringVar(&deploymentName, "deployment", "", "Name of the deployment to edit label")
	flag.StringVar(&namespace, "namespace", "", "Namespace of the deployment to edit label")
	flag.BoolVar(&goalOfOperationIsToRemoveLabel, "remove-label", false, "If true, the label will be removed instead of be added/edited")
	flag.StringVar(&labelToChangeKey, "label", "app.kubernetes.io/name", "Name of the label")
	flag.StringVar(&labelToChangeValue, "value", "", "Value of the label")
	flag.Parse()

	if deploymentName == "" {
		utils.LogError("Deployment name is mandatory")
		os.Exit(1)
	}
	if namespace == "" {
		utils.LogError("Namespace is mandatory")
		os.Exit(1)
	}
	if labelToChangeValue == "" && !goalOfOperationIsToRemoveLabel {
		utils.LogError("label value is mandatory")
		os.Exit(1)
	}

	fmt.Print("\033[H\033[2J")
	ascii := figlet4go.NewAsciiRender()
	renderStr, _ := ascii.Render("k8s labels migrator")
	fmt.Print(renderStr)

	utils.LogInfo("Analyzing your cluster...")
	resourcesAnalyze(clientset, istioClient, crdClient, namespace, deploymentName, labelToChangeKey)
	utils.LogSuccess("Cluster ready")
	displaySummary(
		namespace,
		deploymentName,
		labelToChangeKey,
		labelToChangeValue,
		goalOfOperationIsToRemoveLabel,
	)

	c := utils.AskForConfirmation("Do you validate these parameters?")
	if !c {
		utils.LogInfo("Operation aborted by the user")
		os.Exit(0)
	}
	c2 := utils.AskForConfirmation("I confirm that I have no gitops tool overriding my config (e.g. ArgoCD auto-sync)")
	if !c2 {
		utils.LogInfo("Operation aborted by the user")
		os.Exit(0)
	}

	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	utils.LogWarning("PLEASE DO NOT INTERRUPT THE PROCESS UNTIL THE END!")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")

	MigrationWorkflow(
		namespace,
		clientset,
		istioClient,
		crdClient,
		deploymentName,
		labelToChangeKey,
		labelToChangeValue,
		goalOfOperationIsToRemoveLabel,
	)

	if utils.ArrayContains(matcherLabels, labelToChangeKey) {
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

	fmt.Println("")
	utils.LogSuccess("Migration terminated with success!")
	fmt.Println("")
}
