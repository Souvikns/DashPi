package controller

import (
	"github.com/Souvikns/DashPi/services"
	"github.com/gofiber/fiber/v2"
)

type Reponse struct {
	Uptime      string
	MemoryUsage MemoryStats
	CPULoad     float64
}

type MemoryStats struct {
	Total      uint64
	Free       uint64
	Available  uint64
	Percentage float64
}

func GetSystemInfo(ctx *fiber.Ctx) error {
	uptime, err := services.GetUptime()
	if err != nil {
		ctx.SendStatus(404)
	}

	memoryUsage, err := services.GetMemoryStats()
	if err != nil {
		ctx.SendStatus(404)
	}

	calpercent := services.CalculateRAMUsage(memoryUsage)

	load := <-services.CalcCPULoad()

	response := Reponse{Uptime: uptime, MemoryUsage: MemoryStats{
		Total:      memoryUsage.Total / 1000000,
		Free:       memoryUsage.Free / 1000000,
		Available:  memoryUsage.Available / 1000000,
		Percentage: calpercent,
	}, CPULoad: load}

	return ctx.JSON(&response)
}
