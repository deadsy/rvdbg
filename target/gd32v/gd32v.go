//-----------------------------------------------------------------------------
/*

GD32V is a development platform using a GD32VF103VBT6 RISC-V RV32.

See: https://www.seeedstudio.com/SeeedStudio-GD32-RISC-V-Dev-Board-p-4302.html

*/
//-----------------------------------------------------------------------------

package gd32v

import (
	"errors"
	"fmt"
	"gigadevice/gd32vf103"
	"os"

	cli "github.com/deadsy/go-cli"
	"github.com/deadsy/rvdbg/cpu/riscv"
	"github.com/deadsy/rvdbg/cpu/riscv/rv"
	"github.com/deadsy/rvdbg/cpu/riscv/rv13"
	"github.com/deadsy/rvdbg/flash"
	"github.com/deadsy/rvdbg/gpio"
	"github.com/deadsy/rvdbg/i2c"
	"github.com/deadsy/rvdbg/itf"
	"github.com/deadsy/rvdbg/jtag"
	"github.com/deadsy/rvdbg/mem"
	"github.com/deadsy/rvdbg/soc"
	"github.com/deadsy/rvdbg/target"
)

//-----------------------------------------------------------------------------

// Info is target information.
var Info = target.Info{
	Name:     "gd32v",
	Descr:    "GD32V Board (GigaDevice GD32VF103VBT6 RISC-V RV32)",
	DbgMode:  itf.ModeJtag,
	DbgSpeed: 4000,
	Volts:    3300,
}

//-----------------------------------------------------------------------------

// menuRoot is the root menu.
var menuRoot = cli.Menu{
	{"cpu", riscv.Menu, "cpu functions"},
	{"csr", riscv.CmdCSR, riscv.CsrHelp},
	{"da", riscv.CmdDisassemble, riscv.DisassembleHelp},
	{"dbg", rv13.Menu, "debugger functions"},
	{"exit", target.CmdExit},
	{"flash", flash.Menu, "flash functions"},
	{"gpio", gpio.Menu, "gpio functions"},
	{"gpr", riscv.CmdGpr},
	{"halt", riscv.CmdHalt},
	{"hart", riscv.CmdHart, riscv.HartHelp},
	{"help", target.CmdHelp},
	{"history", target.CmdHistory, cli.HistoryHelp},
	{"i2c", i2c.Menu, "i2c functions"},
	{"jtag", jtag.Menu, "jtag functions"},
	{"map", soc.CmdMap},
	{"mem", mem.Menu, "memory functions"},
	{"regs", soc.CmdRegs, soc.RegsHelp},
	{"resume", riscv.CmdResume},
}

//-----------------------------------------------------------------------------
// GPIO names

var gpioNames = map[string]string{
	// switches
	"PA0":  "WKUP",
	"PC13": "TAMPER",
	// leds
	"PB0": "LED_G", // 0 == on, blue LED on my board
	"PB1": "LED_B", // 0 == on
	"PB5": "LED_R", // 0 == on, a blue LED on my board
	// LCD-FSMC-8080 mode
	"PE1":  "LCD_RST",
	"PD10": "FSMC_D15",
	"PD9":  "FSMC_D14",
	"PD8":  "FSMC_D13",
	"PE15": "FSMC_D12",
	"PE14": "FSMC_D11",
	"PE13": "FSMC_D10",
	"PE12": "FSMC_D09",
	"PE11": "FSMC_D08",
	"PE10": "FSMC_D07",
	"PE9":  "FSMC_D06",
	"PE8":  "FSMC_D05",
	"PE7":  "FSMC_D04",
	"PD1":  "FSMC_D03",
	"PD0":  "FSMC_D02",
	"PD15": "FSMC_D01",
	"PD14": "FSMC_D00",
	"PD4":  "FSMC_NOE",
	"PD5":  "FSMC_NWE",
	"PD11": "FSMC_A16",
	"PD7":  "FSMC_NE1",
	"PE0":  "CPT_IO23",
	"PD13": "CPT_IO24",
	"PE2":  "CPT_IO25",
	"PE3":  "CPT_IO26",
	"PE4":  "CPT_IO27",
	"PD12": "LCD_BL",
	// sd card
	"PB12": "TF_CS",
	"PB15": "SPI1_MOSI",
	"PB13": "SPI1_SCK",
	"PB14": "SPI1_MISO",
	// spi flash
	"PC0": "FLASH_CS",
	"PA5": "SPI0_SCK",
	"PA7": "SPI0_MOSI",
	"PA6": "SPI0_MISO",
	// i2c
	"PB6": "I2C0_SCL",
	"PB7": "I2C0_SDA",
	// jtag
	"PB4":  "TRST",
	"PA15": "TDI",
	"PA13": "TMS",
	"PA14": "TCK",
	"PB3":  "TDO",
	// cn2
	"PA9":  "USART0_TX",
	"PA10": "USART0_RX",
	//"PD11": "EXMC_A16",
	//"PD13": "EXMC_A18", ??
	//"PD14": "EXMC_D0",
	//"PD0":  "EXMC_D2",
	//"PE7":  "EXMC_D4",
	//"PE9":  "EXMC_D6",
	//"PE11": "EXMC_D8",
	//"PE13": "EXMC_D10",
	//"PE15": "EXMC_D12",
	//"PD9":  "EXMC_D14",
	"PB2": "BOOT1",
	"PA2": "USART1_TX",
	"PA3": "USART1_RX",
	//"PD12": "EXMC_A17", ??
	//"PD5":  "EXMC_NWE",
	//"PD15": "EXMC_D1",
	//"PD1":  "EXMC_D3",
	//"PE8":  "EXMC_D5",
	//"PE10": "EXMC_D7",
	//"PE12": "EXMC_D9",
	//"PE14": "EXMC_D11",
	//"PD8":  "EXMC_D13",
	//"PD10": "EXMC_D15",
	//"PE1":  "EXMC_NBL1", ??
	// cn3
	"PC8":  "TIMER2_CH2",
	"PC9":  "TIMER2_CH3",
	"PC10": "UART3_TX",
	"PC11": "UART3_RX",
	"PD2":  "TIMER2_ETI",
	"PC12": "UART4_TX",
	// usb
	"PA11": "USBFS_DM", // USART0_CTS, CAN0_RX, USBFS_DM, TIMER0_CH3
	"PA12": "USBFS_DP", // USART0_RTS, USBFS_DP, CAN0_TX, TIMER0_ETI
	"PD6":  "usb",      // EXMC_NWAIT, USART1_RX
}

