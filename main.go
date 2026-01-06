package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	renderer "github.com/thirst154/go-render/Renderer"
)

const (
	W = 320
	H = 240
)

type Game struct {
	render *renderer.Renderer
	scene  *renderer.Scene
}

func (g *Game) Update() error {
	// Example Spheres

	// Keypress handling using ebiten
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.render.Cam.Position.X -= 0.1
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.render.Cam.Position.X += 0.1
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.render.Cam.Position.Z += 0.1
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.render.Cam.Position.Z -= 0.1
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.render.Cam.Rotation.X += 0.1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.render.Cam.Rotation.X -= 0.1
	}

	g.render.Render(*g.scene)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.WritePixels(g.render.Canvas.GetImage().Pix)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return W, H
}

func randomSphere(number int) []renderer.Sphere {
	spheres := []renderer.Sphere{}

	for i := 0; i < number; i++ {
		spheres = append(spheres, renderer.Sphere{
			Center: renderer.NewVec3(float64(rand.Intn(W)), float64(rand.Intn(H)), float64(rand.Intn(H))),
			Radius: rand.Intn(100),
			Color:  renderer.RandomColor(),
		})
	}

	return spheres

}

func main() {
	spheres := []renderer.Sphere{
		{Center: renderer.NewVec3(0, -1, 3), Radius: 2, Color: renderer.RED, Specular: 500},  // Red sphere SHINY
		{Center: renderer.NewVec3(2, 0, 4), Radius: 1, Color: renderer.BLUE, Specular: 500},  // Blue sphere
		{Center: renderer.NewVec3(-2, 0, 4), Radius: 1, Color: renderer.GREEN, Specular: 10}, // Green sphere
		{Center: renderer.NewVec3(0, -5001, 0), Radius: 5000, Color: renderer.YELLOW, Specular: 1000},
	}

	lights := []renderer.Light{
		renderer.NewAmbientLight(0.1),
		renderer.NewDirectionalLight(0.6, renderer.NewVec3(0, 10, -1)),
	}

	render := renderer.NewRenderer(W, H)
	ebiten.RunGame(&Game{render: render, scene: &renderer.Scene{Spheres: spheres, Lights: lights}})
}
