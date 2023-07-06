package hardware

import (
	"fmt"
	"strings"

	"github.com/Sam36502/4RCH/src/util"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Controller struct {
	PlayerNr   int
	Keys       []string
	DPadNybble Nybble
	BtnsNybble Nybble
}

const PNUM_CHAR = "#"

// Key indices
const (
	CKI_UP    = 0
	CKI_DOWN  = 1
	CKI_LEFT  = 2
	CKI_RIGHT = 3
	CKI_A     = 4
	CKI_B     = 5
	CKI_X     = 6
	CKI_Y     = 7
)

var controllerKeys = []string{
	"kc_p#_up",
	"kc_p#_down",
	"kc_p#_left",
	"kc_p#_right",
	"kc_p#_a",
	"kc_p#_b",
	"kc_p#_x",
	"kc_p#_y",
}

func NewController(playerNr int) *Controller {
	c := Controller{
		PlayerNr:   playerNr,
		Keys:       []string{},
		DPadNybble: 0b0000,
		BtnsNybble: 0b0000,
	}
	for _, key := range controllerKeys {
		c.Keys = append(c.Keys, strings.ReplaceAll(key, PNUM_CHAR, fmt.Sprint(playerNr)))
	}
	return &c
}

func (c *Controller) Tick(vm *Machine) {
	c.DPadNybble = 0 |
		c.ButtonStatus(CKI_DOWN)<<3 |
		c.ButtonStatus(CKI_UP)<<2 |
		c.ButtonStatus(CKI_RIGHT)<<1 |
		c.ButtonStatus(CKI_LEFT)<<0
	c.BtnsNybble = 0 |
		c.ButtonStatus(CKI_Y)<<3 |
		c.ButtonStatus(CKI_X)<<2 |
		c.ButtonStatus(CKI_B)<<1 |
		c.ButtonStatus(CKI_A)<<0

	addr := FPG_P1_DPAD + (c.PlayerNr-1)*2
	vm.RAM[PERIPHERAL_PAGE][addr+0] = c.DPadNybble
	vm.RAM[PERIPHERAL_PAGE][addr+1] = c.BtnsNybble
}

// Returns the status (0/1) of whether a given button is pressed
func (c *Controller) ButtonStatus(cki int) Nybble {
	keycode := util.GlobalOptions.Inputs[c.Keys[cki]]
	if rl.IsKeyDown(keycode) {
		return 1
	}
	return 0
}

func (c *Controller) GetListener(vm *Machine) ([]byte, RAMListener) {
	return []byte{}, func(val Nybble) {}
}

func (c *Controller) Reset() {}