//-----------------------------------------------------------------------------

// Target is the application structure for the target.
type Target struct {
	jtagDevice  *jtag.Device
	rvDebug     rv.Debug
	socDevice   *soc.Device
	socDriver   *socDriver
	memDriver   *memDriver
	csrDriver   *csrDriver
	gpioDriver  *gd32vf103.GpioDriver
	flashDriver *gd32vf103.FlashDriver
}

// New returns a new gd32v target.
func New(jtagDriver jtag.Driver) (target.Target, error) {

	// get the JTAG state
	state, err := jtagDriver.GetState()
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

	// make the jtag chain
	jtagChain, err := jtag.NewChain(jtagDriver, gd32vf103.Chain)
	if err != nil {
		return nil, err
	}

	// make the jtag device for the cpu core
	jtagDevice, err := jtagChain.GetDevice(gd32vf103.CoreIndex)
	if err != nil {
		return nil, err
	}

	rvDebug, err := riscv.NewDebug(jtagDevice)
	if err != nil {
		return nil, err
	}

	// create the SoC device
	socDevice := gd32vf103.NewSoC(gd32vf103.VB).Setup()
	socDriver := newSocDriver(rvDebug)

	// gpio driver
	gpioDriver, err := gd32vf103.NewGpioDriver(socDriver, socDevice, gpioNames)
	if err != nil {
		return nil, err
	}

	// flash driver
	flashDriver, err := gd32vf103.NewFlashDriver(socDriver, socDevice)
	if err != nil {
		return nil, err
	}

	return &Target{
		jtagDevice:  jtagDevice,
		rvDebug:     rvDebug,
		socDevice:   socDevice,
		socDriver:   socDriver,
		memDriver:   newMemDriver(rvDebug, socDevice),
		csrDriver:   newCsrDriver(rvDebug),
		gpioDriver:  gpioDriver,
		flashDriver: flashDriver,
	}, nil
}

//-----------------------------------------------------------------------------

// GetPrompt returns the target prompt string.
func (t *Target) GetPrompt() string {
	return t.rvDebug.GetPrompt(Info.Name)
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

// GetRiscvDebug returns a RISC-V debug driver for this target.
func (t *Target) GetRiscvDebug() rv.Debug {
	return t.rvDebug
}

// GetSoC returns the SoC device and driver.
func (t *Target) GetSoC() (*soc.Device, soc.Driver) {
	return t.socDevice, t.socDriver
}

// GetCSR returns the CSR device and driver.
func (t *Target) GetCSR() (*soc.Device, soc.Driver) {
	return t.rvDebug.GetCurrentHart().CSR, t.csrDriver
}

// GetJtagDevice returns the JTAG device.
func (t *Target) GetJtagDevice() *jtag.Device {
	return t.jtagDevice
}

//-----------------------------------------------------------------------------
