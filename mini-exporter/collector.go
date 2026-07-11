package main

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type NodeCollector struct {
	mu       sync.Mutex
	prevCPU  cpuTimes
	hasPrev  bool
	cpuUsage *prometheus.Desc
	memTotal *prometheus.Desc
	memFree  *prometheus.Desc
	memUsed  *prometheus.Desc
	memUsage *prometheus.Desc
}

type cpuTimes struct {
	idle  uint64
	total uint64
}

type memoryStats struct {
	total     uint64
	available uint64
}

func NewNodeCollector() *NodeCollector {
	labels := prometheus.Labels{"collector": "mini_node"}
	return &NodeCollector{
		cpuUsage: prometheus.NewDesc(
			"mini_node_cpu_usage_ratio",
			"CPU usage ratio calculated from idle and total CPU time since the previous scrape.",
			nil,
			labels,
		),
		memTotal: prometheus.NewDesc(
			"mini_node_memory_total_bytes",
			"Total physical memory in bytes.",
			nil,
			labels,
		),
		memFree: prometheus.NewDesc(
			"mini_node_memory_available_bytes",
			"Available physical memory in bytes.",
			nil,
			labels,
		),
		memUsed: prometheus.NewDesc(
			"mini_node_memory_used_bytes",
			"Used physical memory in bytes.",
			nil,
			labels,
		),
		memUsage: prometheus.NewDesc(
			"mini_node_memory_usage_ratio",
			"Memory usage ratio calculated as used divided by total.",
			nil,
			labels,
		),
	}
}

func (c *NodeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.cpuUsage
	ch <- c.memTotal
	ch <- c.memFree
	ch <- c.memUsed
	ch <- c.memUsage
}

func (c *NodeCollector) Collect(ch chan<- prometheus.Metric) {
	if cpu, err := readCPUTimes(); err == nil {
		ch <- prometheus.MustNewConstMetric(c.cpuUsage, prometheus.GaugeValue, c.cpuUsageRatio(cpu))
	} else {
		ch <- prometheus.NewInvalidMetric(c.cpuUsage, err)
	}

	mem, err := readMemoryStats()
	if err != nil {
		ch <- prometheus.NewInvalidMetric(c.memTotal, err)
		return
	}

	used := mem.total - mem.available
	usage := 0.0
	if mem.total > 0 {
		usage = float64(used) / float64(mem.total)
	}

	ch <- prometheus.MustNewConstMetric(c.memTotal, prometheus.GaugeValue, float64(mem.total))
	ch <- prometheus.MustNewConstMetric(c.memFree, prometheus.GaugeValue, float64(mem.available))
	ch <- prometheus.MustNewConstMetric(c.memUsed, prometheus.GaugeValue, float64(used))
	ch <- prometheus.MustNewConstMetric(c.memUsage, prometheus.GaugeValue, usage)
}

func (c *NodeCollector) cpuUsageRatio(current cpuTimes) float64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.hasPrev {
		c.prevCPU = current
		c.hasPrev = true
		return 0
	}

	totalDelta := current.total - c.prevCPU.total
	idleDelta := current.idle - c.prevCPU.idle
	c.prevCPU = current

	if totalDelta == 0 || idleDelta > totalDelta {
		return 0
	}
	return float64(totalDelta-idleDelta) / float64(totalDelta)
}
