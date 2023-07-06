package hardware

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Screen struct {
	vram  [16]uint16
	fg    rl.Color
	bg    rl.Color
	scale int32
}

func NewScreen(fg, bg rl.Color, scale int32) *Screen {
	return &Screen{
		vram:  [16]uint16{},
		fg:    fg,
		bg:    bg,
		scale: scale,
	}
}

func (s *Screen) Clear() {
	s.vram = [16]uint16{}
	rl.ClearBackground(s.bg)
}

func (s *Screen) Invert() {
	for i, row := range s.vram {
		s.vram[i] = ^row
	}
}

func (s *Screen) DrawBMP(bmp [16]uint16) {
	rl.ClearBackground(s.bg)
	var x, y int32
	for y = 0; y < 16; y++ {
		for x = 0; x < 16; x++ {
			if (bmp[y]>>(15-x))%2 == 1 {
				ps := s.scale
				rl.DrawRectangle(
					x*ps, y*ps,
					ps, ps,
					s.fg,
				)
			}
		}
	}
}

func (s *Screen) DrawVRAM() {
	s.DrawBMP(s.vram)
}

func (s *Screen) Tick(vm *Machine) {
	s.DrawVRAM()

	// Check if top opt bits are used
	opt := vm.RAM[PERIPHERAL_PAGE][FPG_SCR_OPT]
	if (opt>>3)%2 == 1 {
		s.Clear()
		opt &= 0b0111
	}
	if (opt>>2)%2 == 1 {
		s.Invert()
		opt &= 0b1011
	}
	vm.RAM[PERIPHERAL_PAGE][FPG_SCR_OPT] = opt

	// Update state of the screen_value address
	x := vm.RAM[PERIPHERAL_PAGE][FPG_SCR_X]
	y := vm.RAM[PERIPHERAL_PAGE][FPG_SCR_Y]
	var val Nybble = 0
	switch opt % 4 {

	case 0b00:
		val = Nybble(s.vram[y] >> (8 - 1 - x) % 2)

	case 0b01:
		val = Nybble(s.vram[y] >> (8 - 4 - x) % 16)

	case 0b10:
		for i := Nybble(0); i < 4; i++ {
			val |= Nybble(s.vram[y+i]>>(8-1-x)%2) << (4 - i)
		}

	case 0b11:
		val |= Nybble(s.vram[y+0]>>(8-2-x)%4) << 2
		val |= Nybble(s.vram[y+1]>>(8-2-x)%4) << 0
	}
	vm.RAM[PERIPHERAL_PAGE][FPG_SCR_VAL] = val

}

func (s *Screen) GetListener(vm *Machine) ([]byte, RAMListener) {
	return []byte{(PERIPHERAL_PAGE << 4) | FPG_SCR_VAL},
		func(val Nybble) {
			x := vm.RAM[PERIPHERAL_PAGE][FPG_SCR_X]
			y := vm.RAM[PERIPHERAL_PAGE][FPG_SCR_Y]
			opt := vm.RAM[PERIPHERAL_PAGE][FPG_SCR_OPT]
			switch opt % 4 {

			case 0b00:
				s.vram[y] ^= uint16(val) << (16 - 1 - x)

			case 0b01:
				s.vram[y] ^= uint16(val) << (16 - 4 - x)

			case 0b10:
				for i := Nybble(0); i < 4; i++ {
					s.vram[y+i] ^= uint16(val) << (16 - 1 - x) << (4 - i)
				}

			case 0b11:
				s.vram[y+0] ^= uint16(val) << (16 - 2 - x) << 2
				s.vram[y+1] ^= uint16(val) << (16 - 2 - x) << 0
			}
		}
}

func (s *Screen) Reset() {
	s.vram = [16]uint16{}
	s.Clear()
}
