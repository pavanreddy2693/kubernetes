package controller

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	NamespaceEnv      = "NAMESPACE"
	ScalingThreshold  = 70 // Percentage
	ScaleFactor       = 2  // Scale factor when resizing
	CheckInterval     = 30 * time.Second
)

func StartPVCScalerController(config *rest.Config) error {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	namespace := os.Getenv(NamespaceEnv)
	if namespace == "" {
		namespace = "default"
	}

	log.Printf("Monitoring PVCs in namespace: %s", namespace)

	for {
		err := monitorAndScalePVCs(clientset, namespace)
		if err != nil {
			log.Printf("Error monitoring PVCs: %v", err)
		}
		time.Sleep(CheckInterval)
	}
}

func monitorAndScalePVCs(clientset *kubernetes.Clientset, namespace string) error {
	pvcs, err := clientset.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list PVCs: %w", err)
	}

	for _, pvc := range pvcs.Items {
		usagePercentage, err := getPVCUsagePercentage(clientset, pvc)
		if err != nil {
			log.Printf("Error getting usage for PVC %s: %v", pvc.Name, err)
			continue
		}

		if usagePercentage > ScalingThreshold {
			newSize := scalePVCSize(pvc)
			log.Printf("Scaling PVC %s from %s to %s", pvc.Name, pvc.Spec.Resources.Requests.Storage().String(), newSize.String())
			err := resizePVC(clientset, &pvc, newSize)
			if err != nil {
				log.Printf("Error resizing PVC %s: %v", pvc.Name, err)
			}
		}
	}

	return nil
}

func getPVCUsagePercentage(clientset *kubernetes.Clientset, pvc v1.PersistentVolumeClaim) (int, error) {
	// This function would use metrics from the PVC's corresponding PV or CSI driver
	// Placeholder logic, actual implementation depends on storage solution
	return 75, nil // Simulated usage percentage for demo purposes
}

func scalePVCSize(pvc v1.PersistentVolumeClaim) *resource.Quantity {
	originalSize := pvc.Spec.Resources.Requests.Storage()
	newSize := originalSize.DeepCopy()
	newSize.Set(newSize.Value() * ScaleFactor)
	return &newSize
}

func resizePVC(clientset *kubernetes.Clientset, pvc *v1.PersistentVolumeClaim, newSize *resource.Quantity) error {
	pvc.Spec.Resources.Requests[v1.ResourceStorage] = *newSize
	_, err := clientset.CoreV1().PersistentVolumeClaims(pvc.Namespace).Update(context.TODO(), pvc, metav1.UpdateOptions{})
	return err
}

