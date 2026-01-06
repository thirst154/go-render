package renderer

import "math"

type LightType string

const (
	Ambient     LightType = "ambient"
	Directional LightType = "directional"
	Point       LightType = "point"
)

type Light struct {
	Type      LightType
	Intensity float64
	Position  Vec3
	Direction Vec3
}

func NewAmbientLight(intensity float64) Light {
	return Light{
		Type:      Ambient,
		Intensity: intensity,
	}
}

func NewDirectionalLight(intensity float64, direction Vec3) Light {
	return Light{
		Type:      Directional,
		Intensity: intensity,
		Direction: direction,
	}
}

func NewPointLight(intensity float64, position Vec3) Light {
	return Light{
		Type:      Point,
		Intensity: intensity,
		Position:  position,
	}
}

// P is the position of the point being lit
// N is the normal of the point being lit
// V is the view vector
// s is the specular vector
func ComputeLighting(P Vec3, N Vec3, V Vec3, s float64, scene Scene) float64 {
	i := 0.0
	N = Vec3Normalize(N) // ensure normal is unit
	V = Vec3Normalize(V)

	for _, light := range scene.Lights {
		if light.Type == Ambient {
			i += light.Intensity
			continue
		}

		var L Vec3
		var tMax float64
		if light.Type == Point {
			L = Vec3Subtract(light.Position, P)
			tMax = 1
		} else { // Directional
			L = light.Direction
			tMax = math.Inf(1)
		}

		// normalize light direction
		L = Vec3Normalize(L)

		// Shadow check
		shadow_sphere, _ := ClosestIntersection(P, L, 0.001, tMax, scene)
		if shadow_sphere != nil {
			continue
		}

		// Diffuse
		nDotL := Vec3Dot(N, L)
		if nDotL > 0 {
			i += light.Intensity * nDotL
		}

		// Specular R = 2 * N * dot(N, L) - L
		if s != -1 && nDotL > 0 {
			R := Vec3Subtract(Vec3Scale(N, 2*nDotL), L)
			rDotV := Vec3Dot(Vec3Normalize(R), V)
			if rDotV > 0 {
				i += light.Intensity * math.Pow(rDotV, s)
			}
		}

	}

	return i
}
