package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	renderer "github.com/thirst154/go-render/Renderer"
)

type Player struct {
	Position renderer.Vec3
	Rotation renderer.Vec3
	PosVel   renderer.Vec3
	Speed    float64

	// Mouse / input state
	lastMouseX    int
	lastMouseY    int
	mouseCaptured bool
	mouseInit     bool
	Sensitivity   float64
}

func NewPlayer(position renderer.Vec3, rotation renderer.Vec3) Player {
	return Player{
		Position:      position,
		Rotation:      rotation,
		PosVel:        renderer.NewVec3(0, 0, 0),
		Speed:         0.12,
		Sensitivity:   0.0025,
		mouseCaptured: false,
		mouseInit:     false,
	}
}

func (p *Player) MoveCamera(cam *renderer.Camera) {
	cam.Position = p.Position
	cam.Rotation = p.Rotation
}

func (p *Player) Update() {
	// Capture mouse once
	if !p.mouseCaptured {
		ebiten.SetCursorMode(ebiten.CursorModeCaptured)
		p.mouseCaptured = true
	}

	// read mouse delta
	mx, my := ebiten.CursorPosition()
	if !p.mouseInit {
		p.lastMouseX = mx
		p.lastMouseY = my
		p.mouseInit = true
	}

	dx := mx - p.lastMouseX
	dy := my - p.lastMouseY

	// update rotation: rotation.X = pitch, rotation.Y = yaw
	p.Rotation.X -= float64(dy) * p.Sensitivity // invert Y so moving mouse up looks up
	p.Rotation.Y += float64(dx) * p.Sensitivity

	// clamp pitch to avoid flipping
	limit := math.Pi/2 - 0.01
	if p.Rotation.X > limit {
		p.Rotation.X = limit
	}
	if p.Rotation.X < -limit {
		p.Rotation.X = -limit
	}

	// compute forward and right vectors from yaw/pitch
	forward := renderer.RotateVector(renderer.NewVec3(0, 0, 1), renderer.NewVec3(p.Rotation.X, p.Rotation.Y, 0))
	forward = renderer.Vec3Normalize(forward)

	right := renderer.RotateVector(renderer.NewVec3(1, 0, 0), renderer.NewVec3(0, p.Rotation.Y, 0))
	right = renderer.Vec3Normalize(right)

	// WASD + space/shift for vertical movement
	move := renderer.NewVec3(0, 0, 0)
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		move = renderer.Vec3Add(move, renderer.Vec3Scale(forward, p.Speed))
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		move = renderer.Vec3Add(move, renderer.Vec3Scale(forward, -p.Speed))
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		move = renderer.Vec3Add(move, renderer.Vec3Scale(right, -p.Speed))
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		move = renderer.Vec3Add(move, renderer.Vec3Scale(right, p.Speed))
	}
	// vertical
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		move = renderer.Vec3Add(move, renderer.NewVec3(0, p.Speed, 0))
	}
	if ebiten.IsKeyPressed(ebiten.KeyShift) || ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
		move = renderer.Vec3Add(move, renderer.NewVec3(0, -p.Speed, 0))
	}

	// apply movement
	p.Position = renderer.Vec3Add(p.Position, move)

	// store last mouse pos
	p.lastMouseX = mx
	p.lastMouseY = my
}
