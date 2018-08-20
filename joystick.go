package main

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/joystick"
)

func main() {
	joystickAdaptor := joystick.NewAdaptor()
	stick := joystick.NewDriver(joystickAdaptor, "/home/pi/go/src/gobot.io/x/gobot/platforms/joystick/configs/magicseer1.json",
	)

	work := func() {
		// buttons
		stick.On(joystick.APress, func(data interface{}) {
			fmt.Println("a Pressed")
		})
		stick.On(joystick.ARelease, func(data interface{}) {
			fmt.Println("a Released")
		})
		stick.On(joystick.BPress, func(data interface{}) {
			fmt.Println("b Pressed")
		})
		stick.On(joystick.BRelease, func(data interface{}) {
			fmt.Println("b Released")
		})
		stick.On("c_press", func(data interface{}) {
			fmt.Println("c Pressed")
		})
		stick.On("c_release", func(data interface{}) {
			fmt.Println("c Released")
		})
		stick.On("d_press", func(data interface{}) {
			fmt.Println("d Pressed")
		})
		stick.On("d_release", func(data interface{}) {
			fmt.Println("d Released")
		})


		// joysticks
		stick.On(joystick.LeftX, func(data interface{}) {
			fmt.Println("left_x", data)
		})
		stick.On(joystick.LeftY, func(data interface{}) {
			fmt.Println("left_y", data)
		})
		stick.On(joystick.RightX, func(data interface{}) {
			fmt.Println("right_x", data)
		})
		stick.On(joystick.RightY, func(data interface{}) {
			fmt.Println("right_y", data)
		})

		// triggers
		stick.On(joystick.R1Press, func(data interface{}) {
			fmt.Println("R1Press", data)
		})
		stick.On(joystick.R2Press, func(data interface{}) {
			fmt.Println("R2Press", data)
		})
		stick.On(joystick.L1Press, func(data interface{}) {
			fmt.Println("L1Press", data)
		})
		stick.On(joystick.L2Press, func(data interface{}) {
			fmt.Println("L2Press", data)
		})
	}

	robot := gobot.NewRobot("joystickBot",
		[]gobot.Connection{joystickAdaptor},
		[]gobot.Device{stick},
		work,
	)

	robot.Start()
}
