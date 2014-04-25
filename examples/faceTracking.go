package main

import (
	"fmt"
	cv "github.com/hybridgroup/go-opencv/opencv"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-ardrone"
	"github.com/hybridgroup/gobot-opencv"
	"math"
	"path"
	"runtime"
	"time"
)

func main() {
	_, currentfile, _, _ := runtime.Caller(0)
	cascade := path.Join(path.Dir(currentfile), "haarcascade_frontalface_alt.xml")

	opencv := new(gobotOpencv.Opencv)
	opencv.Name = "opencv"

	window := gobotOpencv.NewWindow(opencv)
	window.Name = "window"

	camera := gobotOpencv.NewCamera(opencv)
	camera.Name = "camera"
	camera.Source = "tcp://192.168.1.1:5555"

	ardroneAdaptor := new(gobotArdrone.ArdroneAdaptor)
	ardroneAdaptor.Name = "Drone"

	drone := gobotArdrone.NewArdrone(ardroneAdaptor)
	drone.Name = "Drone"

	work := func() {
		detect := false
		drone.TakeOff()
		var image *cv.IplImage
		gobot.On(camera.Events["Frame"], func(data interface{}) {
			image = data.(*cv.IplImage)
			if detect == false {
				window.ShowImage(image)
			}
		})
		gobot.On(drone.Events["Flying"], func(data interface{}) {
			gobot.After("1s", func() { drone.Up(0.2) })
			gobot.After("2s", func() { drone.Hover() })
			gobot.After("5s", func() {
				detect = true
				go func() {
					for {
						drone.Hover()
						i := image.Clone()
						faces := gobotOpencv.DetectFaces(cascade, i)
						biggest := 0
						var face *cv.Rect
						for _, f := range faces {
							if f.Width() > biggest {
								biggest = f.Width()
								face = f
							}
						}
						if face != nil && (face.Width() <= 100 && face.Width() >= 45) {
							gobotOpencv.DrawRectangles(i, []*cv.Rect{face}, 0, 255, 0, 5)
							center_x := float64(i.Width()) * 0.5
							turn := -(float64(face.X()) - center_x) / center_x
							fmt.Println("turning:", turn)
							if turn < 0 {
								drone.Clockwise(math.Abs(turn))
							} else {
								drone.CounterClockwise(math.Abs(turn))
							}
						}
						window.ShowImage(i)
						time.Sleep(40 * time.Millisecond)
					}
				}()
				gobot.After("20s", func() { drone.Land() })
			})
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{opencv, ardroneAdaptor},
		Devices:     []gobot.Device{window, camera, drone},
		Work:        work,
	}

	robot.Start()
}
