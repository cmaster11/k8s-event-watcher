package k8seventwatcher

import (
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8seventwatcher/internal"
	"log"
	"sync"
	"time"
)

type K8sEventWatcher struct {
	config     *internal.Config
	launchTime v12.Time
	logger     *log.Logger

	kubeInformerFactory informers.SharedInformerFactory

	chStop   chan struct{}
	lock     sync.Mutex
	callback func(event *v1.Event)

	Debug bool
}

func NewK8sEventWatcher(
// Config path of event watcher
	configPath string,
// Config path for k8s cluster, can be empty
	kubeConfigPath *string,
	logWriter io.Writer,
) (*K8sEventWatcher, error) {
	if configPath == "" {
		return nil, errorf("empty config path provided")
	}

	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, errorf("failed to read Config file: %v", err)
	}

	config := &internal.Config{}
	err = yaml.Unmarshal(configData, config)
	if err != nil {
		return nil, errorf("failed to unmarshal Config: %v", err)
	}

	var k8sConfig *rest.Config
	if kubeConfigPath != nil {
		k8sConfig, err = clientcmd.BuildConfigFromFlags("", *kubeConfigPath)
		if err != nil {
			return nil, errorf("failed to build k8s config: %v", err)
		}
	} else {
		k8sConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, errorf("failed to build in-cluster k8s config: %v", err)
		}
	}
	clientSet, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, errorf("failed to initialize k8s client: %v", err)
	}

	launchTime := v12.Now()
	kubeInformerFactory := informers.NewSharedInformerFactory(clientSet, time.Second*30)
	evtInformer := kubeInformerFactory.Core().V1().Events().Informer()

	var logger *log.Logger

	if logWriter != nil {
		logger = log.New(logWriter, "", log.LstdFlags)
	}

	watcher := &K8sEventWatcher{
		config:              config,
		launchTime:          launchTime,
		logger:              logger,
		kubeInformerFactory: kubeInformerFactory,
	}

	evtInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: watcher.onAddEvent,
	}, )

	return watcher, nil
}

func (w *K8sEventWatcher) onAddEvent(obj interface{}) {
	evt, ok := obj.(*v1.Event)
	if !ok {
		w.logEntryError("failed to cast event: %+v", obj)
		return
	}

	if w.config.SinceNow && evt.CreationTimestamp.Before(&w.launchTime) {
		// Old event
		w.logEntryDebug("discarded old event: %+v", evt)
		return
	}

	if !w.config.MatchesEvent(evt) {
		w.logEntryDebug("discarded event: %+v", evt)
		return
	}

	w.logEntryDebug("matched event: %+v", evt)

	w.callback(evt)
}

func (w *K8sEventWatcher) Start(callback func(event *v1.Event)) error {
	if callback == nil {
		return errorf("callback cannot be null")
	}

	w.lock.Lock()
	defer w.lock.Unlock()

	if w.chStop != nil {
		return errorf("already started")
	}

	w.callback = callback
	w.chStop = make(chan struct{})

	go w.kubeInformerFactory.Start(w.chStop)

	w.logEntryInfo("started (%s) with Config:\n%s", w.launchTime.String(), w.config.Dump())

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
