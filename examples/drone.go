package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-ardrone"
)

func main() {

	ardroneAdaptor := new(gobotArdrone.ArdroneAdaptor)
	ardroneAdaptor.Name = "Drone"

	drone := gobotArdrone.NewArdrone(ardroneAdaptor)
	drone.Name = "Drone"

	work := func() {
		drone.TakeOff()
		gobot.On(drone.Events["Flying"], func(data interface{}) {
			gobot.After("15s", func() {
				drone.Land()
			})
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{ardroneAdaptor},
		Devices:     []gobot.Device{drone},
		Work:        work,
	}

	robot.Start()
}
