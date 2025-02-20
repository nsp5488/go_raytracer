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

type SolidColor struct {
	Albedo *vec.Vec3
}

func NewSolidColor(albedo *vec.Vec3) *SolidColor {
	return &SolidColor{Albedo: albedo}
}
func NewSolidColorRGB(r, g, b float64) *SolidColor {
	return &SolidColor{Albedo: vec.New(r, g, b)}
}

func (sc *SolidColor) Value(u, v float64, point *vec.Vec3) *vec.Vec3 {
	return sc.Albedo
}

type Checkerboard struct {
	inv_scale float64
	even      Texture
	odd       Texture
}

func NewCheckerboard(scale float64, even, odd Texture) *Checkerboard {
	return &Checkerboard{
		inv_scale: 1 / scale,
		even:      even,
		odd:       odd,
	}
}
func NewCheckerboardColors(scale float64, even, odd *vec.Vec3) *Checkerboard {
	return &Checkerboard{
		inv_scale: 1 / scale,
		even:      NewSolidColor(even),
		odd:       NewSolidColor(odd),
	}
}

func (cb *Checkerboard) Value(u, v float64, point *vec.Vec3) *vec.Vec3 {
	x := int(math.Floor(cb.inv_scale * point.X()))
	y := int(math.Floor(cb.inv_scale * point.Y()))
	z := int(math.Floor(cb.inv_scale * point.Z()))

	if (x+y+z)%2 == 0 {
		return cb.even.Value(u, v, point)
	} else {
		return cb.odd.Value(u, v, point)
	}
}

type ImageTexture struct {
	img ImageLoader.RTImage
}

func NewImageTexture(filename string) *ImageTexture {
	return &ImageTexture{img: *ImageLoader.LoadImage(filename)}
}

func (it *ImageTexture) Value(u, v float64, point *vec.Vec3) *vec.Vec3 {
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
