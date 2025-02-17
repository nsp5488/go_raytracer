package main

import (
	"log"
	"os"

	"github.com/nsp5488/go_raytracer/internal/camera"
	"github.com/nsp5488/go_raytracer/internal/hittable"
	"github.com/nsp5488/go_raytracer/internal/vec"
)

func main() {
	file, err := os.Create("image.ppm")
	if err != nil {
		log.Fatal("Error creating output file\n")
	}
	c := camera.Camera{}
	c.Width = 400
	c.AspectRatio = float64(16) / float64(9)
	c.Out = file
	c.SamplesPerPixel = 100
	c.MaxDepth = 50

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
	world.Init(4)
	world.Add(&hittable.Sphere{Center: *vec.New(0, -100.5, -1), Radius: 100, Material: &ground})
	world.Add(&hittable.Sphere{Center: *vec.New(0, 0, -1.2), Radius: 0.5, Material: &center})
	world.Add(&hittable.Sphere{Center: *vec.New(-1, 0, -1), Radius: 0.5, Material: &left})
	world.Add(&hittable.Sphere{Center: *vec.New(-1, 0, -1), Radius: .4, Material: &bubble})
	world.Add(&hittable.Sphere{Center: *vec.New(1, 0, -1), Radius: 0.5, Material: &right})

	c.Render(world)
}
