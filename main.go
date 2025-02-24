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
	"github.com/nsp5488/go_raytracer/internal/util"
	"github.com/nsp5488/go_raytracer/internal/vec"
)

// Creates a simple test world.
func testWorld(c *camera.Camera) {
	// define our materials
	// matte
	ground := hittable.NewLambertian(vec.New(0.8, 0.8, 0))
	center := hittable.NewLambertian(vec.New(0.1, .2, .5))
	// glass
	left := hittable.NewDielectric(1.50)

	// models an air bubble
	bubble := hittable.NewDielectric(1.0 / 1.5)

	// metal
	right := hittable.NewMetal(vec.New(0.8, 0.6, 0.2), 1.0)
	// Define the "world"
	world := &hittable.HittableList{}
	world.Init(5)
	world.Add(hittable.NewSphere(*vec.New(0, -100.5, -1), 100, ground))
	world.Add(hittable.NewSphere(*vec.New(0, 0, -1.2), 0.5, center))
	world.Add(hittable.NewSphere(*vec.New(-1, 0, -1), 0.5, left))
	world.Add(hittable.NewSphere(*vec.New(-1, 0, -1), 0.4, bubble))
	world.Add(hittable.NewSphere(*vec.New(1, 0, -1), 0.5, right))
	c.Background = vec.New(0.70, 0.80, 1.00)

	c.Render(hittable.NewHittableList([]hittable.Hittable{hittable.BuildBVH(world)}))
}

// Creates the world from the cover of Ray Tracing in One Weekend.
func coverWorld(c *camera.Camera) {
	c.AspectRatio = float64(16) / float64(9)
	c.Width = 400
	c.SamplesPerPixel = 50
	c.MaxDepth = 50

	c.VerticalFOV = 20
	c.PositionCamera(vec.New(13, 2, 3), vec.New(0, 0, 0), vec.New(0, 1, 0))

	c.DefocusAngle = 0.6
	c.FocusDistance = 10.0
	c.Background = vec.New(0.70, 0.80, 1.00)

	world := &hittable.HittableList{}
	world.Init(4 + 22*21)
	glass := hittable.Dielectric{RefractionIndex: 1.5}
	checker := hittable.NewCheckerboardColors(0.32, vec.New(.2, .3, .1), vec.New(.9, .9, .9))
	// ground := hittable.Lambertian{Albedo: *vec.New(0.5, 0.5, 0.5)}
	world.Add(hittable.NewSphere(*vec.New(0, -1000, 0), 1000, hittable.NewTexturedLambertian(checker)))
	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			mat := rand.Float64()
			center := vec.New(float64(a)+0.9*rand.Float64(), 0.2, float64(b)+0.9*rand.Float64())

			if center.Add(vec.New(4, 0.2, 0).Negate()).Length() > 0.9 {
				var material hittable.Material

				if mat < 0.8 {
					albedo := vec.Random().Multiply(vec.Random())
					material = hittable.NewLambertian(albedo)
					world.Add(hittable.NewMotionSphere(*center, *center.Add(vec.New(0, util.RangeRange(0, 0.5), 0)), 0.2, material))

				} else if mat < 0.95 {
					albedo := vec.RangeRandom(0.5, 1.0)
					fuzz := rand.Float64()
					material = hittable.Metal{Albedo: *albedo, Fuzz: fuzz}
					world.Add(hittable.NewSphere(*center, 0.2, material))

				} else {
					world.Add(hittable.NewSphere(*center, 0.2, &glass))
				}
			}
		}
	}
	world.Add(hittable.NewSphere(*vec.New(0, 1, 0), 1.0, glass))

	mat2 := hittable.NewLambertian(vec.New(0.4, 0.2, 0.1))
	world.Add(hittable.NewSphere(*vec.New(-4, 1, 0), 1.0, mat2))

	mat3 := hittable.Metal{Albedo: *vec.New(.7, .6, .5), Fuzz: 0.0}
	world.Add(hittable.NewSphere(*vec.New(4, 1, 0), 1.0, mat3))

	b := hittable.BuildBVH(world)
	w := &hittable.HittableList{}
	w.Init(1)
	w.Add(b)
	c.Render(w)
}

