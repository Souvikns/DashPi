package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CPUStats struct {
	user, nice, system, idle, iowait, irq, softirq, steal uint64
}

func GetCPUStats() (CPUStats, error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return CPUStats{}, fmt.Errorf("failed to read /proc/stat: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "cpu ") {
			fields := strings.Fields(line)[1:] // Skip the "cpu" prefix
			if len(fields) < 8 {
				return CPUStats{}, fmt.Errorf("unexpected format in /proc/stat")
			}

			// Parse fields into uint64 values
			values := make([]uint64, len(fields))
			for i, field := range fields {
				values[i], err = strconv.ParseUint(field, 10, 64)
				if err != nil {
					return CPUStats{}, fmt.Errorf("failed to parse field %q: %w", field, err)
				}
			}

			return CPUStats{
				user:    values[0],
				nice:    values[1],
				system:  values[2],
				idle:    values[3],
				iowait:  values[4],
				irq:     values[5],
				softirq: values[6],
				steal:   values[7],
			}, nil
		}
	}

	return CPUStats{}, fmt.Errorf("no CPU stats found in /proc/stat")
}

func CalculateCPULoad(prev, curr CPUStats) float64 {
	// Total time is the sum of all CPU times
	prevTotal := prev.user + prev.nice + prev.system + prev.idle + prev.iowait + prev.irq + prev.softirq + prev.steal
	currTotal := curr.user + curr.nice + curr.system + curr.idle + curr.iowait + curr.irq + curr.softirq + curr.steal

	// Idle time
	prevIdle := prev.idle + prev.iowait
	currIdle := curr.idle + curr.iowait

	// Calculate deltas
	deltaTotal := currTotal - prevTotal
	deltaIdle := currIdle - prevIdle

	// Calculate CPU load
	if deltaTotal == 0 {
		return 0.0
	}

	return float64(deltaTotal-deltaIdle) / float64(deltaTotal) * 100.0
}
