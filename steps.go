package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	keda "github.com/Tchoupinax/k8s-labels-migrator/resources"
	utils "github.com/Tchoupinax/k8s-labels-migrator/utils"
	istio "istio.io/client-go/pkg/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

func MigrationWorkflow(
	namespace string,
	clientset *kubernetes.Clientset,
	istioClient *istio.Clientset,
	crdClient *dynamic.DynamicClient,
	deploymentName string,
	changingLabelKey string,
	changingLabelValue string,
	removeLabel bool,
) {
	currentService, _ := clientset.CoreV1().Services(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	currentDeployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	currentDestinationRule, _ := istioClient.NetworkingV1alpha3().DestinationRules(namespace).Get(context.TODO(), deploymentName, v1.GetOptions{})
	if err != nil {
		utils.LogError("No deployment found.")
		os.Exit(1)
	}

	utils.LogInfo("1. Creating the clone deployment")
	var temporalDeployment = *currentDeployment
	temporalDeployment.GenerateName = fmt.Sprintf("%s-%s", currentDeployment.Name, "changing-label-tmp")
	temporalDeployment.Name = fmt.Sprintf("%s-%s", currentDeployment.Name, "changing-label-tmp")
	// It's required because an error is thrown, we can not create a deployment with this property provided
	temporalDeployment.ResourceVersion = ""

	_, err = clientset.AppsV1().Deployments(namespace).Create(context.TODO(), &temporalDeployment, metav1.CreateOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "already exists, the server was not able to generate a unique name for the object") {
			utils.LogWarning("1. Temporary deployment already created. Continue...")
		}
	} else {
		utils.LogSuccess("1. Deployment replicated")
	}

	utils.LogInfo("2. Updating the service")
	var temporalService = *currentService
	delete(temporalService.Spec.Selector, changingLabelKey)
	_, err = clientset.CoreV1().Services(namespace).Update(context.TODO(), &temporalService, metav1.UpdateOptions{})
	utils.Check(err)
	utils.LogSuccess("2. Service updated")

	utils.LogInfo("2.1 Updating Istio destination rules")
	var temporalDestinationRule = currentDestinationRule
	delete(temporalDestinationRule.Spec.Subsets[0].Labels, changingLabelKey)
	_, err = istioClient.NetworkingV1alpha3().DestinationRules(namespace).Update(context.TODO(), temporalDestinationRule, metav1.UpdateOptions{})
	utils.Check(err)
	utils.LogSuccess("2.1 Istio destination rules updated")

	utils.LogInfo("2.2 Keda")
	keda.PauseScaledObject(crdClient, clientset, deploymentName, namespace)

	utils.LogBlocking("3. Waiting while temporals pods are not totally ready to handle traffic")
	areAllPodReady := false
	for !areAllPodReady {
		utils.LogBlockingDot()
		time.Sleep(1 * time.Second)
		areAllPodReady =
			waitUntilAllPodAreReady(clientset, namespace, currentDeployment.Name) &&
				waitUntilAllPodAreReady(clientset, namespace, fmt.Sprintf("%s-%s", currentDeployment.Name, "changing-label-tmp"))
	}
	fmt.Println("")

	utils.LogInfo("4. Delete the old deployment")
	// Delete the old deployment
	deleteError := clientset.AppsV1().Deployments(namespace).Delete(context.TODO(), currentDeployment.Name, *metav1.NewDeleteOptions(0))
	utils.Check(deleteError)
	utils.LogSuccess("4. Old deployment deleted")

	utils.LogInfo("5. Creating the original deployment with modified label")
	var futureOfficialDeployment = *currentDeployment
	futureOfficialDeployment.GenerateName = deploymentName
	futureOfficialDeployment.Name = deploymentName
	// It's required because an error is thrown, we can not create a deployment with this property provided
	futureOfficialDeployment.ResourceVersion = ""

	// Update the label of the deploy and pods
	// If the key is empty, remove the label because we consider empty string is the order to remove the label
	if removeLabel {
		delete(futureOfficialDeployment.ObjectMeta.Labels, changingLabelKey)
		delete(futureOfficialDeployment.Spec.Template.ObjectMeta.Labels, changingLabelKey)
		delete(futureOfficialDeployment.Spec.Selector.MatchLabels, changingLabelKey)
	} else {
		// Label of the deployment
		futureOfficialDeployment.ObjectMeta.Labels[changingLabelKey] = changingLabelValue
		// Label of the pod created by the deployment
		futureOfficialDeployment.Spec.Template.ObjectMeta.Labels[changingLabelKey] = changingLabelValue
		// Then we must include the label in the matchSelector for the deployment to find pods
		futureOfficialDeployment.Spec.Selector.MatchLabels[changingLabelKey] = changingLabelValue
	}

	_, err = clientset.AppsV1().Deployments(namespace).Create(context.TODO(), &futureOfficialDeployment, metav1.CreateOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "already exists, the server was not able to generate a unique name for the object") {
			fmt.Println("⚠️  Temporary deployment already created. Continue...")
		}
		fmt.Println(err)
	}
	utils.LogSuccess("5. Deployment created")

	utils.LogBlocking("6. Waiting while final pods are not totally ready to handle traffic")
	areAllPodReady = false
	for !areAllPodReady {
		utils.LogBlockingDot()
		time.Sleep(1 * time.Second)
		areAllPodReady = waitUntilAllPodAreReady(clientset, namespace, currentDeployment.Name)
	}
	fmt.Println("")

	utils.LogInfo("7. Deleting temporal deployment")
	time.Sleep(1 * time.Second)
	// Delete the temporal deployment
	errDeleteTmpDeploy := clientset.AppsV1().Deployments(namespace).Delete(context.TODO(), fmt.Sprintf("%s-%s", currentDeployment.Name, "changing-label-tmp"), metav1.DeleteOptions{})
	utils.Check(errDeleteTmpDeploy)
	utils.LogSuccess("7. Temporary deployment deleted")

	utils.LogInfo("8 Keda")
	keda.ResumeScaledObject(crdClient, clientset, deploymentName, namespace)
}

