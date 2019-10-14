package main

import (
	"bufio"
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8seventwatcher"
	"log"
	"os"
)

func main() {
	kubeConfigPath := flag.String("kubeconfig", "", "path of k8s k8sConfig file to use")
	configPath := flag.String("config", "config.yaml", "path of k8sConfig file to use")

	flag.Parse()

	watcher, err := k8seventwatcher.NewK8sEventWatcher(
		*configPath,
		kubeConfigPath,
		os.Stdout,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := watcher.Start(func(event *v1.Event) {
		log.Printf("got event: %+v\n", event)
	}); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Press 'Enter' to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	watcher.Stop()
}
