package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/cmaster11/k8s-event-watcher"
	"k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	kubeConfigPath := flag.String("kubeconfig", "", "path of kubeconfig file to use")
	configPath := flag.String("config", "config.yaml", "path of event watcher config file to use")

	flag.Parse()

	watcher, err := k8seventwatcher.NewK8sEventWatcher(
		*configPath,
		kubeConfigPath,
		os.Stdout,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := watcher.Start(func(event *v1.Event, eventFilter *k8seventwatcher.EventFilter, matchedFields map[string]interface{}) {
		log.Printf("got event (%s):\nmatchedFields: %+v\nevent: %+v\n", eventFilter.String(), matchedFields, event)
	}); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Press 'Enter' to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	watcher.Stop()
}
