package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	hw "github.com/Sam36502/4RCH/src/hardware"
	"github.com/Sam36502/4RCH/src/util"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const ()

func main() {

	// Load Options from file
	err := util.LoadOptions()
	if err != nil {
		fmt.Printf("[ERROR]: Failed to load options file '%s'\nPlease make sure the options file is in the same directory as the Arch40 binary.\n", util.OPT_FILE)
		return
	}

	// Initialise Console window
	windowSize := 16 * util.GlobalOptions.PixelSize
	rl.InitWindow(windowSize, windowSize, "4RCH")
	if util.GlobalOptions.TargetFPS > 0 {
		rl.SetTargetFPS(util.GlobalOptions.TargetFPS)
	}
	rl.SetExitKey(util.GlobalOptions.Inputs[util.KC_POWER])

	// Seed RNG
	rand.Seed(time.Now().UnixMicro())
	rand.Int()

	var configMenu *util.ConfigMenu
	if util.GlobalOptions.ConfigKeys {
		configMenu = util.NewConfigMenu()
		rl.SetExitKey(0)
	}

	// Initialise Hardware
	screen := hw.NewScreen(
		util.GlobalOptions.ColourFG,
		util.GlobalOptions.ColourBG,
		util.GlobalOptions.PixelSize,
	)
	soundCard := hw.NewSoundCard(util.GlobalOptions.MasterVol, util.GlobalOptions.Samples)
	p1Controller := hw.NewController(1)
	p2Controller := hw.NewController(2)

	vm := hw.NewMachine()
	vm.PlugIn(screen)
	vm.PlugIn(soundCard)
	vm.PlugIn(p1Controller)
	vm.PlugIn(p2Controller)
	if len(util.GlobalOptions.MonAddrs) > 0 {
		monitor := hw.NewMemMonitor(util.GlobalOptions.MonAddrs)
		vm.PlugIn(monitor)
	}

	// Try to load cartridge provided as argument if available
	cartridgeFile := ""
	if len(os.Args) > 1 {
		cartridgeFile = os.Args[1]
		vm.LoadCartridgeFile(cartridgeFile)
	}

	// Register System key-listeners
	util.AddKeyListener(util.KC_RESET, vm.Reset)

	// Main loop
	var blink uint16 = 0
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		// Handle Config menu if active
		if configMenu != nil && !configMenu.IsDone() {
			configMenu.ConfigureInputs()
			if configMenu.IsDone() {
				configMenu = nil
			}
			rl.SetExitKey(util.GlobalOptions.Inputs[util.KC_POWER])
			continue
		}

		// Try to load a cartridge from dropped files
		if vm.Cart == nil {
			if blink%4096 < 2048 {
				screen.DrawBMP(hw.BMP_NO_CART)
			} else {
				screen.Clear()
			}
			blink++
		}
		if rl.IsFileDropped() {
			fs := rl.LoadDroppedFiles()
			cartridgeFile = fs[len(fs)-1]
			vm.LoadCartridgeFile(cartridgeFile)
		}

		if util.GlobalOptions.DebugMode && vm.Cart != nil {
			cmd := vm.Program[vm.InsPointer]
			rl.ClearBackground(util.GlobalOptions.ColourBG)
			vm.TickPeripherals()
			rl.DrawText(fmt.Sprintf("[%03d](A %02d): %v --> %v\n", vm.InsPointer, vm.Accumulator, cmd.Ins, cmd.Args), 10, 10, 20, rl.Pink)
			if rl.IsMouseButtonPressed(0) {
				vm.Tick()
			}
		} else {
			vm.Tick()
		}

		util.HandleInputs()
		rl.EndDrawing()
	}

	soundCard.Terminate()
	util.SaveOptions()
	rl.CloseWindow()
}
