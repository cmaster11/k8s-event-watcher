package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cmaster11/k8s-event-watcher/pkg"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var log = logrus.StandardLogger()

func main() {
	kubeConfigPath := flag.String("kubeconfig", "", "path of kubeconfig file to use")
	configPath := flag.String("config", "config.yaml", "path of event watcher config file to use")

	flag.Parse()

	if *configPath == "" {
		log.Fatal("empty config path provided")
	}

	configData, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.WithError(err).Fatal("failed to read config file")
	}

	config := &pkg.Config{}
	err = yaml.Unmarshal(configData, config)
	if err != nil {
		log.WithError(err).Fatal("failed to unmarshal config")
	}

	watcher, err := pkg.NewK8sEventWatcher(
		config,
		kubeConfigPath,
		os.Stdout,
	)
	if err != nil {
		log.WithError(err).Fatal("failed to initialize k8s event watcher")
	}

	if err := watcher.Start(func(event *v1.Event, eventFilter *pkg.EventFilter, matchResult *pkg.MatchResult) {
		fields := logrus.Fields{
			"event":       event,
			"eventFilter": eventFilter,
			"matchResult": matchResult,
		}

		jsonBytes, _ := json.MarshalIndent(fields, "", "  ")

		log.Printf("--- Event ---\n%s\n", string(jsonBytes))
	}); err != nil {
		log.Fatal(err)
	}

	log.WithFields(logrus.Fields{
		"launchTime": watcher.LaunchTime(),
		"config":     config,
	}).Info("started")

	fmt.Println("Press 'Enter' to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	watcher.Stop()
}
