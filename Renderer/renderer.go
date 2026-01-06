package renderer

import (
	"math"
)

type Camera struct {
	Position Vec3 // O
	Rotation Vec3
}

type Viewport struct {
	Width, Height int // Vw, Vh
	Distance      int // d

}

type Renderer struct {
	Cam      Camera
	Viewport Viewport
	Canvas   *Canvas
}

type Sphere struct {
	Center   Vec3
	Radius   int
	Color    Color
	Specular float64
}

type Scene struct {
	Spheres []Sphere
	Lights  []Light
}

// Cx, Cy are positions on the canvas. The 2d Pixels on the screen.

func NewRenderer(screenW, screenH int) *Renderer {

	canvas := NewCanvas(screenW, screenH)

	cam := Camera{
		Position: NewVec3(0, 0, 0),
		Rotation: NewVec3(0, 0, 0),
	}

	viewport := Viewport{
		Width:    1,
		Height:   1,
		Distance: 1,
	}

	return &Renderer{
		Cam:      cam,
		Canvas:   canvas,
		Viewport: viewport,
	}
}

func CanvasToViewport(x int, y int, canvas *Canvas, viewport *Viewport) Vec3 {
	return NewVec3(
		float64(x)*float64(viewport.Width)/float64(canvas.Width),
		float64(y)*float64(viewport.Height)/float64(canvas.Height),
		float64(viewport.Distance),
	)
}

// Rotate a vector by Euler angles (rotation.X = pitch, rotation.Y = yaw, rotation.Z = roll)
// Rotation angles are in radians.
func RotateVector(v Vec3, rotation Vec3) Vec3 {
	// Rotate around X (pitch)
	sinX := math.Sin(rotation.X)
	cosX := math.Cos(rotation.X)
	y1 := v.Y*cosX - v.Z*sinX
	z1 := v.Y*sinX + v.Z*cosX

	// Rotate around Y (yaw)
	sinY := math.Sin(rotation.Y)
	cosY := math.Cos(rotation.Y)
	x2 := v.X*cosY + z1*sinY
	z2 := -v.X*sinY + z1*cosY

	// Rotate around Z (roll)
	sinZ := math.Sin(rotation.Z)
	cosZ := math.Cos(rotation.Z)
	x3 := x2*cosZ - y1*sinZ
	y3 := x2*sinZ + y1*cosZ

	return NewVec3(x3, y3, z2)
}

// T is any real number from -infinity to infinity
func RayEquation(CamPos Vec3, ViewportPos Vec3, t float64) Vec3 {

	rayDir := RayDirection(CamPos, ViewportPos)

	return vec3Add(CamPos, vec3Scale(rayDir, t))
}

func RayDirection(CamPos Vec3, ViewportPos Vec3) Vec3 {
	return vec3Subtract(ViewportPos, CamPos)
}

// Points on a sphere satisfy the equation (P - C) · (P - C) = r^2

func SphereEquation(P Vec3, C Vec3, r int) bool {
	d := vec3Subtract(P, C)

	return d.X*d.X+d.Y*d.Y+d.Z*d.Z <= float64(r*r)
}

// Points on an Ray satisfy the equation P = O + tD

func IntersectRaySphere(O Vec3, D Vec3, sphere Sphere) (float64, float64) {
	r := sphere.Radius
	CO := vec3Subtract(O, sphere.Center)

	a := vec3Dot(D, D)
	b := 2 * vec3Dot(CO, D)
	c := vec3Dot(CO, CO) - float64(r*r)

	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return math.Inf(1), math.Inf(1)
	}

	t1 := (-b + math.Sqrt(discriminant)) / (2 * a)
	t2 := (-b - math.Sqrt(discriminant)) / (2 * a)

	return t1, t2
}

// simple util to make TraceRay function cleaner
func in(l, min, max float64) bool {
	return l >= min && l <= max
}

func ClosestIntersection(O Vec3, D Vec3, tMin, tMax float64, scene Scene) (*Sphere, float64) {
	closestT := math.Inf(1)
	var closestSphere *Sphere = nil

	// For each sphere in the scene
	// Check for intersection with the ray
	// If there is an intersection, check if it is the closest one
	// If it is, update closestT and closestSphere
	for _, sphere := range scene.Spheres {
		t1, t2 := IntersectRaySphere(O, D, sphere)
		if in(t1, tMin, tMax) && t1 < closestT {
			closestT = t1
			closestSphere = &sphere
		}

		if in(t2, tMin, tMax) && t2 < closestT {
			closestT = t2
			closestSphere = &sphere
		}

	}

	return closestSphere, closestT
}

func TraceRay(O Vec3, D Vec3, tMin, tMax float64, scene Scene) Color {

	closestSphere, closestT := ClosestIntersection(O, D, tMin, tMax, scene)

	if closestSphere == nil {
		return BACKGROUND_COLOR // Background color
	}

	//return closestSphere.Color

	// Lighting
	P := vec3Add(O, vec3Scale(D, closestT))
	N := vec3Subtract(P, closestSphere.Center) // Sphere normal
	N = vec3Normalize(N)

	DNegated := NewVec3(-D.X, -D.Y, -D.Z)
	illum := ComputeLighting(P, N, DNegated, closestSphere.Specular, scene)
	if illum < 0 {
		illum = 0
	}
	if illum > 1 {
		illum = 1
	}

	color := closestSphere.Color
	// scale each channel by the illumination (preserve hue) and clamp to 0..255
	r := math.Min(255, math.Max(0, float64(color.R)*illum))
	g := math.Min(255, math.Max(0, float64(color.G)*illum))
	b := math.Min(255, math.Max(0, float64(color.B)*illum))

	color.R = uint8(r)
	color.G = uint8(g)
	color.B = uint8(b)

	return color

}

func (r *Renderer) Render(scene Scene) {

	// Canvas.SetPixel expects 0-based top-left coordinates.
	// Iterate over image coordinates (0..Width-1, 0..Height-1)
	// and map them to centered canvas coordinates for the ray math.
	halfW := r.Canvas.Width / 2
	halfH := r.Canvas.Height / 2
	for x := 0; x < r.Canvas.Width; x++ {
		for y := 0; y < r.Canvas.Height; y++ {
			// convert to centered coordinates expected by CanvasToViewport
			cx := x - halfW
			cy := halfH - y
			D := CanvasToViewport(cx, cy, r.Canvas, &r.Viewport)
			// apply camera rotation to the viewport ray direction
			D = RotateVector(D, r.Cam.Rotation)
			color := TraceRay(r.Cam.Position, D, 1, 10000, scene)
			r.Canvas.SetPixel(x, y, color)
		}
	}
}
