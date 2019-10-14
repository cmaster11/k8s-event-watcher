package k8seventwatcher

import "fmt"

func (w *K8sEventWatcher) logEntryDebug(format string, args ...interface{}) {
	if ! w.Debug {
		return
	}
	if w.logger != nil {
		w.logger.Printf("[K8EW] (DEBUG): "+format+"\n", args...)
	}
}

func (w *K8sEventWatcher) logEntryInfo(format string, args ...interface{}) {
	if w.logger != nil {
		w.logger.Printf("[K8EW] (INFO): "+format+"\n", args...)
	}
}

func (w *K8sEventWatcher) logEntryError(format string, args ...interface{}) {
	if w.logger != nil {
		w.logger.Printf("[K8EW] (ERROR): "+format+"\n", args...)
	}
}

func errorf(format string, args ...interface{}) error {
	return fmt.Errorf("[K8EW]: "+format+"\n", args)
}
