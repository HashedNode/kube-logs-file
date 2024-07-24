package main

import (
	"context"
	"flag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
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

	namespace := flag.String("namespace", "", "Kubernetes namespace")
	podName := flag.String("pod", "", "Kubernetes pod name")
	
	flag.Parse()

	//use current contex

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %s", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error building clientset: %s", err)
	}

	if *namespace == "" {
		log.Fatalf("Namespace not specified")
	}

	if *podName == "" {
		log.Fatalf("Pod name not specified")
	}

	podlist, err := clientset.CoreV1().Pods(*namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing pods: %s", err)
	}

	for _, pod := range podlist.Items {
		if strings.Contains(pod.Name, *podName) {
			log.Printf("Found pod %s", pod.Name)
		}
	}

	/*
		req := clientset.CoreV1().Pods(namespace).GetLogs(pod, &v1.PodLogOptions{})

		podLogs, err := req.Stream(context.Background())
		if err != nil {
			log.Fatalf("Error opening stream: %s", err)
		}
		defer podLogs.Close()

		logs, err := io.ReadAll(podLogs)
		if err != nil {
			log.Fatalf("Error reading pod logs: %s", err)
		}

		logFileName := strings.Join([]string{pod, "log"}, ".")
		err = os.WriteFile(logFileName, logs, 0644)
		if err != nil {
			log.Fatalf("Error writing log file: %s", err)
		}
		fmt.Printf("Logs written to %s\n", logFileName)*/
}
