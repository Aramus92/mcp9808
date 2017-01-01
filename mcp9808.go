package mcp9808

import (
	"log"

	"github.com/learnaddict/go-i2c"
)

const (
	hwBaseAddr  = 0x18
	hwBus       = 1 // Depends on Raspberry Pi version (0 or 1)
	hwTempAddr  = 0x05
	hwManuAddr  = 0x06
	hwDevIDAddr = 0x07
	chkManuID   = 0x54
	chkDevID    = 0x0400
)

// Find any modules and return the addresses of the detected
func Find() (m []uint8) {
	for i := uint8(0); i < 8; i++ {
		a := hwBaseAddr + i
		d := Check(a)
		if d == true {
			log.Printf("MCP9808 detected at 0x%x", a)
			m = append(m, a)
		}
	}
	return
}

// Check for a valid MCP9808 sensor at the i2c address
func Check(address uint8) bool {
	i, err := i2c.NewI2C(address, hwBus)
	if err != nil {
		return false
	}
	defer i.Close()

	if v, err := i.ReadRegU16BE(hwManuAddr); err != nil || v != chkManuID {
		return false
	}
	if v, err := i.ReadRegU16BE(hwDevIDAddr); err != nil || v != chkDevID {
		return false
	}
	return true
}

// Read the current temperature from the MCP9808 sensor at the i2c address
func Read(address uint8) (float32, error) {
	i, err := i2c.NewI2C(address, hwBus)
	if err != nil {
		return 0, err
	}
	defer i.Close()

	t, err := i.ReadRegU16BE(hwTempAddr)
	if err != nil {
		return 0, err
	}
	return float32(t&0x0FFF) / float32(16), nil
}
