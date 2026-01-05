package renderer

import (
	"image"
	"image/color"
)

type Canvas struct {
	Width, Height int
	Pixles        *image.RGBA
}

func NewCanvas(width, height int) *Canvas {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	return &Canvas{
		Width:  width,
		Height: height,
		Pixles: img,
	}
}

func (c *Canvas) SetPixel(x, y int, col Color) {
	c.Pixles.SetRGBA(x, y, col.RGBA())
}

func (c *Canvas) GetPixel(x, y int) color.RGBA {
	return c.Pixles.RGBAAt(x, y)
}

func (c *Canvas) GetImage() *image.RGBA { return c.Pixles }