func checkeredSpheres(c *camera.Camera) {
	world := &hittable.HittableList{}
	world.Init(2)
	checker := hittable.NewCheckerboardColors(0.32, vec.New(.2, .3, .1), vec.New(.9, .9, .9))

	world.Add(hittable.NewSphere(*vec.New(0, -10, 0), 10, hittable.NewTexturedLambertian(checker)))
	world.Add(hittable.NewSphere(*vec.New(0, 10, 0), 10, hittable.NewTexturedLambertian(checker)))
	b := hittable.BuildBVH(world)
	w := &hittable.HittableList{}
	w.Init(1)
	w.Add(b)
	c.AspectRatio = float64(16) / float64(9)
	c.Width = 400
	c.SamplesPerPixel = 50
	c.MaxDepth = 50

	c.VerticalFOV = 20
	c.PositionCamera(vec.New(13, 2, 3), vec.New(0, 0, 0), vec.New(0, 1, 0))
	c.Background = vec.New(0.70, 0.80, 1.00)

	c.DefocusAngle = 0
	c.Render(w)
}

func earth(c *camera.Camera) {
	c.AspectRatio = float64(16) / float64(9)
	c.Width = 500
	c.SamplesPerPixel = 50
	c.MaxDepth = 50

	c.VerticalFOV = 20
	c.PositionCamera(vec.New(0, 0, 12), vec.New(0, 0, 0), vec.New(0, 1, 0))
	c.Background = vec.New(0.70, 0.80, 1.00)

	c.DefocusAngle = 0
	l := &hittable.HittableList{}
	l.Init(1)
	earth_tex := hittable.NewImageTexture("earthmap.jpg")
	earth_surf := hittable.NewTexturedLambertian(earth_tex)
	l.Add(hittable.NewSphere(*vec.Empty(), 2, earth_surf))
	c.Render(l)
}

func perlin(cam *camera.Camera) {
	world := &hittable.HittableList{}
	world.Init(2)
	p := hittable.NewNoiseTextureWithType(4, hittable.MARBLE)
	s1 := hittable.NewSphere(*vec.New(0, -1000, 0), 1000, hittable.NewTexturedLambertian(p))
	s2 := hittable.NewSphere(*vec.New(0, 2, 0), 2, hittable.NewTexturedLambertian(p))
	world.Add(s1)
	world.Add(s2)
	cam.AspectRatio = 16.0 / 9.0
	cam.Width = 400
	cam.SamplesPerPixel = 100
	cam.MaxDepth = 50

	cam.VerticalFOV = 20
	cam.PositionCamera(vec.New(13, 2, 3), vec.New(0, 0, 0), vec.New(0, 1, 0))
	cam.Background = vec.New(0.70, 0.80, 1.00)

	cam.DefocusAngle = 0

	cam.Render(world)
}

func quads(cam *camera.Camera) {
	world := &hittable.HittableList{}
	world.Init(5)

	leftRed := hittable.NewLambertian(vec.New(1, 0.2, 0.2))
	backGreen := hittable.NewTexturedLambertian(hittable.NewNoiseTextureWithType(5, hittable.MARBLE))
	rightBlue := hittable.NewLambertian(vec.New(0.2, 0.2, 1.0))
	upperOrange := hittable.NewLambertian(vec.New(1.0, 0.5, 0.0))
	lowerTeal := hittable.NewLambertian(vec.New(0.2, 0.8, 0.8))

	world.Add(hittable.NewQuad(vec.New(-3, -2, 5), vec.New(0, 0, -4), vec.New(0, 4, 0), leftRed))
	world.Add(hittable.NewQuad(vec.New(-2, -2, 0), vec.New(4, 0, 0), vec.New(0, 4, 0), backGreen))
	world.Add(hittable.NewQuad(vec.New(3, -2, 1), vec.New(0, 0, 4), vec.New(0, 4, 0), rightBlue))
	world.Add(hittable.NewQuad(vec.New(-2, 3, 1), vec.New(4, 0, 0), vec.New(0, 0, 4), upperOrange))
	world.Add(hittable.NewQuad(vec.New(-2, -3, 5), vec.New(4, 0, 0), vec.New(0, 0, -4), lowerTeal))
	bvh := hittable.BuildBVH(world)
	w := &hittable.HittableList{}
	w.Init(1)
	w.Add(bvh)

	cam.AspectRatio = 1.0
	cam.Width = 400
	cam.SamplesPerPixel = 100
	cam.MaxDepth = 50
	cam.Background = vec.New(0.70, 0.80, 1.00)

	cam.VerticalFOV = 80
	cam.PositionCamera(vec.New(0, 0, 9), vec.New(0, 0, 0), vec.New(0, 1, 0))
	cam.DefocusAngle = 0
	cam.Render(w)
}

