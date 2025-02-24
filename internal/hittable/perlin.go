package hittable

import (
	"math"
	"math/rand"

	"github.com/nsp5488/go_raytracer/internal/vec"
)

const pointCount = 256

type Perlin struct {
	randVec [pointCount]*vec.Vec3
	permX   [pointCount]int
	permY   [pointCount]int
	permZ   [pointCount]int
}

// Generates a new Perlin noise texture
func NewPerlin() *Perlin {
	p := &Perlin{}
	for i := 0; i < pointCount; i++ {
		p.randVec[i] = vec.RangeRandom(-1, 1).UnitVector()
	}
	p.generatePerm()
	return p
}

// Returns the value of this randomized perlin noise at the given point
func (p *Perlin) Noise(point *vec.Vec3) float64 {
	u := point.X() - math.Floor(point.X())
	v := point.Y() - math.Floor(point.Y())
	w := point.Z() - math.Floor(point.Z())

	i := int(math.Floor(point.X()))
	j := int(math.Floor(point.Y()))
	k := int(math.Floor(point.Z()))
	c := [2][2][2]*vec.Vec3{}
	for di := 0; di < 2; di++ {
		for dj := 0; dj < 2; dj++ {
			for dk := 0; dk < 2; dk++ {
				c[di][dj][dk] = p.randVec[p.permX[(i+di)&255]^
					p.permY[(j+dj)&255]^
					p.permZ[(k+dk)&255]]
			}
		}
	}

	return perlinInterpolation(&c, u, v, w)
}

// Sums repeated calls to Noise to generate a turbulent texture
func (p *Perlin) Turbulence(point *vec.Vec3, depth int) float64 {
	accum := 0.0
	temp_p := &vec.Vec3{}
	*temp_p = *point
	weight := 1.0

	for i := 0; i < depth; i++ {
		accum += weight * p.Noise(temp_p)
		weight *= 0.5
		temp_p.ScaleInplace(2)
	}
	return math.Abs(accum)
}

// Generates data then permutes it
func (p *Perlin) generatePerm() {
	for i := 0; i < pointCount; i++ {
		p.permX[i] = i
		p.permY[i] = i
		p.permZ[i] = i
	}

	permute(&p.permX)
	permute(&p.permY)
	permute(&p.permZ)
}

// Helper method for generatePerm, shuffles the values in the provided array randomly
func permute(p *[pointCount]int) {
	for i := len(p) - 1; i > 0; i-- {
		target := rand.Intn(i)
		p[i], p[target] = p[target], p[i]
	}
}

// Calculates the perlin interpolation of the provided floats
func perlinInterpolation(c *[2][2][2]*vec.Vec3, u, v, w float64) float64 {
	// voodoo magic AKA Hermitian smoothing
	uu := u * u * (3 - 2*u)
	vv := v * v * (3 - 2*v)
	ww := w * w * (3 - 2*w)

	accumulator := 0.0
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				weight := vec.New(u-float64(i), v-float64(j), w-float64(k))
				accumulator += ((float64(i)*uu + float64(1-i)*(1-uu)) *
					(float64(j)*vv + float64((1-j))*(1-vv)) *
					(float64(k)*ww + float64((1-k))*(1-ww)) * c[i][j][k].Dot(weight))
			}
		}
	}
	return accumulator
}
