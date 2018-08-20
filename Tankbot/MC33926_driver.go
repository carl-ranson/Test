package Tankbot

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

// MotorDriver Represents a Motor
type MC33926Driver struct {
	name             string
	connection       gpio.DigitalWriter
	connection2		 gpio.PwmWriter
	EnablePin		 string
	DirectionPin     string
	PwmPin           string
	CurrentState     byte
	CurrentDirection string
	isOn             bool
}

// NewMotorDriver return a new MotorDriver given a DigitalWriter and pin
func NewMC33926Driver(a gpio.DigitalWriter, b gpio.PwmWriter, enablePin, directionPin, pwmPin string) *MC33926Driver {
	return &MC33926Driver{
		name:             gobot.DefaultName("Motor"),
		connection:       a,
		connection2:	  b,
		EnablePin:        enablePin,
		DirectionPin:     directionPin,
		PwmPin:           pwmPin,
		CurrentState:     0,
		CurrentDirection: "forward",
		isOn:             false,
	}
}

// Name returns the MotorDrivers name
func (m *MC33926Driver) Name() string { return m.name }

// SetName sets the MotorDrivers name
func (m *MC33926Driver) SetName(n string) { m.name = n }

// Connection returns the MotorDrivers Connection
func (m *MC33926Driver) Connection() gobot.Connection { return m.connection.(gobot.Connection) }

// Start implements the Driver interface
func (m *MC33926Driver) Start() (err error) { return }

// Halt implements the Driver interface
func (m *MC33926Driver) Halt() (err error) { return }


// Forward sets the forward pin to the specified speed
func (m *MC33926Driver) Forward(speed byte) (err error) {
	err = m.connection.DigitalWrite(m.DirectionPin, 1)
	if err != nil {
		return err
	}
	return m.connection2.PwmWrite(m.PwmPin, speed)
}

// Backward sets the backward pin to the specified speed
func (m *MC33926Driver) Backward(speed byte) (err error) {
	err = m.connection.DigitalWrite(m.DirectionPin, 0)
	if err != nil {
		return err
	}
	return m.connection2.PwmWrite(m.PwmPin, speed)
}

// Forward sets the forward pin to the specified speed
func (m *MC33926Driver) Enable() (err error) {
	err = m.connection.DigitalWrite(m.EnablePin, 1)
	return
}

// Backward sets the backward pin to the specified speed
func (m *MC33926Driver) Disable() (err error) {
	err = m.connection.DigitalWrite(m.EnablePin, 0)
	return
}

