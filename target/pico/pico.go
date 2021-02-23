//-----------------------------------------------------------------------------
/*

Raspberry Pi Pico is a microcontroller board using a RP2040 ARM Cortex-M0+

See: https://www.raspberrypi.org/products/raspberry-pi-pico/

*/
//-----------------------------------------------------------------------------

package pico

import (
	"errors"
	"fmt"
	"os"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/cpu/arm/cm"
	"github.com/deadsy/rvdbg/flash"
	"github.com/deadsy/rvdbg/gpio"
	"github.com/deadsy/rvdbg/itf"
	"github.com/deadsy/rvdbg/mem"
	"github.com/deadsy/rvdbg/soc"
	"github.com/deadsy/rvdbg/target"
	"github.com/deadsy/rvdbg/vendor/rpi/rp20xx"
)

//-----------------------------------------------------------------------------

// Info is target information.
var Info = target.Info{
	Name:     "pico",
	Descr:    "RPi Pico Board (RP2040 Dual Core Cortex-M0+)",
	DbgMode:  itf.ModeSwd,
	DbgSpeed: 4000,
	Volts:    3300,
}

//-----------------------------------------------------------------------------

// menuRoot is the root menu.
var menuRoot = cli.Menu{
	{"cpu", cm.Menu, "cpu functions"},
	{"exit", target.CmdExit},
	{"flash", flash.Menu, "flash functions"},
	{"gpio", gpio.Menu, "gpio functions"},
	{"help", target.CmdHelp},
	{"history", target.CmdHistory, cli.HistoryHelp},
	{"swd", swd.Menu, "jtag functions"},
	{"map", soc.CmdMap},
	{"mem", mem.Menu, "memory functions"},
	{"regs", soc.CmdRegs, soc.RegsHelp},
}

//-----------------------------------------------------------------------------
// GPIO names

var gpioNames = map[string]string{
	"GPIO0":  "GP0",
	"GPIO1":  "GP1",
	"GPIO2":  "GP2",
	"GPIO3":  "GP3",
	"GPIO4":  "GP4",
	"GPIO5":  "GP5",
	"GPIO6":  "GP6",
	"GPIO7":  "GP7",
	"GPIO8":  "GP8",
	"GPIO9":  "GP9",
	"GPIO10": "GP10",
	"GPIO11": "GP11",
	"GPIO12": "GP12",
	"GPIO13": "GP13",
	"GPIO14": "GP14",
	"GPIO15": "GP15",
	"GPIO16": "GP16",
	"GPIO17": "GP17",
	"GPIO18": "GP18",
	"GPIO19": "GP19",
	"GPIO20": "GP20",
	"GPIO21": "GP21",
	"GPIO22": "GP22",
	//"GPIO23": "",
	//"GPIO24": "",
	//"GPIO25": "",
	"GPIO26": "GP26_A0",
	"GPIO27": "GP27_A1",
	"GPIO28": "GP28_A2",
	//"GPIO29": "",
}

//-----------------------------------------------------------------------------

// Target is the application structure for the target.
type Target struct {
	swdDevice   *swd.Device
	cmDebug     cm.Debug
	socDevice   *soc.Device
	socDriver   *socDriver
	memDriver   *memDriver
	gpioDriver  *rp20xx.GpioDriver
	flashDriver *rp20xx.FlashDriver
}

// New returns a new pico target.
func New(swdDriver swd.Driver) (target.Target, error) {

	// get the SWD state
	state, err := swdDriver.GetState()
	if err != nil {
		return nil, err
	}

	// check the voltage
	if state.TargetVoltage >= 0 {
		if float32(state.TargetVoltage) < 0.9*float32(Info.Volts) {
			return nil, fmt.Errorf("target voltage is too low (%dmV), is the target connected and powered?", state.TargetVoltage)
		}
	}

	// check the ~SRST state
	if !state.Srst {
		return nil, errors.New("target ~SRST line asserted, target is held in reset")
	}

	cmDebug, err := cm.NewDebug(swdDevice)
	if err != nil {
		return nil, err
	}

	// create the SoC device
	socDevice := rp20xx.NewSoC(rp20xx.RP2040).Setup()
	socDriver := newSocDriver(cmDebug)

	// gpio driver
	gpioDriver, err := rp20xx.NewGpioDriver(socDriver, socDevice, gpioNames)
	if err != nil {
		return nil, err
	}

	// flash driver
	flashDriver, err := rp20xx.NewFlashDriver(socDriver, socDevice)
	if err != nil {
		return nil, err
	}

	return &Target{
		swdDevice:   swdDevice,
		cmDebug:     cmDebug,
		socDevice:   socDevice,
		socDriver:   socDriver,
		memDriver:   newMemDriver(cmDebug, socDevice),
		gpioDriver:  gpioDriver,
		flashDriver: flashDriver,
	}, nil
}

//-----------------------------------------------------------------------------

// GetPrompt returns the target prompt string.
func (t *Target) GetPrompt() string {
	return t.cmDebug.GetPrompt(Info.Name)
}

// GetMenuRoot returns the target root menu.
func (t *Target) GetMenuRoot() []cli.MenuItem {
	return menuRoot
}

// Shutdown shuts down the target application.
func (t *Target) Shutdown() {
}

// Put outputs a string to the user application.
func (t *Target) Put(s string) {
	os.Stdout.WriteString(s)
}

//-----------------------------------------------------------------------------

// GetMemoryDriver returns a memory driver for this target.
func (t *Target) GetMemoryDriver() mem.Driver {
	return t.memDriver
}

// GetGpioDriver returns a GPIO driver for this target.
func (t *Target) GetGpioDriver() gpio.Driver {
	return t.gpioDriver
}

// GetFlashDriver returns a Flash driver for this target.s
func (t *Target) GetFlashDriver() flash.Driver {
	return t.flashDriver
}

// GetCmDebug returns a Cortex-M debug driver for this target.
func (t *Target) GetCmDebug() cm.Debug {
	return t.cmDebug
}

// GetSoC returns the SoC device and driver.
func (t *Target) GetSoC() (*soc.Device, soc.Driver) {
	return t.socDevice, t.socDriver
}

// GetSwdDevice returns the SWD device.
func (t *Target) GetSwdDevice() *swd.Device {
	return t.swdDevice
}

//-----------------------------------------------------------------------------
