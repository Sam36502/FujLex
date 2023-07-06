package util

import (
	"encoding/json"
	"io/ioutil"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	OPT_FILE      = "options.json"
	OPT_FILE_MODE = 0650
	OPT_INDENT    = "    "
)

type Options struct {
	MasterVol  float32   `json:"master_vol"`
	PixelSize  int32     `json:"pixel_size"`
	TargetFPS  int32     `json:"target_fps"`
	ColourFG   rl.Color  `json:"colour_fg"`
	ColourBG   rl.Color  `json:"colour_bg"`
	Samples    [4]string `json:"samples"`
	Inputs     Inputs    `json:"inputs"`
	ConfigKeys bool      `json:"configure_key"`
	MonAddrs   []string  `json:"mon_addrs"`
	DebugMode  bool      `json:"debug_mode"`
}

var GlobalOptions Options

func LoadOptions() error {
	data, err := ioutil.ReadFile(OPT_FILE)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &GlobalOptions)
}

func SaveOptions() error {
	data, err := json.MarshalIndent(GlobalOptions, "", OPT_INDENT)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(OPT_FILE, data, OPT_FILE_MODE)
}
