package main

import (
	"log"
	"os/user"

	"github.com/Souvikns/DashPi/controller"
	"github.com/gofiber/fiber/v2"
	"github.com/zcalusic/sysinfo"
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
	app.Get("/api/sysinfo", controller.GetSystemInfo)

	app.Listen(":8080")
}
