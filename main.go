package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	renderer "github.com/thirst154/go-render/Renderer"
)

const (
	W = 320
	H = 240
)

type Game struct {
	renderer *renderer.Renderer
}

func (g *Game) Update() error {

	g.renderer.Render([]renderer.Sphere{})
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.WritePixels(g.renderer.Canvas.GetImage().Pix)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return W, H
}

func main() {
	renderer := renderer.NewRenderer(W, H)
	ebiten.RunGame(&Game{renderer: renderer})
}
