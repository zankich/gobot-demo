package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-digispark"
	"github.com/hybridgroup/gobot-gpio"
)

func main() {
	master := gobot.GobotMaster()
	gobot.Api(master)

	digispark := new(gobotDigispark.DigisparkAdaptor)
	digispark.Name = "digispark"

	red := gobotGPIO.NewLed(digispark)
	red.Name = "red"
	red.Pin = "0"

	blue := gobotGPIO.NewLed(digispark)
	blue.Name = "blue"
	blue.Pin = "2"

	work := func() {
		gobot.Every("1s", func() {
			red.Toggle()
		})
	}

	master.Robots = append(master.Robots, gobot.Robot{
		Name:        "digispark",
		Connections: []gobot.Connection{digispark},
		Devices:     []gobot.Device{red, blue},
		Work: work,
	})

	master.Start()
}