# k8s-event-watcher

This library lets you watch for Kubernetes events, and triggers a callback whenever one event matching the criteria is found.

An example of configuration is:

```yaml
sinceNow: true
filters:
- rules:
    involvedObject.kind: Job
    involvedObject.name: "^*.fail"
    reason: BackoffLimitExceeded
  errorRules:
    # Always mark this event as an error
    type: .*
```

This would match a `Job`-failed `Event`:

```yaml
apiVersion: v1
kind: Event
...
involvedObject:
  apiVersion: batch/v1
  kind: Job
  name: job-fail
  ...
lastTimestamp: "2019-10-14T07:11:45Z"
message: Job has reached the specified backoff limit
metadata:
  ...
reason: BackoffLimitExceeded
...
type: Warning
```

## Configuration

Configurable keys:

* `sinceNow`: if `true`, only processes events generated after the program starts.
* `filters`: a list of `rules`. Each `rules` object is evaluated **independently**, in an `OR` fashion. Any event is
    evaluated on all sets of `rules`. The first matching filter will cause the trigger of the callback. 
* `filters.[*].rules`: a map of regular expressions. Each regular expression is evaluated against the provided object key.
    If **all** the regular expressions match, then the event will be sent to the callback, and the `MatchedFields` 
    map will be populated with the matching fields.
* `filters.[*].errorRules`: a map of regular expressions. Each regular expression is evaluated against the provided object key.
    If **all** the regular expressions match, then the event will be considered an error, and the `MatchedErrorFields` 
    map will be populated with the matching error fields.

## Local test

You can test the execution of this library by running the example program:

```bash
go run ./example --kubeconfig PATH_TO_A_KUBECONFIG_FILE --config ./example/config.yaml
```

While the example program is running, you can then start a failing job with:

```bash
kubectl apply -f ./example/job-fail-k8s.yaml
```

The example program will then pick up the failure and show the matching event.

## Webhook sender

You can also use the webhook sender application (`./cmd/webhook`) to trigger a webhook whenever an event is matched.