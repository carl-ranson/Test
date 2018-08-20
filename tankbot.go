package main

import (
	"time"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/raspi"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/api"
	"gobot.io/x/gobot/platforms/joystick"
	"fmt"
	"os"
	"../src/Tankbot"
)

const GPIO_01 = "28"
const GPIO_02 = "3"
const GPIO_03 = "5"
const GPIO_04 = "7"
const GPIO_05 = "29"
const GPIO_06 = "31"
const GPIO_07 = "26"
const GPIO_08 = "24"
const GPIO_09 = "21"
const GPIO_10 = "19"
const GPIO_11 = "23"
const GPIO_12 = "32"
const GPIO_13 = "33"
const GPIO_14 = "8"
const GPIO_15 = "10"
const GPIO_16 = "36"
const GPIO_17 = "11"
const GPIO_18 = "12"
const GPIO_19 = "35"
const GPIO_20 = "38"
const GPIO_21 = "40"
const GPIO_22 = "15"
const GPIO_23 = "16"
const GPIO_24 = "18"
const GPIO_25 = "22"
const GPIO_26 = "37"
const GPIO_27 = "13"



const M1EnablePin = GPIO_22
const M1DirectionPin = GPIO_24
const M1PWMPin = GPIO_12
const M2EnablePin = GPIO_23
const M2DirectionPin = GPIO_25
const M2PWMPin = GPIO_13
const LeftSensor = GPIO_16
const RightSensor = GPIO_17
const RedPin = GPIO_14
const GreenPin = GPIO_02
const BluePin = GPIO_03

const TurnRight = 1
const TurnLeft = 2
const Forward = 3
const Backward = 4
const Stop = 0

type Action interface {
	handleEvent(event string, data interface{}) 
	getState(state *State) 
	getEvents() []string
}

type ShutdownAction struct {
	enterDown bool
	backDown bool
} 
func (b *ShutdownAction) handleEvent(event string, data interface{}) {
	
	fmt.Println("Handle event", event)
	if event == "enter_press" {
		b.enterDown = true
	}
	if event == "enter_release" {
		b.enterDown = false
	}
	if event == "back_press" {
		b.backDown = true
	}
	if event == "back_release" {
		b.backDown = false
	}
}
	
func (b *ShutdownAction) getState(state *State) {
	if !state.shutdownSet && b.enterDown && b.backDown {
		state.shutdown = true
		state.shutdownSet = true
	}	
}
func (b *ShutdownAction) getEvents() []string {
	return []string{"enter_press","enter_release","back_press","back_release"}
}
type DriveAction struct {
	driveDir byte
	speed byte
	speedStr string
} 
func (b *DriveAction) handleEvent(event string, data interface{}) {
	if event == "left_x" {
		if data.(int16) < 0 {
			b.driveDir = TurnLeft
		} else if data.(int16) > 0 {
			b.driveDir = TurnRight
		} else {
			b.driveDir = Stop
		}
	}
	if event == "left_y" {
		if data.(int16) < 0 {
			b.driveDir = Forward
		} else if data.(int16) > 0 {
			b.driveDir = Backward
		} else {
			b.driveDir = Stop
		}
	}
	if event == "a_press" {
		b.speed = 64
		b.speedStr = "low speed"
	}
	if event == "b_press" {
		b.speed = 96
		b.speedStr = "mid speed"
	}
	if event == "c_press" {
		b.speed = 128
		b.speedStr = "high speed"
	}
	if event == "d_press" {
		b.speed = 255
		b.speedStr = "max speed"
	}
}
func (b *DriveAction) getState(state *State) {
	if !state.driveDirSet {
		state.driveDir = b.driveDir
		state.driveDirSet = true
	}
	if !state.speedSet {
		state.speed = b.speed
		state.speedStr = b.speedStr
		state.speedSet = true
	}
}
func (b *DriveAction) getEvents() []string {
	return []string{"left_x", "left_y", "a_press", "b_press" ,"c_press", "d_press"}
}

type State struct {
	r,g,b bool 
	rgbSet bool 
	driveDir byte
	driveDirSet bool
	speed byte
	speedStr string
	speedSet bool
	shutdown bool
	shutdownSet bool
}

