package hardware

import (
	"io/ioutil"
)

const (
	SIZE_PROG = 256
	SIZE_DATA = 4096
	FILE_MODE = 0650
)

type Cart struct {
	Program    [SIZE_PROG]Nybble
	Data       [SIZE_DATA]Nybble
	IsWritable bool
}

func LoadCartFromFile(filename string) (*Cart, error) {
	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if len(fileData) == 0 {
		return nil, nil
	}

	cart := Cart{}

	// Handle flags
	flagByte := fileData[0]
	if (flagByte>>0)%2 == 1 {
		cart.IsWritable = true
	}

	fileData = fileData[1:]
	for bi, b := range fileData {
		for hn := 0; hn < 2; hn++ {
			ni := (bi << 1) | hn // Calculate nybble index
			n := Nybble(b>>((1-hn)*4)) % 16

			if ni < SIZE_PROG {
				cart.Program[ni] = n
			} else {
				cart.Data[ni-SIZE_PROG] = n
			}
		}
	}

	return &cart, nil
}

func SaveCartToFile(filename string, cart Cart) error {

	flagByte := byte(0)
	if cart.IsWritable {
		flagByte |= (1 << 0)
	}

	fileData := make([]byte, (SIZE_PROG/2)+1)
	fileData[0] = flagByte

	// Write Program
	for ni := 0; ni < SIZE_PROG; ni++ {
		bi := (ni >> 1) + 1
		hn := ni % 2
		b := byte(cart.Program[ni])
		fileData[bi] |= b << ((1 - hn) * 4)
	}

	// Find end of data block
	size := SIZE_DATA
	for i := SIZE_DATA - 1; i >= 0; i-- {
		if cart.Data[i] != 0 {
			size = i + 1
			break
		}
		if i == 0 {
			size = 0
		}
	}

	// Write Data
	for ni := 0; ni < size; ni++ {
		if ni%2 == 0 {
			fileData = append(fileData, byte(cart.Data[ni]))
		} else {
			bi := 1 + SIZE_PROG + ni/2
			fileData[bi] |= byte(cart.Data[ni]) << 4
		}
	}

	return ioutil.WriteFile(filename, fileData, FILE_MODE)
}

func (c *Cart) Tick(vm *Machine) {
	addr := 0 |
		uint16(vm.RAM[PERIPHERAL_PAGE][FPG_DSK_H])<<8 |
		uint16(vm.RAM[PERIPHERAL_PAGE][FPG_DSK_M])<<4 |
		uint16(vm.RAM[PERIPHERAL_PAGE][FPG_DSK_L])

	vm.RAM[PERIPHERAL_PAGE][FPG_DSK_VAL] = vm.Cart.Data[addr]
}

func (c *Cart) GetListener(vm *Machine) ([]byte, RAMListener) {
	return []byte{(PERIPHERAL_PAGE << 4) | FPG_DSK_VAL}, func(val Nybble) {
		if !vm.Cart.IsWritable {
			c.Tick(vm)
			return
		}

		addr := uint16(0) |
			uint16(vm.RAM[PERIPHERAL_PAGE][FPG_DSK_H])<<8 |
			uint16(vm.RAM[PERIPHERAL_PAGE][FPG_DSK_M])<<4 |
			uint16(vm.RAM[PERIPHERAL_PAGE][FPG_DSK_L])

		vm.Cart.Data[addr] = val
	}
}

func (c *Cart) Reset() {}
