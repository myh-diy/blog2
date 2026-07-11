# Mini Node Exporter

一个用于学习 Prometheus exporter 原理的小型 Linux 采集器，只采集 CPU 和内存。

它使用 Prometheus 官方 Go client：

- `github.com/prometheus/client_golang/prometheus`
- `github.com/prometheus/client_golang/prometheus/promhttp`

## 启动

```bash
go mod tidy
go run .
```

默认监听：

```text
http://localhost:9101/metrics
```

可以用环境变量修改监听地址：

```bash
EXPORTER_ADDR=:9101 go run .
```

## Docker

```bash
docker build -t mini-node-exporter .
docker run --rm -p 9101:9101 mini-node-exporter
```

如果你想采集宿主机的 `/proc`，后面可以再扩展成挂载宿主机 `/proc` 的模式。当前版本用于理解 exporter 和 Prometheus scrape 流程。

## 指标

```text
mini_node_cpu_usage_ratio
mini_node_memory_total_bytes
mini_node_memory_available_bytes
mini_node_memory_used_bytes
mini_node_memory_usage_ratio
```

CPU 使用率基于两次 scrape 之间的 CPU idle/total 差值计算，所以第一次采集通常是 `0`，第二次开始才有真实比例。

## Prometheus 配置

参考 `prometheus.example.yml`：

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: "mini-node-exporter"
    static_configs:
      - targets:
          - "localhost:9101"
```
