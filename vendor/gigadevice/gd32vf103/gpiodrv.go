//-----------------------------------------------------------------------------
/*

GigaDevice gd32vf103 GPIO Driver

This code implements the gpio.Driver interface.

*/
//-----------------------------------------------------------------------------

package gd32vf103

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/deadsy/rvdbg/cpu/riscv/rv"
	"github.com/deadsy/rvdbg/soc"
)

//-----------------------------------------------------------------------------

// contains returns whether a string list contains a string.
func contains(list []string, s string) bool {
	for i := range list {
		if s == list[i] {
			return true
		}
	}
	return false
}

//-----------------------------------------------------------------------------

type GpioConfig struct {
	Pin  string // port/bit name (Pxy) where x = A,B,C..., y = 0..31
	Mode string // pin mode (i,o)
	Name string // target pin name
}

type GpioDriver struct {
	cfg   []GpioConfig
	dbg   rv.Debug
	dev   *soc.Device
	ports []string // available port names
}

func NewGpioDriver(dbg rv.Debug, dev *soc.Device, cfg []GpioConfig) *GpioDriver {
	return &GpioDriver{
		cfg:   cfg,
		dbg:   dbg,
		dev:   dev,
		ports: []string{"GPIOA", "GPIOB", "GPIOC", "GPIOD", "GPIOE"},
	}
}

// Init initialises the GPIO sub-system
func (drv *GpioDriver) Init() error {
	return nil
}

// Status returns a status string for GPIOs
func (drv *GpioDriver) Status() string {
	return ""
}

// Pin converts a pin name to a port/bit tuple
func (drv *GpioDriver) Pin(name string) (string, uint, error) {
	name = strings.ToUpper(name)
	if !strings.HasPrefix(name, "P") {
		return "", 0, errors.New("pin name must start with \"P\"")
	}
	if len(name) < 3 || len(name) > 4 {
		return "", 0, fmt.Errorf("pin name \"%s\" has the wrong length", name)
	}
	port := fmt.Sprintf("GPIO%s", name[1:2])
	if !contains(drv.ports, port) {
		return "", 0, fmt.Errorf("no port \"%s\" on this device", port)
	}
	bit, err := strconv.ParseUint(name[2:], 10, 8)
	if err != nil {
		return "", 0, fmt.Errorf("could not parse gpio bit in \"%s\"", name)
	}
	if bit > 15 {
		return "", 0, errors.New("gpio bit is > 15")
	}
	return "", uint(bit), nil
}

// Set an output bit
func (drv *GpioDriver) Set(port string, bit uint) error {
	r, err := drv.dev.GetPeripheralRegister(port, "BOP")
	if err != nil {
		return err
	}
	return r.Wr(1 << bit)
}

// Clear an output bit
func (drv *GpioDriver) Clr(port string, bit uint) error {
	r, err := drv.dev.GetPeripheralRegister(port, "BC")
	if err != nil {
		return err
	}
	return r.Wr(1 << bit)
}

//-----------------------------------------------------------------------------