type Behaviour struct {
	signalChan chan bool
	nextBehaviour *Behaviour
	name string
	eventCaught bool
	action Action
}

func NewBehaviour(name string, signalChan chan bool, action Action, nextBehaviour *Behaviour) *Behaviour {
	return &Behaviour{
		nextBehaviour: nextBehaviour,
		name: name,
		signalChan: signalChan, 
		action: action,
	}
}
func makeEventHandler(event string, b *Behaviour) func(data interface{}) {
	return func(data interface{}) {
		b.handleEvent(event, data) 
	}
}
func (b *Behaviour) registerEvents(t *Trackbot) {	
	fmt.Println("Register events for ", b.name)
	for _, event := range b.action.getEvents() {
		t.on(event, makeEventHandler(event, b))
	}
}
func (b *Behaviour) handleEvent(event string, data interface{}) {
	fmt.Println("HandleEvent", event, "on", b.name)
	b.action.handleEvent(event, data)
	b.eventCaught = true
	b.signalChan <- true
}
func (b *Behaviour) getState(state *State) {
	if b.eventCaught {
		b.action.getState(state)
	}
	if b.nextBehaviour != nil {
		b.nextBehaviour.getState(state)
	}
	b.eventCaught = false
	
	//fmt.Println("getState", "on", b.name, "state = ", state)
}

