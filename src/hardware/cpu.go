package hardware

import (
	"fmt"
	"math"
	"time"
)

type Instruction struct {
	Opcode string
	Binary Nybble
	Nargs  int
}

type Command struct {
	Ins  Instruction
	Args []Nybble
}

func (vm *Machine) ExecuteCommand(cmd Command) {
	if !vm.IsRunning {
		return
	}

	// Check correct no. of args were passed
	if cmd.Ins.Nargs != len(cmd.Args) {
		fmt.Printf("[ERROR]: Passed an invalid no. of arguments")
		return
	}

	// Dereference memory references ahead of time for convenience
	memVal := Nybble(0)
	if len(cmd.Args) == 2 {
		memVal = vm.RAM[cmd.Args[1]][cmd.Args[0]]
	}

	nextInsPointer := vm.InsPointer + 1
	switch cmd.Ins.Binary {

	case 0b0000:
		vm.IsRunning = false

	case 0b0001:
		vm.Accumulator = cmd.Args[0]

	case 0b0010:
		vm.Accumulator = memVal

	case 0b0011:
		vm.RAM[cmd.Args[1]][cmd.Args[0]] = vm.Accumulator

		// Call RAM-Listeners if any
		if lnrs, ex := vm.RAMListeners[byte(cmd.Args[1]<<4)|byte(cmd.Args[0])]; ex {
			for _, lnr := range lnrs {
				lnr(vm.Accumulator)
			}
		}

	case 0b0100:
		vm.Accumulator += cmd.Args[0]
		vm.Accumulator %= 16
		vm.Accumulator -= cmd.Args[1]
		vm.Accumulator %= 16

	case 0b0101:
		vm.Accumulator += memVal
		vm.Accumulator %= 16

	case 0b0110:
		fallthrough
	case 0b0111:
		fmt.Printf("[WARNING]: IP:%d; Unassigned instruction called\n", vm.InsPointer)

	case 0b1000:
		vm.Accumulator = ^vm.Accumulator

	case 0b1001:
		vm.Accumulator |= memVal

	case 0b1010:
		vm.Accumulator &= memVal

	case 0b1011:
		if (cmd.Args[0]>>3)%2 == 0 {
			if (cmd.Args[0]>>2)%2 == 0 {
				vm.Accumulator <<= cmd.Args[1] % 4
				vm.Accumulator %= 0xF
			} else {
				vm.Accumulator >>= cmd.Args[1] % 4
			}
		} else {
			if (cmd.Args[0]>>2)%2 == 0 {
				vm.Accumulator <<= cmd.Args[1] % 4
				vm.Accumulator |= (vm.Accumulator >> 4) % 0xF
			} else {
				vm.Accumulator <<= 4
				vm.Accumulator >>= cmd.Args[1] % 4
				vm.Accumulator |= (vm.Accumulator >> 4) % 0xF
			}
		}

	case 0b1100:
		scale := (cmd.Args[0] >> 1) % 8
		mul := math.Pow10(int(scale) - 4)
		length := (cmd.Args[0]%2)<<4 | cmd.Args[1]
		dur := time.Duration(mul * float64(length) * float64(time.Second.Nanoseconds()))

		vm.IsRunning = false
		time.AfterFunc(dur, func() {
			vm.IsRunning = true
		})

	case 0b1101:
		if vm.Accumulator != cmd.Args[0] {
			nextInsPointer += byte(cmd.Args[1])
		}

	case 0b1110:
		nextInsPointer = byte(cmd.Args[1]<<4) | byte(cmd.Args[0])

	case 0b1111:
		cm := vm.RAM[0][cmd.Args[0]]
		pg := vm.RAM[0][cmd.Args[1]]
		nextInsPointer = byte(pg<<4) | byte(cm)

	}

	vm.InsPointer = nextInsPointer
}

func (i Instruction) String() string {
	return fmt.Sprintf("%s (0b%04b, 0x%X) A%d", i.Opcode, i.Binary, i.Binary, i.Nargs)
}
