package main

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/zcalusic/sysinfo"
	"log"
	"os/user"
)

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
		if err != nil {
			log.Fatal(err)
		}
		return ctx.SendString(string(u))
	})

	app.Listen(":8080")
}
