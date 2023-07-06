package util

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func StepTowardsColour(col rl.Color, trgt rl.Color) rl.Color {
	new := col
	new.R += DiffSign(col.R, trgt.R)
	new.G += DiffSign(col.G, trgt.G)
	new.B += DiffSign(col.B, trgt.B)
	return new
}

func DiffSign(a, b uint8) uint8 {
	n := int(b) - int(a)
	if n > 0 {
		return 1
	}
	if n < 0 {
		return 255
	}
	return 0
}

func PopupBox(msg string, clr rl.Color) {
	darkClr := clr
	darkClr.R -= 50
	darkClr.G -= 50
	darkClr.B -= 50

	width := 300
	height := 150
	x := (int(GlobalOptions.PixelSize) * 16 / 2) - width/2
	y := (int(GlobalOptions.PixelSize) * 16 / 2) - height/2
	rec := rl.Rectangle{X: float32(x), Y: float32(y), Width: float32(width), Height: float32(height)}
	rl.DrawRectangleRec(rec, clr)
	rl.DrawRectangleLinesEx(rec, 5, darkClr)

	rl.DrawText(msg, int32(x+25), int32(y+25), 20, darkClr)
	rl.EndDrawing()
}
