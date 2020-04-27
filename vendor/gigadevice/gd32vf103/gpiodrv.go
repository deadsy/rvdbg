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

	"github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/soc"
	"github.com/deadsy/rvdbg/util"
)

//-----------------------------------------------------------------------------

// GpioDriver is a GPIO driver for the gd32vf103.
type GpioDriver struct {
	names map[string]string // standard pin names to target pin names
	cache map[string]string // cache of pin mode/value strings
	drv   soc.Driver
	dev   *soc.Device
}

// NewGpioDriver retuns a new GPIO driver for the gd32vf103.
func NewGpioDriver(drv soc.Driver, dev *soc.Device, names map[string]string) *GpioDriver {
	return &GpioDriver{
		names: names,
		cache: make(map[string]string),
		drv:   drv,
		dev:   dev,
	}
}

// targetName returns a taget specific name for the GPIO pin name.
func (drv *GpioDriver) targetName(name string) string {
	if s, ok := drv.names[name]; ok {
		return s
	}
	return ""
}

// changed returns true if the pin mode/value has changed.
func (drv *GpioDriver) changed(name, mode string) bool {
	rc := false
	if m, ok := drv.cache[name]; ok {
		rc = m != mode
	}
	drv.cache[name] = mode
	return rc
}

// Status returns a status string for GPIOs
func (drv *GpioDriver) Status() string {

	s := [][]string{}

	// look for ports GPIOA..GPIOH
	for i := 0; i < 8; i++ {
		port := fmt.Sprintf("GPIO%c", 'A'+i)

		if drv.dev.GetPeripheral(port) == nil {
			continue
		}

		ctl0, err := drv.dev.RdPeripheralRegister(drv.drv, port, "CTL0")
		if err != nil {
			s = append(s, []string{"", "", "", fmt.Sprintf("could not read %s.CTL0: %s", port, err)})
			continue
		}

		ctl1, err := drv.dev.RdPeripheralRegister(drv.drv, port, "CTL1")
		if err != nil {
			s = append(s, []string{"", "", "", fmt.Sprintf("could not read %s.CTL1: %s", port, err)})
			continue
		}

		istat, err := drv.dev.RdPeripheralRegister(drv.drv, port, "ISTAT")
		if err != nil {
			s = append(s, []string{"", "", "", fmt.Sprintf("could not read %s.ISTAT: %s", port, err)})
			continue
		}

		octl, err := drv.dev.RdPeripheralRegister(drv.drv, port, "OCTL")
		if err != nil {
			s = append(s, []string{"", "", "", fmt.Sprintf("could not read %s.OCTL: %s", port, err)})
			continue
		}

		lock, err := drv.dev.RdPeripheralRegister(drv.drv, port, "LOCK")
		if err != nil {
			s = append(s, []string{"", "", "", fmt.Sprintf("could not read %s.LOCK: %s", port, err)})
			continue
		}

		prefix := fmt.Sprintf("P%c", 'A'+i)

		var md, ctl uint
		for j := 0; j < 16; j++ {
			if j < 8 {
				// ctl0
				md = (ctl0 >> (j * 4)) & 3
				ctl = (ctl0 >> ((j * 4) + 2)) & 3

			} else {
				// ctl1
				md = (ctl1 >> (j * 4)) & 3
				ctl = (ctl1 >> ((j * 4) + 2)) & 3
			}

			mode := ""
			cfg := []string{}

			switch md {
			case 0:
				mode = "in"
			case 1, 2, 3:
				mode = "out"
			}

			if mode == "in" {
				switch ctl {
				case 0:
					cfg = append(cfg, "analog")
				case 1:
					cfg = append(cfg, "float")
				case 2:
					cfg = append(cfg, []string{"pull-down", "pull-up"}[util.BoolToInt(octl&(1<<j) != 0)])
				}
			} else {
				switch ctl {
				case 0:
					cfg = append(cfg, "push-pull")
				case 1:
					cfg = append(cfg, "open-drain")
				case 2:
					mode = "af"
					cfg = append(cfg, "push-pull")
				case 3:
					mode = "af"
					cfg = append(cfg, "open-drain")
				}
				switch md {
				case 1:
					cfg = append(cfg, "10 MHz")
				case 2:
					cfg = append(cfg, "2 MHz")
				case 3:
					cfg = append(cfg, "50 MHz")
				}
			}
			if lock&(1<<j) != 0 {
				cfg = append(cfg, "locked")
			}

			if mode == "in" {
				mode += []string{"(0)", "(1)"}[util.BoolToInt(istat&(1<<j) != 0)]
			}

			if mode == "out" {
				mode += []string{"(0)", "(1)"}[util.BoolToInt(octl&(1<<j) != 0)]
			}

			name := fmt.Sprintf("%s%d", prefix, j)
			mode += []string{"", "*"}[util.BoolToInt(drv.changed(name, mode))]
			tgtName := drv.targetName(name)
			cfgStr := fmt.Sprintf("(%s)", strings.Join(cfg, ","))
			s = append(s, []string{name, mode, tgtName, cfgStr})
		}
	}

	return cli.TableString(s, []int{0, 0, 0, 0}, 1)
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
	if drv.dev.GetPeripheral(port) == nil {
		return "", 0, fmt.Errorf("no port \"%s\" on this device", port)
	}
	bit, err := strconv.ParseUint(name[2:], 10, 8)
	if err != nil {
		return "", 0, fmt.Errorf("could not parse gpio bit in \"%s\"", name)
	}
	if bit > 15 {
		return "", 0, errors.New("gpio bit is > 15")
	}
	return port, uint(bit), nil
}

// Set sets an output bit
func (drv *GpioDriver) Set(port string, bit uint) error {
	return drv.dev.WrPeripheralRegister(drv.drv, port, "BOP", 1<<bit)
}

// Clr clears an output bit
func (drv *GpioDriver) Clr(port string, bit uint) error {
	return drv.dev.WrPeripheralRegister(drv.drv, port, "BC", 1<<bit)
}

//-----------------------------------------------------------------------------
