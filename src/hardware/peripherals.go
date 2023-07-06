package hardware

import (
	"fmt"
	"strconv"

	"github.com/Sam36502/4RCH/src/util"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Peripheral interface {
	Tick(vm *Machine)
	GetListener(vm *Machine) ([]byte, RAMListener)
	Reset()
}

type RAMListener func(val Nybble)

type MemMonitor struct {
	addrs  []byte
	vals   []Nybble
	colour []rl.Color
}

var MON_CLR = rl.SkyBlue

const (
	PERIPHERAL_PAGE = 0xF
)

// F-page Addresses
const (
	FPG_P1_DPAD = 0x0 // Player 1 Direction-Pad
	FPG_P1_BTNS = 0x1 // Player 1 Buttons
	FPG_P2_DPAD = 0x2 // Player 2 Direction-Pad
	FPG_P2_BTNS = 0x3 // Player 2 Buttons

	FPG_SCR_X   = 0x4 // Screen X Coord
	FPG_SCR_Y   = 0x5 // Screen Y Coord
	FPG_SCR_VAL = 0x6 // Screen Pixel Value
	FPG_SCR_OPT = 0x7 // Screen Options

	FPG_SND_VOL = 0x8 // Sound Card Volume
	FPG_SND_PTC = 0x9 // Sound Card Pitch
	FPG_SND_OPT = 0xA // Sound Card reserved

	FPG_RAND = 0xB // Pseudo-Random Number

	FPG_DSK_H   = 0xC // High-nyble of disk address   \
	FPG_DSK_M   = 0xD // Middle-nyble of disk address  } 12-bit Address
	FPG_DSK_L   = 0xE // Low-nyble of disk address    /
	FPG_DSK_VAL = 0xF // Value of the selected disk nyble
)

func (vm *Machine) AddRAMListener(trigAddrs []byte, lnr RAMListener) {
	for _, addr := range trigAddrs {
		vm.RAMListeners[addr] = append(vm.RAMListeners[addr], lnr)
	}
}

// Creates a memory monitor from a list of hex addresses
func NewMemMonitor(hexAddrs []string) *MemMonitor {
	mon := MemMonitor{
		addrs:  []byte{},
		vals:   []Nybble{},
		colour: []rl.Color{},
	}
	for _, hex := range hexAddrs {
		mii, err := strconv.ParseUint(hex, 16, 32)
		if err != nil || mii > 255 {
			fmt.Printf("[WARNING]: Invalid monitor address provided '%s', Must be single hex byte.", hex)
			continue
		}
		mon.addrs = append(mon.addrs, byte(mii))
		mon.vals = append(mon.vals, 0)
		mon.colour = append(mon.colour, MON_CLR)
	}
	return &mon
}

func (mon *MemMonitor) Tick(vm *Machine) {
	yPos := int32(10)
	if util.GlobalOptions.DebugMode {
		yPos += 20
	}
	rl.DrawText("Monitored Addresses:", 10, yPos, 20, MON_CLR)
	for i, addr := range mon.addrs {
		yPos += 20
		val := mon.vals[i]
		rl.DrawText(
			fmt.Sprintf("0x%02X: %02d, 0x%X, 0b%04b", addr, val, val, val),
			10, yPos, 20, mon.colour[i],
		)
	}
	for i, col := range mon.colour {
		mon.colour[i] = util.StepTowardsColour(col, MON_CLR)
	}
}

func (mon *MemMonitor) GetListener(vm *Machine) ([]byte, RAMListener) {
	return mon.addrs, func(val Nybble) {
		for i, addr := range mon.addrs {
			val := vm.RAM[addr>>4][addr%16]
			if mon.vals[i] != val {
				mon.colour[i] = rl.Red
				mon.vals[i] = val
			}
		}
	}
}

func (mon *MemMonitor) Reset() {
	mon.vals = []Nybble{}
	mon.colour = []rl.Color{}
}