func simpleLight(cam *camera.Camera) {
	world := &hittable.HittableList{}
	world.Init(4)
	p := hittable.NewNoiseTextureWithType(4, hittable.MARBLE)
	l := hittable.NewDiffuseLight(vec.New(4, 4, 4))

	s1 := hittable.NewSphere(*vec.New(0, -1000, 0), 1000, hittable.NewTexturedLambertian(p))
	s2 := hittable.NewSphere(*vec.New(0, 2, 0), 2, hittable.NewTexturedLambertian(p))
	q := hittable.NewQuad(vec.New(3, 1, -2), vec.New(2, 0, 0), vec.New(0, 2, 0), l)
	s := hittable.NewSphere(*vec.New(0, 7, 0), 2, l)
	world.Add(q)
	world.Add(s1)
	world.Add(s2)
	world.Add(s)
	cam.AspectRatio = 16.0 / 9.0
	cam.Width = 400
	cam.SamplesPerPixel = 100
	cam.MaxDepth = 50
	cam.Background = vec.New(0, 0, 0)

	cam.VerticalFOV = 20
	cam.PositionCamera(vec.New(26, 3, 6), vec.New(0, 2, 0), vec.New(0, 1, 0))

	cam.DefocusAngle = 0

	cam.Render(world)

}

