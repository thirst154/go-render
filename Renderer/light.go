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
func ComputeLighting(P Vec3, N Vec3, V Vec3, s float64, lights []Light) float64 {
	i := 0.0
	N = vec3Normalize(N) // ensure normal is unit
	V = vec3Normalize(V)

	for _, light := range lights {
		if light.Type == Ambient {
			i += light.Intensity
			continue
		}

		var L Vec3
		if light.Type == Point {
			L = vec3Subtract(light.Position, P)
		} else { // Directional
			L = light.Direction
		}

		// normalize light direction
		L = vec3Normalize(L)

		// Diffuse
		nDotL := vec3Dot(N, L)
		if nDotL > 0 {
			i += light.Intensity * nDotL
		}

		// Specular R = 2 * N * dot(N, L) - L
		if s != -1 && nDotL > 0 {
			R := vec3Subtract(vec3Scale(N, 2*nDotL), L)
			rDotV := vec3Dot(vec3Normalize(R), V)
			if rDotV > 0 {
				i += light.Intensity * math.Pow(rDotV, s)
			}
		}

	}

	return i
}
