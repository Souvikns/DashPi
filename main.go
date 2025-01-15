package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zcalusic/sysinfo"
)

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
