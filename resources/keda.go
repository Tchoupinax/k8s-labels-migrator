package keda

import (
	"context"
	"fmt"
	"time"

	"github.com/Tchoupinax/k8s-labels-migrator/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

// https://keda.sh/docs/2.13/concepts/scaling-deployments

func PauseScaledObject(
	crdClient *dynamic.DynamicClient,
	clientset *kubernetes.Clientset,
	deploymentName string,
	namespace string,
) {
	crdGVR := schema.GroupVersionResource{
		Group:    "keda.sh",
		Version:  "v1alpha1",
		Resource: "scaledobjects",
	}
	kedaScaledObject, _ := crdClient.Resource(crdGVR).Namespace(namespace).Get(context.TODO(), deploymentName, v1.GetOptions{})
	if kedaScaledObject != nil {
		utils.LogInfo("2.2 Keda Scaled Object detected")
		// Add the annotation "autoscaling.keda.sh/paused"
		kedaScaledObject.Object["metadata"].(map[string]interface{})["annotations"].(map[string]interface{})["autoscaling.keda.sh/paused"] = "true"
		_, updateKedaError := crdClient.Resource(crdGVR).Namespace(namespace).Update(context.TODO(), kedaScaledObject, v1.UpdateOptions{})
		utils.Check(updateKedaError)

		utils.LogSuccess("Keda object paused ⏸️")
		utils.LogBlocking("Waiting randomly 5 seconds to ensure keda controller registered the update")
		for i := 1; i < 5; i++ {
			utils.LogBlockingDot()
			time.Sleep(1 * time.Second)
		}
		fmt.Println()
	} else {
		utils.LogInfo("2.2 Any Keda Scaled Object detected")
	}
}

func ResumeScaledObject(
	crdClient *dynamic.DynamicClient,
	clientset *kubernetes.Clientset,
	deploymentName string,
	namespace string,
) {
	crdGVR := schema.GroupVersionResource{
		Group:    "keda.sh",
		Version:  "v1alpha1",
		Resource: "scaledobjects",
	}
	kedaScaledObject, _ := crdClient.Resource(crdGVR).Namespace(namespace).Get(context.TODO(), deploymentName, v1.GetOptions{})
	if kedaScaledObject != nil {
		utils.LogInfo("2.2 Keda Scaled Object detected")
		delete(
			kedaScaledObject.Object["metadata"].(map[string]interface{})["annotations"].(map[string]interface{}),
			"autoscaling.keda.sh/paused",
		)
		utils.LogSuccess("Keda object resumed ▶️")
		_, updateKedaError := crdClient.Resource(crdGVR).Namespace(namespace).Update(context.TODO(), kedaScaledObject, v1.UpdateOptions{})
		utils.Check(updateKedaError)
	} else {
		utils.LogInfo("2.2 Any Keda Scaled Object detected")
	}
}