type Trackbot struct {
	leftMotor *Tankbot.MC33926Driver
	rightMotor *Tankbot.MC33926Driver
	leftSensor *gpio.DirectPinDriver
	rightSensor *gpio.DirectPinDriver
	joystickDriver *joystick.Driver
	redPin  *gpio.DirectPinDriver
	greenPin *gpio.DirectPinDriver
	bluePin *gpio.DirectPinDriver

	leftSensorValue int
	rightSensorValue int
	speed byte

}
func NewTrackbot(pi *raspi.Adaptor, joystickAdaptor *joystick.Adaptor) *Trackbot {
	return &Trackbot{
		leftMotor: Tankbot.NewMC33926Driver(pi, pi, M1EnablePin, M1DirectionPin, M1PWMPin),
		rightMotor: Tankbot.NewMC33926Driver(pi, pi, M2EnablePin, M2DirectionPin, M2PWMPin),
		leftSensor: gpio.NewDirectPinDriver(pi, LeftSensor),
		rightSensor: gpio.NewDirectPinDriver(pi, RightSensor),
		joystickDriver: joystick.NewDriver(joystickAdaptor, "/home/pi/projects/src/magicseer1.json"),
		redPin: gpio.NewDirectPinDriver(pi, RedPin),
		greenPin: gpio.NewDirectPinDriver(pi, GreenPin),
		bluePin: gpio.NewDirectPinDriver(pi, BluePin),

		speed: 48,
	}
}
func (t *Trackbot) out(msg string) {
	if msg != "" {
		fmt.Println(msg)
	}
}
func (t *Trackbot) SetSpeed(newSpeed byte, msg string) {
	t.speed = newSpeed
	t.out(msg)
}
func (t *Trackbot) enableMotors() {
	t.leftMotor.Enable()
	t.rightMotor.Enable()
}
func (t *Trackbot) signalLive() {	
	for i := 1; i <= 10; i++ {
		t.green("")
		time.Sleep(100 * time.Millisecond)
		t.blue("")
		time.Sleep(100 * time.Millisecond)
	}
	t.green("")
}
func (t *Trackbot) turnLeft(msg string) {
	t.leftMotor.Backward(t.speed)
	t.rightMotor.Forward(t.speed)
	t.enableMotors()
	t.out(msg)
}
func (t *Trackbot) turnRight(msg string) {
	t.leftMotor.Forward(t.speed)
	t.rightMotor.Backward(t.speed)
	t.enableMotors()
	t.out(msg)	
}
func (t *Trackbot) forward(msg string) {
	t.leftMotor.Forward(t.speed)
	t.rightMotor.Forward(t.speed)
	t.enableMotors()
	t.out(msg)
}
func (t *Trackbot) backward(msg string) {
	t.leftMotor.Backward(t.speed)
	t.rightMotor.Backward(t.speed)	
	t.enableMotors()
	t.out(msg)
}
func (t *Trackbot) red(msg string) {
	t.redPin.On()
	t.greenPin.Off()
	t.bluePin.Off()
	t.out(msg)
}
func (t *Trackbot) green(msg string) {
	t.redPin.Off()
	t.greenPin.On()
	t.bluePin.Off()
	t.out(msg)
}
func (t *Trackbot) blue(msg string) {
	t.redPin.Off()
	t.greenPin.Off()
	t.bluePin.On()
	t.out(msg)
}
func (t *Trackbot) ledOff(msg string) {
	t.redPin.Off()
	t.greenPin.Off()
	t.bluePin.Off()
	t.out(msg)
}
func (t *Trackbot) start() string {
	t.enableMotors()
	return fmt.Sprintf("Robot Going")
}
func (t *Trackbot) stop(msg string) {
	t.leftMotor.Disable()
	t.rightMotor.Disable()
	t.out(msg)
}
func (t *Trackbot) shutdown() string {
	return "Shutting down"
}
func (t *Trackbot) devices() []gobot.Device {
	return []gobot.Device{t.leftMotor, t.rightMotor, t.leftSensor, t.rightSensor,t.joystickDriver, t.redPin, t.bluePin, t.greenPin}
}
func (t *Trackbot) on(id string, f func(interface {})) {
	t.joystickDriver.On(id, f)
}
func boolToInt(b bool) byte {
	if b { return 1 }
	return 0
}
func (t *Trackbot) work() {
	
	t.signalLive()
	
	signalChan := make(chan bool)

	// create behaviour list
	shutdownAction := new(ShutdownAction)
	driveAction := new(DriveAction)	
	
	behaviours := 
		NewBehaviour("drive", signalChan, driveAction,
		NewBehaviour("shutdown", signalChan, shutdownAction,
		nil))
		
	// register behaviour events
	thisBehaviour := behaviours
	for thisBehaviour != nil {
		thisBehaviour.registerEvents(t)
		thisBehaviour = thisBehaviour.nextBehaviour
	}
	
	for {
		_ = <- signalChan
		//fmt.Println("State change is ready")
		
		var state State
		behaviours.getState(&state)
		//fmt.Println(state)
		
		if state.rgbSet {
			t.redPin.DigitalWrite(boolToInt(state.r))
			t.greenPin.DigitalWrite(boolToInt(state.g))
			t.bluePin.DigitalWrite(boolToInt(state.b))
		}
		if state.driveDirSet {
			if state.driveDir == TurnLeft {
				t.turnLeft("Turn Left")
			}
			if state.driveDir == TurnRight {
				t.turnRight("Turn Right")
			}
			if state.driveDir == Forward {
				t.forward("Forward")
			}
			if state.driveDir == Backward {
				t.backward("Backward")
			}
			if state.driveDir == Stop {
				t.stop("Stop")
			}
		}
		if state.speedSet {
			t.SetSpeed(state.speed, state.speedStr)
		}
		if state.shutdownSet {
			fmt.Println("Killing")
			time.Sleep(10000 * time.Millisecond)
			os.Exit(0)
		}
	}
	/*
	t.on("error", func(data interface{}) { 
		fmt.Println("error", data)
		time.Sleep(30000 * time.Millisecond)
	}) */
}
func main() {
	fmt.Println("Waiting to boot..")
	time.Sleep(20000 * time.Millisecond)
	fmt.Println("Booting...")
	
	gbot := gobot.NewMaster()
	server := api.NewAPI(gbot)
    server.Start()
	
	//joystickAdaptor := joystick.NewAdaptor("vnc-abspointer")
	joystickAdaptor := joystick.NewAdaptor("Magicsee R1")

	r := raspi.NewAdaptor()
	trackbot := NewTrackbot(r, joystickAdaptor)

	
	robot := gobot.NewRobot("tankbot", []gobot.Connection{r,joystickAdaptor}, trackbot.devices(), trackbot.work, )
	robot.AddCommand("Start", func(params map[string]interface{}) interface{} {	return trackbot.start() })
	robot.AddCommand("Stop", func(params map[string]interface{}) interface{} { return trackbot.shutdown() })

	//robot.Start()
	gbot.AddRobot(robot)
	gbot.Start()
}
