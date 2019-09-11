---
title: bring your own prometheus
date: "2015-05-01T22:12:03.284Z"
---

Prometheus is a
By default, Istio deploys Prometheus into the `istio-system` namespace, and configures Prometheus to scrape metrics for both the Istio control plane and the sidecar proxies for your workloads.

But what if you're already running Prometheus on your Kubernetes cluster, and you've customized the install for scaling, or for alerts like Slack notifactions?

Not to worry. You can use your own Prometheus instance for Istio metrics. Let's see how.

Prometheus relies on a scrape configuration


But you might already have your own installation, configured for other things, or with add-ons like Slack notifications using Alert Manager

https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/



and enables the Prometheus adapter. This allows metrics to flow from your Envoy sidecar proxies to Prometheus






```bash
$ kubectl get service -n monitoring

NAME         TYPE           CLUSTER-IP   EXTERNAL-IP      PORT(S)          AGE
prometheus   LoadBalancer   10.0.3.155   <IP>             9090:32352/TCP   21m
```

