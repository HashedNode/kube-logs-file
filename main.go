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
	"sync"
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
	podNames := flag.String("pods", "", "Comma-separated list of pods names with no spaces")

	flag.Parse()

	if *namespace == "" {
		log.Fatalf("Namespace not specified")
	}

	if *podNames == "" {
		log.Fatalf("Pod names not specified")
	}

	podsNamesList := strings.Split(*podNames, ",")

	var (
		config *restclient.Config
		err    error
	)

	//use current contex
	if contextName != nil && *contextName != "" {
		log.Printf("Switching Kubernetes context to %s", *contextName)
		configOverrides := &clientcmd.ConfigOverrides{CurrentContext: *contextName}
		config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(&clientcmd.ClientConfigLoadingRules{ExplicitPath: *kubeconfig}, configOverrides).ClientConfig()
	} else {
		log.Printf("Using default Kubernetes context")
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	}

	if err != nil {
		log.Fatalf("Error building kubeconfig: %s", err)
	}

	clientSet, err := kubernetes.NewForConfig(config)

	if err != nil {
		log.Fatalf("Error building clientSet: %s", err)
	}

	podlist, err := clientSet.CoreV1().Pods(*namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing pods: %s", err)
	}

	var wg sync.WaitGroup

	for _, podName := range podsNamesList {
		wg.Add(1)
		go func(podName string) {
			defer wg.Done()

			log.Printf("Start search pod %s", podName)

			var realPod v1.Pod

			for _, pod := range podlist.Items {
				if strings.Contains(pod.Name, strings.TrimSpace(podName)) {
					realPod = pod
					break
				}
			}

			if realPod.Name != "" {
				req := clientSet.CoreV1().Pods(*namespace).GetLogs(realPod.Name, &v1.PodLogOptions{})
				getAndWriteLogs(req, realPod)
			} else {
				log.Printf("Pod with name %s not found", podName)
			}

		}(strings.TrimSpace(podName))
	}
	wg.Wait()
	fmt.Printf("finished download logs")

}

func getAndWriteLogs(req *restclient.Request, realPod v1.Pod) {
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
