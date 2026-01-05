package renderer

import "image/color"

type Color struct {
	R, G, B, A uint8
}

var (
	BLACK            = Color{0, 0, 0, 255}
	WHITE            = Color{255, 255, 255, 255}
	RED              = Color{255, 0, 0, 255}
	GREEN            = Color{0, 255, 0, 255}
	BLUE             = Color{0, 0, 255, 255}
	BACKGROUND_COLOR = WHITE
)

func (c Color) RGBA() color.RGBA {
	return color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}
}
