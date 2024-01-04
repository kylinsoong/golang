package main


import (
    "context"
    "fmt"
    "os"
    "time"
    
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/tools/record"
)

func main() {

    config, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
    if err != nil {
        fmt.Printf("Error building kubeconfig: %v\n", err)
        os.Exit(1)
    }
   
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        fmt.Printf("Error creating Kubernetes client: %v\n", err)
        os.Exit(1)
    }

    eventBroadcaster := record.NewBroadcaster()
    eventBroadcaster.StartLogging(func(format string, args ...interface{}) {
        fmt.Printf(format+"\n", args...)
    })
    eventBroadcaster.StartRecordingToSink(&corev1.EventSinkImpl{
        Interface: clientset.CoreV1().Events("default"),
    })

    recorder := eventBroadcaster.NewRecorder(config, corev1.EventSource{Component: "example-controller"})

    pod := &corev1.Pod{
        ObjectMeta: metav1.ObjectMeta{
            Name:      "example-pod",
	    Namespace: "default",
        },
    }

    message := "This is a sample event"
    recorder.Event(pod, corev1.EventTypeNormal, "ExampleEvent", message)

    time.Sleep(2 * time.Second)

    events, err := clientset.CoreV1().Events("default").List(context.TODO(), metav1.ListOptions{
	FieldSelector: fmt.Sprintf("involvedObject.name=%s", pod.Name),
    })
    if err != nil {
	fmt.Printf("Error listing events: %v\n", err)
	os.Exit(1)
    }

    fmt.Println("Events for the pod:")
    for _, event := range events.Items {
        fmt.Printf(" - %s: %s\n", event.Reason, event.Message)
    }

}

