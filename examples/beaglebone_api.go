package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-beaglebone"
	"github.com/hybridgroup/gobot-gpio"
)

func main() {
	master := gobot.GobotMaster()
	api := gobot.Api(master)
	api.Port = "3001"
	beaglebone := new(gobotBeaglebone.Beaglebone)
	beaglebone.Name = "beaglebone"

	led := gobotGPIO.NewLed(beaglebone)
	led.Name = "led"
	led.Pin = "P9_12"

	work := func() {
		gobot.Every("1s", func() { led.Toggle() })
	}

	master.Robots = append(master.Robots, &gobot.Robot{
		Name:        "beaglebone",
		Connections: []gobot.Connection{beaglebone},
		Devices:     []gobot.Device{led},
		Work:        work,
	})

	master.Start()
}
