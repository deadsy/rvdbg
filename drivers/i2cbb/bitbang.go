//-----------------------------------------------------------------------------
/*

I2C BitBang GPIO Driver

This implements the i2c.Driver interface.

*/
//-----------------------------------------------------------------------------

package i2cbb

import (
	"errors"
	"time"

	"github.com/deadsy/rvdbg/i2c"
)

//-----------------------------------------------------------------------------

// Driver is the lowlevel GPIO based api.
type Driver interface {
	SclRel()     // set high impedance of the SCL line
	SdaRel()     // set high impedance of the SDA line
	SclLo()      // drive SCL line low
	SdaLo()      // drive SDA line low
	SclRd() byte // read the SCL line
	SdaRd() byte // read the SDA line
	Init() error // initialise the IO lines
	Delay()      // delay to control the clock rate
}

//-----------------------------------------------------------------------------

type bitBang struct {
	drv Driver
}

// New returns an I2C driver using bit-banged GPIO.
func New(drv Driver) i2c.Driver {
	return &bitBang{
		drv: drv,
	}
}

// start creates the start condition- SDA goes low while SCL is high.
// On Exit- SDA and SCL are held low.
func (i2c *bitBang) start() error {
	// release the clock and data lines
	i2c.drv.SclRel()
	i2c.drv.SdaRel()
	// check that scl and sda are both high (no bus contention)
	i2c.drv.Delay()
	if i2c.drv.SdaRd() == 0 || i2c.drv.SclRd() == 0 {
		return errors.New("bus error")
	}
	i2c.drv.SdaLo()
	i2c.drv.Delay()
	i2c.drv.SclLo()
	i2c.drv.Delay()
	return nil
}

// stop creates the stop condition- SDA goes high while SCL is high.
// On Exit- SDA and SCL are released.
func (i2c *bitBang) stop() {
	i2c.drv.SclLo()
	i2c.drv.Delay()
	i2c.drv.SdaLo()
	i2c.drv.Delay()
	i2c.drv.SclRel()
	i2c.drv.Delay()
	i2c.drv.SdaRel()
	i2c.drv.Delay()
}

// clock SCL and read SDA at clock high.
// On Entry- SCL is held low.
// On Exit- SCL is held low, SDA =0/1 is returned.
func (i2c *bitBang) clock() (byte, error) {
	i2c.drv.Delay()
	i2c.drv.SclRel()
	// wait for any slave clock stretching
	delay := 100
	for i2c.drv.SclRd() == 0 && delay > 0 {
		time.Sleep(100 * time.Microsecond)
		delay--
	}
	if delay == 0 {
		i2c.stop()
		return 0, errors.New("slave timeout")
	}
	// read the data
	i2c.drv.Delay()
	val := i2c.drv.SdaRd()
	i2c.drv.SclLo()
	i2c.drv.Delay()
	return val, nil
}

// wrByte writes a byte of data to the slave.
// On Entry- SCL is held low.
// On Exit- SDA is released, SCL is held low.
func (i2c *bitBang) wrByte(val byte) error {
	mask := byte(0x80)
	for mask != 0 {
		if val&mask != 0 {
			i2c.drv.SdaRel()
		} else {
			i2c.drv.SdaLo()
		}
		_, err := i2c.clock()
		if err != nil {
			return err
		}
		mask >>= 1
	}
	i2c.drv.SdaRel()
	return nil
}

// rdByte reads a byte from a slave.
// On Entry- SCL is held low.
// On Exit- SDA is released, SCL is held low
func (i2c *bitBang) rdByte() (byte, error) {
	i2c.drv.SdaRel()
	val := byte(0)
	for i := 0; i < 8; i++ {
		val <<= 1
		bit, err := i2c.clock()
		if err != nil {
			return 0, err
		}
		val |= bit
	}
	return val, nil
}

// wrAck sends an ack to the slave.
func (i2c *bitBang) wrAck() error {
	i2c.drv.SdaLo()
	_, err := i2c.clock()
	if err != nil {
		return err
	}
	i2c.drv.SdaRel()
	return nil
}

// rdAck clocks in the SDA level from the slave.
// Return: false = no ack, true = ack
func (i2c *bitBang) rdAck() (bool, error) {
	i2c.drv.SdaRel()
	bit, err := i2c.clock()
	return bit == 0, err
}

// Read a buffer of bytes from device adr.
func (i2c *bitBang) Read(adr byte, buf []byte) (int, error) {
	// start a read cycle
	err := i2c.start()
	if err != nil {
		return 0, err
	}
	defer i2c.stop()
	// address the device
	i2c.wrByte(adr | 1)
	ack, err := i2c.rdAck()
	if err != nil {
		return 0, err
	}
	if !ack {
		return 0, errors.New("address error")
	}
	// read data
	for i := 0; i < len(buf); i++ {
		buf[i], err = i2c.rdByte()
		if err != nil {
			return i, err
		}
		// The last byte from the slave is not acked
		if i < len(buf)-1 {
			err := i2c.wrAck()
			if err != nil {
				return i, err
			}
		}
	}
	return len(buf), nil
}

// Write a buffer of bytes to device adr.
func (i2c *bitBang) Write(adr byte, buf []byte) (int, error) {
	// start a write cycle
	err := i2c.start()
	if err != nil {
		return 0, err
	}
	defer i2c.stop()
	// address the device
	i2c.wrByte(adr & ^byte(1))
	ack, err := i2c.rdAck()
	if err != nil {
		return 0, err
	}
	if !ack {
		return 0, errors.New("address error")
	}
	// write data
	for i := 0; i < len(buf); i++ {
		err := i2c.wrByte(buf[i])
		if err != nil {
			return i, err
		}
		ack, err := i2c.rdAck()
		if err != nil {
			return i, err
		}
		if !ack {
			// no ack from slave
			return i, errors.New("nak error")
		}
	}
	return len(buf), nil
}

// Init initialises the i2c bitbang driver.
func (i2c *bitBang) Init() error {
	return i2c.drv.Init()
}

//-----------------------------------------------------------------------------
