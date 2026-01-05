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
	Center Vec3
	Radius int
	Color  Color
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

func in(l, min, max float64) bool {
	return l >= min && l <= max
}

func TraceRay(O Vec3, D Vec3, tMin, tMax float64, spheres []Sphere) Color {
	closestT := math.Inf(1)
	var closestSphere *Sphere = nil

	// For each sphere in the scene
	// Check for intersection with the ray
	// If there is an intersection, check if it is the closest one
	// If it is, update closestT and closestSphere
	for _, sphere := range spheres {
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
	if closestSphere == nil {
		return BACKGROUND_COLOR // Background color
	}

	return closestSphere.Color

}

func (r *Renderer) Render(spheres []Sphere) {

	// Example Spheres
	spheres = []Sphere{
		{Center: NewVec3(0, -1, 3), Radius: 1, Color: RED},   // Red sphere
		{Center: NewVec3(2, 0, 4), Radius: 1, Color: BLUE},   // Blue sphere
		{Center: NewVec3(-2, 0, 4), Radius: 1, Color: GREEN}, // Green sphere
	}

	// Canvas.SetPixel expects 0-based top-left coordinates.
	// Iterate over image coordinates (0..Width-1, 0..Height-1)
	// and map them to centered canvas coordinates for the ray math.
	halfW := r.Canvas.Width / 2
	halfH := r.Canvas.Height / 2
	for x := 0; x < r.Canvas.Width; x++ {
		for y := 0; y < r.Canvas.Height; y++ {
			// convert to centered coordinates expected by CanvasToViewport
			cx := x - halfW
			cy := y - halfH
			D := CanvasToViewport(cx, cy, r.Canvas, &r.Viewport)
			color := TraceRay(r.Cam.Position, D, 1, 10000, spheres)
			r.Canvas.SetPixel(x, y, color)
		}
	}
}
