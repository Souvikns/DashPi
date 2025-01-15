package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type MemoryStats struct {
	total     uint64
	free      uint64
	available uint64
}

func GetMemoryStats() (MemoryStats, error) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return MemoryStats{}, fmt.Errorf("failed to read /proc/meminfo: %w", err)
	}

	stats := MemoryStats{}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		key := fields[0]
		value, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return MemoryStats{}, fmt.Errorf("failed to parse value for %s: %w", key, err)
		}

		switch key {
		case "MemTotal:":
			stats.total = value
		case "MemFree:":
			stats.free = value
		case "MemAvailable:":
			stats.available = value
		}
	}

	if stats.total == 0 {
		return MemoryStats{}, fmt.Errorf("failed to parse total memory")
	}

	return stats, nil
}

func CalculateRAMUsage(stats MemoryStats) float64 {
	usedMemory := stats.total - stats.available
	return float64(usedMemory) / float64(stats.total) * 100.0
}
