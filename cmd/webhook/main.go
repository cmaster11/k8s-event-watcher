package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/cmaster11/k8s-event-watcher/pkg"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var log = logrus.StandardLogger()

type CmdWebhook struct {
	httpClient *http.Client
	config     *Config
	watcher    *pkg.K8sEventWatcher
}

func main() {
	kubeConfigPath := flag.String("kubeconfig", "", "path of kubeconfig file to use")
	configPath := flag.String("config", "config.yaml", "path of event watcher config file to use")
	maxRetries := flag.Int("maxRetries", 3, "max retries for webhook sender")
	webhookTimeout := flag.Duration("webhookTimeout", 10*time.Second, "the max timeout when sending a webhook")

	flag.Parse()

	config, err := parseConfig(configPath)
	if err != nil {
		log.WithError(err).Fatal("failed to parse config")
	}

	watcher, err := pkg.NewK8sEventWatcher(
		&config.Config,
		kubeConfigPath,
		os.Stdout,
	)
	if err != nil {
		log.WithError(err).Fatal("failed to initialize event watcher")
	}

	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Timeout = *webhookTimeout
	retryClient.RetryMax = *maxRetries

	cmdWebhook := CmdWebhook{
		httpClient: retryClient.StandardClient(),
		config:     config,
		watcher:    watcher,
	}

	done := make(chan bool, 1)

	cmdWebhook.start()

	log.WithFields(logrus.Fields{
		"launchTime": watcher.LaunchTime(),
		"config":     config,
	}).Info("started")

	waitForShutdown(done)
	<-done

	log.Info("exiting")

	watcher.Stop()
}

func (c *CmdWebhook) start() {
	if err := c.watcher.Start(c.onEvent); err != nil {
		log.WithError(err).Fatal("failed to start watch")
	}
}

func (c *CmdWebhook) onEvent(event *v1.Event, eventFilter *pkg.EventFilter, matchResult *pkg.MatchResult) {
	fields := logrus.Fields{
		"event":       event,
		"eventFilter": eventFilter,
		"matchResult": matchResult,
	}

	log.WithFields(fields).Info("got event")

	bodyJSON, err := json.Marshal(fields)
	if err != nil {
		log.WithError(err).Error("failed to marshal event fields")
		return
	}

	for _, wh := range c.config.Webhooks {
		c.sendWebhook(wh, bodyJSON)
	}
}

func (c *CmdWebhook) sendWebhook(wh *Webhook, bodyJSON []byte) {
	req, err := http.NewRequest("POST", wh.Url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		log.WithError(err).Error("failed to initialize webhook request")
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(bodyJSON)))
	for key, val := range wh.Headers {
		req.Header.Add(key, val)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.WithError(err).Error("failed to send webhook request")
		return
	}
	defer resp.Body.Close()

	responseBodyBytes, _ := ioutil.ReadAll(resp.Body)

	logrus.WithFields(logrus.Fields{
		"statusCode":   resp.StatusCode,
		"responseBody": string(responseBodyBytes),
	}).Info("webhook sent")
}

func waitForShutdown(done chan bool) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()
}
