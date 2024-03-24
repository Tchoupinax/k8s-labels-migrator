package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func Step1(
	namespace string,
	clientset *kubernetes.Clientset,
	deploymentName string,
	changingLabelKey string,
	changingLabelValue string,
	removeLabel bool,
) {
	currentService, _ := clientset.CoreV1().Services(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	currentDeployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logError("No deployment found.")
		os.Exit(1)
	}

	logInfo("1. Creating the clone deployment")
	var temporalDeployment = *currentDeployment
	temporalDeployment.GenerateName = fmt.Sprintf("%s-%s", currentDeployment.Name, "changing-label-tmp")
	temporalDeployment.Name = fmt.Sprintf("%s-%s", currentDeployment.Name, "changing-label-tmp")
	// It's required because an error is thrown, we can not create a deployment with this property provided
	temporalDeployment.ResourceVersion = ""

	_, err = clientset.AppsV1().Deployments(namespace).Create(context.TODO(), &temporalDeployment, metav1.CreateOptions{});
	if err != nil {
		if strings.Contains(err.Error(), "already exists, the server was not able to generate a unique name for the object") {
			logWarning("1. Temporary deployment already created. Continue...")
		}
	} else {
		logSuccess("1. Deployment replicated")
	}

	logInfo("2. Updating the service...")
	var temporalService = *currentService
	delete(temporalService.Spec.Selector, changingLabelKey)
	_, err = clientset.CoreV1().Services(namespace).Update(context.TODO(), &temporalService, metav1.UpdateOptions{});
	check(err)
	logSuccess("2. Service updated")

	logBlocking("3. Waiting while pods are not totally ready to handle traffic")
	areAllPodReady := false
	for !areAllPodReady {
		logBlockingDot()
		time.Sleep(1 * time.Second)
		areAllPodReady = 
			waitUntilAllPodAreReady(clientset, namespace, "api") && 
			waitUntilAllPodAreReady(clientset, namespace, fmt.Sprintf("%s-%s", currentDeployment.Name, "changing-label-tmp"))
	}
	fmt.Println("")

	logInfo("4. Delete the old deployment...")
	// Delete the old deployment
	deleteError := clientset.AppsV1().Deployments(namespace).Delete(context.TODO(), currentDeployment.Name, *metav1.NewDeleteOptions(0))
	check(deleteError)
	logSuccess("4. Old deployment deleted")

	logInfo("5. Creating the original deployment with modified label")
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
		} else {
		futureOfficialDeployment.ObjectMeta.Labels[changingLabelKey] = changingLabelValue
		futureOfficialDeployment.Spec.Template.ObjectMeta.Labels[changingLabelKey] = changingLabelValue
	}

	_, err = clientset.AppsV1().Deployments(namespace).Create(context.TODO(), &futureOfficialDeployment, metav1.CreateOptions{});
	if err != nil {
		if strings.Contains(err.Error(), "already exists, the server was not able to generate a unique name for the object") {
			fmt.Println("⚠️  Temporary deployment already created. Continue...")
		}
		fmt.Println(err)
	}
	logSuccess("5. Deployment created")

	logBlocking("6. Waiting while pods are not totally ready to handle traffic")
	areAllPodReady = false
	for !areAllPodReady {
		logBlockingDot()
		time.Sleep(1 * time.Second)
		areAllPodReady = waitUntilAllPodAreReady(clientset, namespace, "api")
	}
	fmt.Println("")

	logInfo("7. Deleting temporal deployment...")
	time.Sleep(1 * time.Second)
	// Delete the temporal deployment
	errDeleteTmpDeploy := clientset.AppsV1().Deployments(namespace).Delete(context.TODO(), fmt.Sprintf("%s-%s", currentDeployment.Name, "changing-label-tmp"), metav1.DeleteOptions{})
	check(errDeleteTmpDeploy)
	logSuccess("7. Temporary deployment deleted")
}

func AddLabelToServiceSelector (
	namespace string,
	clientset *kubernetes.Clientset,
	applicationName string,
	changingLabelKey string,
	changingLabelValue string,
	removeLabel bool,
) {
	logInfo("====== Additionnal step ====================================")
	logInfo("8. Add the label as a selector in the service...")
	// Get the current service
	currentService, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), applicationName, metav1.GetOptions{})
	check(err)


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
	check(updateServiceError)

	logSuccess("8. Service configured")
	logInfo("============================================================")
}
