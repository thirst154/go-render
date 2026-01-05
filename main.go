package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	W = 320
	H = 240
)

type Game struct {
	img *image.RGBA
}

func (g *Game) Update() error {
	// put pixel example
	g.img.Set(10, 10, color.RGBA{255, 0, 0, 255})
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.WritePixels(g.img.Pix)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return W, H
}

func main() {
	img := image.NewRGBA(image.Rect(0, 0, W, H))
	ebiten.RunGame(&Game{img: img})
}
