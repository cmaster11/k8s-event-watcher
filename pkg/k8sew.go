package pkg

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

var log = logrus.StandardLogger()

type K8sEventWatcher struct {
	config     *Config
	launchTime v12.Time

	kubeInformerFactory informers.SharedInformerFactory

	chStop   chan struct{}
	lock     sync.Mutex
	callback func(event *v1.Event, eventFilter *EventFilter, matchResult *MatchResult)

	Debug bool
}

func NewK8sEventWatcher(
	// Config of event watcher
	config *Config,
	// Config path for k8s cluster, can be empty
	kubeConfigPath *string,
	logWriter io.Writer,
) (*K8sEventWatcher, error) {
	var err error

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	var k8sConfig *rest.Config
	if kubeConfigPath != nil {
		k8sConfig, err = clientcmd.BuildConfigFromFlags("", *kubeConfigPath)
		if err != nil {
			return nil, fmt.Errorf("failed to build k8s config: %w", err)
		}
	} else {
		k8sConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to build in-cluster k8s config: %w", err)
		}
	}
	clientSet, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize k8s client: %w", err)
	}

	launchTime := v12.Now()
	kubeInformerFactory := informers.NewSharedInformerFactory(clientSet, time.Second*30)
	evtInformer := kubeInformerFactory.Core().V1().Events().Informer()

	watcher := &K8sEventWatcher{
		config:              config,
		launchTime:          launchTime,
		kubeInformerFactory: kubeInformerFactory,
	}

	evtInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: watcher.onAddEvent,
	})

	return watcher, nil
}

func (w *K8sEventWatcher) onAddEvent(obj interface{}) {
	evt, ok := obj.(*v1.Event)
	if !ok {
		log.WithField("object", obj).Errorf("failed to cast event")
		return
	}

	if w.config.SinceNow && evt.CreationTimestamp.Before(&w.launchTime) {
		// Old event
		log.WithField("event", evt).Debug("discarded old event")
		return
	}

	// Convert the event to a map
	outMap, err := eventToMap(evt)
	if err != nil {
		log.WithError(err).Error("failed to cast event to map")
		return
	}

	filter, matchResult, err := w.config.MatchingEventFilter(outMap)
	if err != nil {
		log.WithError(err).Error("failed to find matching event filter")
		return
	}
	if filter != nil {
		log.WithField("event", evt).Debug("matched event")
		w.callback(evt, filter, matchResult)
		return
	}

	log.WithField("event", evt).Debug("discarded event")
}

func (w *K8sEventWatcher) Start(callback func(event *v1.Event, eventFilter *EventFilter, matchResult *MatchResult)) error {
	if callback == nil {
		return fmt.Errorf("callback cannot be null")
	}

	w.lock.Lock()
	defer w.lock.Unlock()

	if w.chStop != nil {
		return fmt.Errorf("already started")
	}

	w.callback = callback
	w.chStop = make(chan struct{})

	go w.kubeInformerFactory.Start(w.chStop)

	return nil
}

func (w *K8sEventWatcher) Stop() {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.chStop == nil {
		return
	}

	close(w.chStop)
	w.chStop = nil
}

func (w *K8sEventWatcher) LaunchTime() v12.Time {
	return w.launchTime
}
