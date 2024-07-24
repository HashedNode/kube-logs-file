package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	contextName := flag.String("context", "", "Kubernetes context name")
	namespace := flag.String("namespace", "", "Kubernetes namespace")
	podName := flag.String("pod", "", "Kubernetes pod name")

	flag.Parse()

	if *namespace == "" {
		log.Fatalf("Namespace not specified")
	}

	if *podName == "" {
		log.Fatalf("Pod name not specified")
	}

	var (
		config *restclient.Config
		err    error
	)

	//use current contex
	if contextName != nil && *contextName != "" {
		configOverrides := &clientcmd.ConfigOverrides{CurrentContext: *contextName}
		config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(&clientcmd.ClientConfigLoadingRules{ExplicitPath: *kubeconfig}, configOverrides).ClientConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	}

	if err != nil {
		log.Fatalf("Error building kubeconfig: %s", err)
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		log.Fatalf("Error building clientset: %s", err)
	}

	podlist, err := clientset.CoreV1().Pods(*namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing pods: %s", err)
	}

	var realPod v1.Pod

	for _, pod := range podlist.Items {
		if strings.Contains(pod.Name, *podName) {
			log.Printf("Found pod %s", pod.Name)
			realPod = pod
			break
		}
	}

	if realPod.Name == "" {
		log.Fatalf("No Pod found with name %s", *podName)
	}

	req := clientset.CoreV1().Pods(*namespace).GetLogs(realPod.Name, &v1.PodLogOptions{})
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		log.Fatalf("Error opening stream: %s \n for pod %s", err, realPod.Name)
	}

	defer podLogs.Close()

	logs, err := io.ReadAll(podLogs)

	if err != nil {
		log.Fatalf("Error reading pod logs: %s", err)
	}

	fileName := strings.Join([]string{realPod.Name, "log"}, ".")
	err = os.WriteFile(fileName, logs, 0644)
	if err != nil {
		log.Fatalf("Error writing file: %s", err)
	}

	fmt.Printf("Logs written to %s \n", fileName)

}
