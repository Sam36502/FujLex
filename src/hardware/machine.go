package hardware

import (
	"fmt"
	"math/rand"
)

type Nybble uint8

// The actual 4RCH virtual machine
type Machine struct {
	Accumulator  Nybble
	RAM          [16][16]Nybble
	RAMListeners map[byte][]RAMListener
	Program      []Command
	Cart         *Cart
	InsPointer   byte
	peripherals  []Peripheral
	IsRunning    bool
	ticks        int8
}

func NewMachine() *Machine {
	vm := Machine{
		Accumulator:  0,
		RAM:          [16][16]Nybble{},
		RAMListeners: map[byte][]RAMListener{},
		Program:      []Command{},
		Cart:         nil,
		InsPointer:   0,
		IsRunning:    true,
		ticks:        0,
	}
	return &vm
}

func (vm *Machine) PlugIn(prph Peripheral) {
	vm.peripherals = append(vm.peripherals, prph)
	vm.AddRAMListener(prph.GetListener(vm))
}

// Tries to load a cartridge from a file
func (vm *Machine) LoadCartridgeFile(filename string) {
	c, err := LoadCartFromFile(filename)
	if err != nil {
		fmt.Printf("[ERROR]: Failed to load cartridge '%s': %v\n", filename, err)
	} else {
		vm.Cart = c
	}

	vm.LoadCartridge(c)
}

func (vm *Machine) LoadCartridge(c *Cart) {
	vm.Program = []Command{}
	for i := 0; i < len(c.Program); i++ {
		ins := ALL_INS[c.Program[i]]
		cmd := Command{
			Ins:  ins,
			Args: c.Program[i+1 : i+1+ins.Nargs],
		}
		i += ins.Nargs
		vm.Program = append(vm.Program, cmd)
	}
	vm.PlugIn(c)
}

// Tells the VM to perform one action
func (vm *Machine) Tick() {
	if vm.Cart == nil {
		return
	}

	vm.TickPeripherals()
	vm.RAM[PERIPHERAL_PAGE][FPG_RAND] = Nybble(rand.Intn(0xF))
	cmd := vm.Program[vm.InsPointer]
	vm.ExecuteCommand(cmd)
}

func (vm *Machine) TickPeripherals() {
	for _, prph := range vm.peripherals {
		prph.Tick(vm)
	}
}

func (vm *Machine) Reset() {
	vm.IsRunning = false
	for _, prph := range vm.peripherals {
		prph.Reset()
	}
	vm.Accumulator = 0
	vm.RAM = [16][16]Nybble{}
	vm.InsPointer = 0
	vm.ticks = 0
	vm.IsRunning = true
}
