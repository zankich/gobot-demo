package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-firmata"
	"github.com/hybridgroup/gobot-gpio"
	"github.com/hybridgroup/gobot-sphero"
	"time"
)

func main() {
	firmata := new(gobotFirmata.FirmataAdaptor)
	firmata.Name = "firmata"
	firmata.Port = "/dev/ttyACM0"

	led := gobotGPIO.NewLed(firmata)
	led.Name = "led"
	led.Pin = "3"

	spheroAdaptor := new(gobotSphero.SpheroAdaptor)
	spheroAdaptor.Name = "Sphero"
	spheroAdaptor.Port = "/dev/rfcomm0"

	sphero := gobotSphero.NewSphero(spheroAdaptor)
	sphero.Name = "sphero"

	work := func() {
		sphero.SetRGB(0, 255, 0)
		led.Brightness(0)
		c := make(chan interface{})

		gobot.On(sphero.Events["Collision"], func(data interface{}) {
			gobot.Publish(c, data)
		})

		gobot.Every("2s", func() {
			sphero.Roll(90, uint16(gobot.Rand(360)))
		})

		gobot.On(c, func(data interface{}) {
			sphero.SetRGB(255, 0, 0)
			brightness := uint8(1)
			fade_amount := uint8(5)

			for brightness != 0 {
				time.Sleep(10 * time.Millisecond)
				led.Brightness(brightness)
				brightness = brightness + fade_amount
				if brightness == 0 || brightness == 255 {
					fade_amount = -fade_amount
				}
			}
			sphero.SetRGB(0, 255, 0)
			led.Brightness(0)
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{firmata, spheroAdaptor},
		Devices:     []gobot.Device{led, sphero},
		Work:        work,
	}

	robot.Start()
}
