//-----------------------------------------------------------------------------
/*

RP20xx GPIO Driver

This code implements the gpio.Driver interface.

*/
//-----------------------------------------------------------------------------

package rp20xx

import (
	"errors"

	"github.com/deadsy/rvdbg/soc"
)

//-----------------------------------------------------------------------------

// GpioDriver is a GPIO driver for the rp2040.
type GpioDriver struct {
	names map[string]string // standard pin names to target pin names
	cache map[string]string // cache of pin mode/value strings
	drv   soc.Driver
	dev   *soc.Device
}

// NewGpioDriver retuns a new GPIO driver for the rp2040.
func NewGpioDriver(drv soc.Driver, dev *soc.Device, names map[string]string) (*GpioDriver, error) {
	return &GpioDriver{
		names: names,
		cache: make(map[string]string),
		drv:   drv,
		dev:   dev,
	}, nil
}

// Status returns a status string for GPIOs
func (drv *GpioDriver) Status() string {
	return ""
}

// Pin converts a pin name to a port/bit tuple
func (drv *GpioDriver) Pin(name string) (string, uint, error) {
	return "", 0, errors.New("TODO")
}

// Set sets an output bit
func (drv *GpioDriver) Set(port string, bit uint) error {
	return errors.New("TODO")
}

// Clr clears an output bit
func (drv *GpioDriver) Clr(port string, bit uint) error {
	return errors.New("TODO")
}

//-----------------------------------------------------------------------------
