package main

import (
	"flag"
	"log"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeConfigPath := flag.String("kubeconfig", "", "path of k8s config file to use")

	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("started")

	kubeInformerFactory := informers.NewSharedInformerFactory(clientSet, time.Second*30)
	evtInformer := kubeInformerFactory.Core().V1().Events().Informer()

	evtInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: onAddEvent,
	}, )

	stop := make(chan struct{})
	defer close(stop)

	go kubeInformerFactory.Start(stop)
	for {
		time.Sleep(time.Second)
	}
}

func onAddEvent(obj interface{})  {
	evt, ok := obj.(*v1.Event)
	if !ok {
		log.Printf("failed to cast event: %+v\n", obj)
		return
	}

	log.Printf("new event: :%+v\n", evt)
}
