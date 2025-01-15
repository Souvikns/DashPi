package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zcalusic/sysinfo"
)

// RAM Usage 

type memoryStats struct {
	total     uint64
	free      uint64
	available uint64
}

func getMemoryStats() (memoryStats, error) {
	data, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return memoryStats{}, fmt.Errorf("failed to read /proc/meminfo: %w", err)
	}

	stats := memoryStats{}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		key := fields[0]
		value, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return memoryStats{}, fmt.Errorf("failed to parse value for %s: %w", key, err)
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
		return memoryStats{}, fmt.Errorf("failed to parse total memory")
	}

	return stats, nil
}

func calculateRAMUsage(stats memoryStats) float64 {
	usedMemory := stats.total - stats.available
	return float64(usedMemory) / float64(stats.total) * 100.0
}

//----

// CPU USAGE
type CPUStats struct {
	user, nice, system, idle, iowait, irq, softirq, steal uint64
}

func getCPUStats() (CPUStats, error) {
	data, err := ioutil.ReadFile("/proc/stat")
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

func calculateCPULoad(prev, curr CPUStats) float64 {
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

// --

// UPTIME

func getUptime()(string, error) {
	data , err := os.ReadFile("/proc/uptime")
	if err != nil {
		return "", fmt.Errorf("failed to read /proc/uptime: %w", err)
	}

	fields := strings.Fields(string(data))

	if len(fields) <1 {
		return "", fmt.Errorf("unexpected format in /proc/uptime")
	}

	uptimeSeconds, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return "", fmt.Errorf("fialed to parse uptime: %w", err)
	}

	uptime := time.Duration(uptimeSeconds) * time.Second
	days := uptime / (24* time.Hour)
	uptime %= 24 * time.Hour
	hours := uptime / time.Hour
	uptime %= time.Hour
	minutes := uptime / time.Minute
	uptime %= time.Minute
	seconds := uptime / time.Second

	return fmt.Sprintf("%d days, %d hours, %d minutes, %d seconds", days, hours, minutes, seconds), nil
}

func main() {
	current, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	if current.Uid != "0" {
		log.Fatal("requires superuser privilege")
	}

	var si sysinfo.SysInfo

	si.GetSysInfo()

	app := fiber.New()

	app.Static("/", "./web/dist")
	app.Get("/api/sysinfo", func(ctx *fiber.Ctx) error {
		u, err := json.MarshalIndent(&si, "", " ")
		if err !=nil {
			log.Fatal(err)
		}
		return ctx.SendString(string(u))
	})

	app.Listen(":8080")
}
