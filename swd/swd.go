//-----------------------------------------------------------------------------
/*

SWD Device Functions

*/
//-----------------------------------------------------------------------------

package swd

import (
	cli "github.com/deadsy/go-cli"
)

//-----------------------------------------------------------------------------
// SWD driver interface

// State is the current SWD interface state
type State struct {
	TargetVoltage int  // Target reference voltage in mV
	Srst          bool // SRST pin state
}

// Driver is the interface for an SWD driver.
type Driver interface {
	GetState() (*State, error)
	Close() error
}

//-----------------------------------------------------------------------------

// Device stores the state for an SWD device.
type Device struct {
	drv Driver // swd driver
}

// GetDevice returns an SWD device.
func GetDevice(drv Driver) (*Device, error) {
	dev := &Device{
		drv: drv,
	}
	return dev, nil
}

//-----------------------------------------------------------------------------

// Menu submenu items
var Menu = cli.Menu{}

//-----------------------------------------------------------------------------
