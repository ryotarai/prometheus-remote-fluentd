# prometheus-remote-fluentd

[![Docker Repository on Quay](https://quay.io/repository/ryotarai/prometheus-remote-fluentd/status "Docker Repository on Quay")](https://quay.io/repository/ryotarai/prometheus-remote-fluentd)

```
+------------+                  +---------------------------+       +---------+
|            |                  |                           |       |         |
| Prometheus +------------------> prometheus-remote-fluentd +-------> Fluentd |
|            |                  |                           |       |         |
+------------+  Remote Storage  +---------------------------+       +---------+
                   (/write)
```

```
$ prometheus-remote-fluentd -fluent-host=localhost -fluent-tag=prometheus.samples
```
