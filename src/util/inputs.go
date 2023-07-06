package util

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Inputs map[string]int32

type KeyListener func()

type ConfigMenu struct {
	index int
	keys  []string
}

const (
	// System Key-Names
	KC_POWER = "kc_power"
	KC_RESET = "kc_reset"

	// Controller Key-Names
	KC_P1_UP    = "kc_p1_up"
	KC_P1_DOWN  = "kc_p1_down"
	KC_P1_LEFT  = "kc_p1_left"
	KC_P1_RIGHT = "kc_p1_right"
	KC_P1_A     = "kc_p1_a"
	KC_P1_B     = "kc_p1_b"
	KC_P1_X     = "kc_p1_x"
	KC_P1_Y     = "kc_p1_y"

	KC_P2_UP    = "kc_p2_up"
	KC_P2_DOWN  = "kc_p2_down"
	KC_P2_LEFT  = "kc_p2_left"
	KC_P2_RIGHT = "kc_p2_right"
	KC_P2_A     = "kc_p2_a"
	KC_P2_B     = "kc_p2_b"
	KC_P2_X     = "kc_p2_x"
	KC_P2_Y     = "kc_p2_y"
)

// Key descriptions for auto-configuration
var KEY_DESCS = map[string]string{
	KC_POWER:    "Power Button (Exits)",
	KC_RESET:    "Reset Button (Reboots)",
	KC_P1_UP:    "Player 1 Up Button",
	KC_P1_DOWN:  "Player 1 Down Button",
	KC_P1_LEFT:  "Player 1 Left Button",
	KC_P1_RIGHT: "Player 1 Right Button",
	KC_P1_A:     "Player 1 A Button",
	KC_P1_B:     "Player 1 B Button",
	KC_P1_X:     "Player 1 X Button",
	KC_P1_Y:     "Player 1 Y Button",
	KC_P2_UP:    "Player 2 Up Button",
	KC_P2_DOWN:  "Player 2 Down Button",
	KC_P2_LEFT:  "Player 2 Left Button",
	KC_P2_RIGHT: "Player 2 Right Button",
	KC_P2_A:     "Player 2 A Button",
	KC_P2_B:     "Player 2 B Button",
	KC_P2_X:     "Player 2 X Button",
	KC_P2_Y:     "Player 2 Y Button",
}

var KEY_LISTENERS = map[string][]KeyListener{}

func AddKeyListener(key string, lnr KeyListener) {
	KEY_LISTENERS[key] = append(KEY_LISTENERS[key], lnr)
}

func HandleInputs() {
	for key, lnrs := range KEY_LISTENERS {
		if kc, ex := GlobalOptions.Inputs[key]; ex && rl.IsKeyPressed(kc) {
			for _, lnr := range lnrs {
				lnr()
			}
		}
	}
}

func NewConfigMenu() *ConfigMenu {
	m := ConfigMenu{
		index: 0,
		keys: []string{
			"kc_power",
			"kc_reset",
			"kc_p1_up",
			"kc_p1_down",
			"kc_p1_left",
			"kc_p1_right",
			"kc_p1_a",
			"kc_p1_b",
			"kc_p1_x",
			"kc_p1_y",
			"kc_p2_up",
			"kc_p2_down",
			"kc_p2_left",
			"kc_p2_right",
			"kc_p2_a",
			"kc_p2_b",
			"kc_p2_x",
			"kc_p2_y",
		},
	}
	return &m
}

// Goes through all keys and prompts the use to press the desired key
func (m *ConfigMenu) ConfigureInputs() {
	key := m.keys[m.index]
	desc := KEY_DESCS[key]
	rl.ClearBackground(rl.RayWhite)
	PopupBox(fmt.Sprintf("Press the key for:\n%s", desc), rl.SkyBlue)

	pressed := rl.GetKeyPressed()
	if pressed != 0 {
		GlobalOptions.Inputs[key] = pressed
		m.index++
	}

	if m.IsDone() {
		GlobalOptions.ConfigKeys = false
		SaveOptions()
	}
}

func (m *ConfigMenu) IsDone() bool {
	return m.index == len(m.keys)
}
