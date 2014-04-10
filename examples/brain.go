package main

import (
	"encoding/json"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-neurosky"
)

var Master *gobot.Master = gobot.GobotMaster()
var eeg gobotNeurosky.EEG

func Eeg(params map[string]interface{}) string {
	b, _ := json.Marshal(eeg)
	return string(b)
}

func main() {

	gobot.Api(Master)

	adaptor := new(gobotNeurosky.NeuroskyAdaptor)
	adaptor.Name = "neurosky"
	adaptor.Port = "/dev/rfcomm0"

	neuro := gobotNeurosky.NewNeurosky(adaptor)
	neuro.Name = "neuro"

	work := func() {
		gobot.On(neuro.Events["EEG"], func(data interface{}) {
			eeg = data.(gobotNeurosky.EEG)
		})
	}

	Master.Robots = append(Master.Robots, gobot.Robot{
		Name:        "brain",
		Connections: []gobot.Connection{adaptor},
		Devices:     []gobot.Device{neuro},
		Work:        work,
		Commands:    map[string]interface{}{"Eeg": Eeg},
	})

	Master.Start()
}
