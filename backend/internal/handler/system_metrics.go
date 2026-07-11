package handler

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type SystemMetrics struct {
	CPUUsageRatio        float64 `json:"cpu_usage_ratio"`
	MemoryTotalBytes     float64 `json:"memory_total_bytes"`
	MemoryAvailableBytes float64 `json:"memory_available_bytes"`
	MemoryUsedBytes      float64 `json:"memory_used_bytes"`
	MemoryUsageRatio     float64 `json:"memory_usage_ratio"`
}

func GetSystemMetrics(exporterURL string) gin.HandlerFunc {
	client := &http.Client{Timeout: 3 * time.Second}

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, exporterURL, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "mini exporter unavailable"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("mini exporter returned %d", resp.StatusCode)})
			return
		}

		metrics, err := parseSystemMetrics(resp.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"metrics": metrics})
	}
}

func parseSystemMetrics(input interface{ Read([]byte) (int, error) }) (SystemMetrics, error) {
	values := map[string]float64{}
	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		name := strings.SplitN(fields[0], "{", 2)[0]
		switch name {
		case "mini_node_cpu_usage_ratio",
			"mini_node_memory_total_bytes",
			"mini_node_memory_available_bytes",
			"mini_node_memory_used_bytes",
			"mini_node_memory_usage_ratio":
			value, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return SystemMetrics{}, err
			}
			values[name] = value
		}
	}
	if err := scanner.Err(); err != nil {
		return SystemMetrics{}, err
	}

	return SystemMetrics{
		CPUUsageRatio:        values["mini_node_cpu_usage_ratio"],
		MemoryTotalBytes:     values["mini_node_memory_total_bytes"],
		MemoryAvailableBytes: values["mini_node_memory_available_bytes"],
		MemoryUsedBytes:      values["mini_node_memory_used_bytes"],
		MemoryUsageRatio:     values["mini_node_memory_usage_ratio"],
	}, nil
}