func AddLabelToServiceSelector(
	namespace string,
	clientset *kubernetes.Clientset,
	applicationName string,
	changingLabelKey string,
	changingLabelValue string,
	removeLabel bool,
) {
	utils.LogInfo("====== Additionnal step ====================================")
	utils.LogInfo("8. Add the label as a selector in the service...")
	// Get the current service
	currentService, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), applicationName, metav1.GetOptions{})
	utils.Check(err)

	var futureService = *currentService
	if removeLabel {
		// Update the value of the label
		delete(futureService.Spec.Selector, changingLabelKey)
	} else {
		futureService.Spec.Selector[changingLabelKey] = changingLabelValue
		// If the string is empty, remove the label (see deployment)
	}

	// Update the service in the cluster
	_, updateServiceError := clientset.CoreV1().Services(namespace).Update(context.TODO(), &futureService, metav1.UpdateOptions{})
	utils.Check(updateServiceError)

	utils.LogSuccess("8. Service configured")
	utils.LogInfo("============================================================")
}

func AddLabelToIstioDestinatonRulesSelector(
	namespace string,
	clientset *kubernetes.Clientset,
	istioClient *istio.Clientset,
	applicationName string,
	changingLabelKey string,
	changingLabelValue string,
	removeLabel bool,
) {
	utils.LogInfo("====== Additionnal step ====================================")
	utils.LogInfo("9. Add the label as a selector in istio destination rules...")
	currentDestinationRule, err := istioClient.NetworkingV1alpha3().DestinationRules(namespace).Get(context.TODO(), applicationName, v1.GetOptions{})
	utils.Check(err)

	var futureDestinationRule = currentDestinationRule
	if removeLabel {
		// Update the value of the label
		delete(futureDestinationRule.Spec.Subsets[0].Labels, changingLabelKey)
	} else {
		futureDestinationRule.Spec.Subsets[0].Labels[changingLabelKey] = changingLabelValue
		// If the string is empty, remove the label (see deployment)
	}

	// Update the service in the cluster
	_, updateDestinationRuleError := istioClient.NetworkingV1alpha3().DestinationRules(namespace).Update(context.TODO(), futureDestinationRule, metav1.UpdateOptions{})
	utils.Check(updateDestinationRuleError)

	utils.LogSuccess("9. Istio destination rules configured")
	utils.LogInfo("============================================================")
}