func cornellBox(cam *camera.Camera) {
	world := &hittable.HittableList{}
	world.Init(8)

	red := hittable.NewLambertian(vec.New(.65, .05, .05))
	white := hittable.NewLambertian(vec.New(.73, .73, .73))
	green := hittable.NewLambertian(vec.New(.12, .45, .15))
	light := hittable.NewDiffuseLight(vec.New(15, 15, 15))

	// walls and light
	world.Add(hittable.NewQuad(vec.New(555, 0, 0), vec.New(0, 555, 0), vec.New(0, 0, 555), green))
	world.Add(hittable.NewQuad(vec.New(0, 0, 0), vec.New(0, 555, 0), vec.New(0, 0, 555), red))
	world.Add(hittable.NewQuad(vec.New(343, 550, 332), vec.New(-130, 0, 0), vec.New(0, 0, -105), light))
	world.Add(hittable.NewQuad(vec.New(0, 0, 0), vec.New(555, 0, 0), vec.New(0, 0, 555), white))
	world.Add(hittable.NewQuad(vec.New(555, 555, 555), vec.New(-555, 0, 0), vec.New(0, 0, -555), white))
	world.Add(hittable.NewQuad(vec.New(0, 0, 555), vec.New(555, 0, 0), vec.New(0, 555, 0), white))

	// boxes
	b1 := hittable.NewBox(vec.New(0, 0, 0), vec.New(165, 330, 165), white)
	b1 = hittable.RotateY(b1, 15)
	b1 = hittable.Translate(b1, vec.New(265, 0, 295))

	b2 := hittable.NewBox(vec.New(0, 0, 0), vec.New(165, 165, 165), white)
	b2 = hittable.RotateY(b2, -18)
	b2 = hittable.Translate(b2, vec.New(130, 0, 65))

	world.Add(b1)
	world.Add(b2)
	cam.AspectRatio = 1.0
	cam.Width = 600
	cam.SamplesPerPixel = 200
	cam.MaxDepth = 50

	cam.Background = vec.Empty()
	cam.VerticalFOV = 40
	cam.PositionCamera(vec.New(278, 278, -800), vec.New(278, 278, 0), vec.New(0, 1, 0))
	cam.DefocusAngle = 0

	w := &hittable.HittableList{}
	w.Init(1)
	w.Add(hittable.BuildBVH(world))
	cam.Render(world)
}
func cornellSmoke(cam *camera.Camera) {
	world := &hittable.HittableList{}
	world.Init(10)

	red := hittable.NewLambertian(vec.New(.65, .05, .05))
	white := hittable.NewLambertian(vec.New(.73, .73, .73))
	green := hittable.NewLambertian(vec.New(.12, .45, .15))
	light := hittable.NewDiffuseLight(vec.New(15, 15, 15))

	// walls and light
	world.Add(hittable.NewQuad(vec.New(555, 0, 0), vec.New(0, 555, 0), vec.New(0, 0, 555), green))
	world.Add(hittable.NewQuad(vec.New(0, 0, 0), vec.New(0, 555, 0), vec.New(0, 0, 555), red))
	world.Add(hittable.NewQuad(vec.New(343, 550, 332), vec.New(-130, 0, 0), vec.New(0, 0, -105), light))
	world.Add(hittable.NewQuad(vec.New(0, 0, 0), vec.New(555, 0, 0), vec.New(0, 0, 555), white))
	world.Add(hittable.NewQuad(vec.New(555, 555, 555), vec.New(-555, 0, 0), vec.New(0, 0, -555), white))
	world.Add(hittable.NewQuad(vec.New(0, 0, 555), vec.New(555, 0, 0), vec.New(0, 555, 0), white))

	// boxes
	b1 := hittable.NewBox(vec.New(0, 0, 0), vec.New(165, 330, 165), white)
	b1 = hittable.RotateY(b1, 15)
	b1 = hittable.Translate(b1, vec.New(265, 0, 295))

	b2 := hittable.NewBox(vec.New(0, 0, 0), vec.New(165, 165, 165), white)
	b2 = hittable.RotateY(b2, -18)
	b2 = hittable.Translate(b2, vec.New(130, 0, 65))

	// smoke
	world.Add(hittable.ConstantMedium(b1, .01, vec.Empty()))
	world.Add(hittable.ConstantMedium(b2, .01, vec.New(1, 1, 1)))

	cam.AspectRatio = 1.0
	cam.Width = 600
	cam.SamplesPerPixel = 200
	cam.MaxDepth = 50

	cam.Background = vec.Empty()
	cam.VerticalFOV = 40
	cam.PositionCamera(vec.New(278, 278, -800), vec.New(278, 278, 0), vec.New(0, 1, 0))
	cam.DefocusAngle = 0

	w := &hittable.HittableList{}
	w.Init(1)
	w.Add(hittable.BuildBVH(world))
	cam.Render(world)
}

func main() {
	cpuprofile := flag.String("cpuprofile", "", "Write cpu profile to file")
	outFile := flag.String("o", "image.ppm", "Specify a custom output file")
	coreCount := flag.Int("N", 1, "Set the number of cores to allocate to rendering")
	// imgWidth := flag.Int("width", 400, "Set the image width, default:1200")
	// samplesPerPix := flag.Int("samples", 50, "Specify the number of samples to take per pixel")
	// maxDepth := flag.Int("depth", 50, "Sets the maximum recursive depth of ray bounces")
	// vfov := flag.Float64("fov", 20, "Sets the vertical FOV of the camera")

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	file, err := os.Create(*outFile)
	defer file.Close()
	if err != nil {
		log.Fatal("Error creating output file\n")
	}

	outBuf := bytes.Buffer{}

	c := camera.Camera{}
	c.Out = &outBuf
	c.MaxThreads = *coreCount
	// testWorld(&c)
	// coverWorld(&c)
	// checkeredSpheres(&c)
	// earth(&c)
	// perlin(&c)
	// quads(&c)
	// simpleLight(&c)
	// cornellBox(&c)
	cornellSmoke(&c)
	file.Write(outBuf.Bytes())
}
