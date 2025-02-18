package main

import (
	"bytes"
	"flag"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"

	"github.com/nsp5488/go_raytracer/internal/camera"
	"github.com/nsp5488/go_raytracer/internal/hittable"
	"github.com/nsp5488/go_raytracer/internal/vec"
)

func testWorld() *hittable.HittableList {
	// define our materials
	// matte
	ground := hittable.Lambertian{Albedo: *vec.New(0.8, 0.8, 0)}
	center := hittable.Lambertian{Albedo: *vec.New(0.1, .2, .5)}

	// glass
	left := hittable.Dielectric{RefractionIndex: 1.50}

	// models an air bubble
	bubble := hittable.Dielectric{RefractionIndex: 1.0 / 1.5}

	// metal
	right := hittable.Metal{Albedo: *vec.New(0.8, 0.6, 0.2), Fuzz: 1.0}

	// Define the "world"
	world := &hittable.HittableList{}
	world.Init(5)

	world.Add(&hittable.Sphere{Center: *vec.New(0, -100.5, -1), Radius: 100, Material: &ground})
	world.Add(&hittable.Sphere{Center: *vec.New(0, 0, -1.2), Radius: 0.5, Material: &center})
	world.Add(&hittable.Sphere{Center: *vec.New(-1, 0, -1), Radius: 0.5, Material: &left})
	world.Add(&hittable.Sphere{Center: *vec.New(-1, 0, -1), Radius: .4, Material: &bubble})
	world.Add(&hittable.Sphere{Center: *vec.New(1, 0, -1), Radius: 0.5, Material: &right})
	return world
}

func coverWorld() *hittable.HittableList {
	world := &hittable.HittableList{}
	world.Init(50)
	glass := hittable.Dielectric{RefractionIndex: 1.5}
	ground := hittable.Lambertian{Albedo: *vec.New(0.5, 0.5, 0.5)}
	world.Add(&hittable.Sphere{Center: *vec.New(0, -1000, -1), Radius: 1000, Material: &ground})
	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			mat := rand.Float64()
			center := vec.New(float64(a)+0.9*rand.Float64(), 0.2, float64(b)+0.9*rand.Float64())

			if center.Add(vec.New(4, 0.2, 0).Negate()).Length() > 0.9 {
				var material hittable.Material

				if mat < 0.8 {
					albedo := vec.Random().Multiply(vec.Random())
					material = hittable.Lambertian{Albedo: *albedo}
					world.Add(&hittable.Sphere{Center: *center, Radius: 0.2, Material: material})
				} else if mat < 0.95 {
					albedo := vec.RangeRandom(0.5, 1.0)
					fuzz := rand.Float64()
					material = hittable.Metal{Albedo: *albedo, Fuzz: fuzz}
					world.Add(&hittable.Sphere{Center: *center, Radius: 0.2, Material: material})
				} else {
					world.Add(&hittable.Sphere{Center: *center, Radius: 0.2, Material: &glass})
				}
			}
		}
	}
	world.Add(&hittable.Sphere{Center: *vec.New(0, 1, 0), Radius: 1.0, Material: glass})
	mat2 := hittable.Lambertian{Albedo: *vec.New(0.4, 0.2, 0.1)}
	world.Add(&hittable.Sphere{Center: *vec.New(-4, 1, 0), Radius: 1.0, Material: mat2})
	mat3 := hittable.Metal{Albedo: *vec.New(.7, .6, .5), Fuzz: 0.0}
	world.Add(&hittable.Sphere{Center: *vec.New(4, 1, 0), Radius: 1.0, Material: mat3})
	return world
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	outBuf := bytes.Buffer{}

	world := coverWorld()

	c := camera.Camera{}
	c.Out = &outBuf

	c.AspectRatio = float64(16) / float64(9)
	c.Width = 120
	c.SamplesPerPixel = 5
	c.MaxDepth = 50

	c.VerticalFOV = 20
	c.PositionCamera(vec.New(13, 2, 3), vec.New(0, 0, 0), vec.New(0, 1, 0))

	c.DefocusAngle = 0.6
	c.FocusDistance = 10.0

	c.Render(world)
	file, err := os.Create("image.ppm")
	if err != nil {
		log.Fatal("Error creating output file\n")
	}
	file.Write(outBuf.Bytes())
}
