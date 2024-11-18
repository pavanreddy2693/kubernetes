package main

import (
	"log"
	"os"

	"github.com/pavanreddy2693/pvc-scaler-controller/pkg/controller"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Get the in-cluster config or kubeconfig
	kubeconfig := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Failed to load kubeconfig: %v", err)
	}

	// Start the PVC scaler controller
	if err := controller.StartPVCScalerController(config); err != nil {
		log.Fatalf("Error starting PVC scaler controller: %v", err)
	}
}

