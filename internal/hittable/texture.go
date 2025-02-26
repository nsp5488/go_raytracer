package hittable

import (
	"fmt"
	"math"

	ImageLoader "github.com/nsp5488/go_raytracer/internal/imageloader"
	"github.com/nsp5488/go_raytracer/internal/interval"
	"github.com/nsp5488/go_raytracer/internal/vec"
)

type Texture interface {
	Value(u, v float64, point *vec.Vec3) *vec.Vec3
}

type solidColor struct {
	Albedo *vec.Vec3
}

func NewSolidColor(albedo *vec.Vec3) *solidColor {
	return &solidColor{Albedo: albedo}
}
func NewSolidColorRGB(r, g, b float64) *solidColor {
	return &solidColor{Albedo: vec.New(r, g, b)}
}

func (sc *solidColor) Value(u, v float64, point *vec.Vec3) *vec.Vec3 {
	return sc.Albedo
}

type checkerboard struct {
	inv_scale float64
	even      Texture
	odd       Texture
}

func NewCheckerboard(scale float64, even, odd Texture) *checkerboard {
	return &checkerboard{
		inv_scale: 1 / scale,
		even:      even,
		odd:       odd,
	}
}
func NewCheckerboardColors(scale float64, even, odd *vec.Vec3) *checkerboard {
	return &checkerboard{
		inv_scale: 1 / scale,
		even:      NewSolidColor(even),
		odd:       NewSolidColor(odd),
	}
}

func (cb *checkerboard) Value(u, v float64, point *vec.Vec3) *vec.Vec3 {
	x := int(math.Floor(cb.inv_scale * point.X()))
	y := int(math.Floor(cb.inv_scale * point.Y()))
	z := int(math.Floor(cb.inv_scale * point.Z()))

	if (x+y+z)%2 == 0 {
		return cb.even.Value(u, v, point)
	} else {
		return cb.odd.Value(u, v, point)
	}
}

type imageTexture struct {
	img *ImageLoader.RTImage
}

func NewImageTexture(filename string) *imageTexture {
	return &imageTexture{img: ImageLoader.LoadImage(filename)}
}

func (it *imageTexture) Value(u, v float64, point *vec.Vec3) *vec.Vec3 {
	if it.img.Height <= 0 {
		return vec.New(0, 1, 1)
	}
	if u < 0 || u > 1 || v < 0 || v > 1 {
		fmt.Printf("Warning: UV coordinates out of bounds: u=%f, v=%f\n", u, v)
	}
	u = interval.New(0, 1).Clamp(u)
	v = 1 - interval.New(0, 1).Clamp(v)

	i := int(u * float64(it.img.Width))
	j := int(v * float64(it.img.Height))
	pixel := it.img.PixelData(i, j)

	scale := 1.0 / 255.0
	return vec.New(float64(pixel.Data[0])*scale, float64(pixel.Data[1])*scale, float64(pixel.Data[2])*scale)
}

type perlinType uint8

const (
	_ perlinType = iota
	PERLIN
	MARBLE
	TURBULENT
)

type noiseTexture struct {
	noise   *perlin
	scale   float64
	variant perlinType
}

func NewNoiseTexture(scale float64) *noiseTexture {
	return &noiseTexture{noise: NewPerlin(), scale: scale, variant: PERLIN}
}

func NewNoiseTextureWithType(scale float64, variant perlinType) *noiseTexture {
	return &noiseTexture{noise: NewPerlin(), scale: scale, variant: variant}

}

func (nt *noiseTexture) Value(u, v float64, point *vec.Vec3) *vec.Vec3 {
	switch nt.variant {
	case PERLIN:
		return vec.New(1, 1, 1).Scale(.5 * (1.0 + nt.noise.Noise(point.Scale(nt.scale))))
	case MARBLE:
		return vec.New(.5, .5, .5).Scale(1 + math.Sin(nt.scale*point.Z()+10*nt.noise.Turbulence(point, 7)))
	case TURBULENT:
		return vec.New(1, 1, 1).Scale(nt.noise.Turbulence(point, 7))
	}

	// should be impossible, but we'll default to perlin
	nt.variant = PERLIN
	return vec.New(1, 1, 1).Scale(.5 * (1.0 + nt.noise.Noise(point.Scale(nt.scale))))
}
