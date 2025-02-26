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
	world := hittable.NewHittableList(5)
	world.Add(hittable.NewSphere(vec.New(0, -100.5, -1), 100, ground))
	world.Add(hittable.NewSphere(vec.New(0, 0, -1.2), 0.5, center))
	world.Add(hittable.NewSphere(vec.New(-1, 0, -1), 0.5, left))
	world.Add(hittable.NewSphere(vec.New(-1, 0, -1), 0.4, bubble))
	world.Add(hittable.NewSphere(vec.New(1, 0, -1), 0.5, right))
	c.Background = vec.New(0.70, 0.80, 1.00)
	b := hittable.BuildBVH(world)
	w := hittable.NewHittableList(1)
	w.Add(b)
	c.Render(w, hittable.NewHittableList(0))
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

	world := hittable.NewHittableList(4 + 22*21)
	glass := hittable.NewDielectric(1.5)
	checker := hittable.NewCheckerboardColors(0.32, vec.New(.2, .3, .1), vec.New(.9, .9, .9))
	// ground := hittable.Lambertian{Albedo: *vec.New(0.5, 0.5, 0.5)}
	world.Add(hittable.NewSphere(vec.New(0, -1000, 0), 1000, hittable.NewTexturedLambertian(checker)))
	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			mat := rand.Float64()
			center := vec.New(float64(a)+0.9*rand.Float64(), 0.2, float64(b)+0.9*rand.Float64())

			if center.Add(vec.New(4, 0.2, 0).Negate()).Length() > 0.9 {
				var material hittable.Material

				if mat < 0.8 {
					albedo := vec.Random().Multiply(vec.Random())
					material = hittable.NewLambertian(albedo)
					world.Add(hittable.NewMotionSphere(center, center.Add(vec.New(0, util.RangeRange(0, 0.5), 0)), 0.2, material))

				} else if mat < 0.95 {
					albedo := vec.RangeRandom(0.5, 1.0)
					fuzz := rand.Float64()
					material = hittable.NewMetal(albedo, fuzz)
					world.Add(hittable.NewSphere(center, 0.2, material))

				} else {
					world.Add(hittable.NewSphere(center, 0.2, glass))
				}
			}
		}
	}
	world.Add(hittable.NewSphere(vec.New(0, 1, 0), 1.0, glass))

	mat2 := hittable.NewLambertian(vec.New(0.4, 0.2, 0.1))
	world.Add(hittable.NewSphere(vec.New(-4, 1, 0), 1.0, mat2))
	mat3 := hittable.NewMetal(vec.New(.7, .6, .5), 0)
	world.Add(hittable.NewSphere(vec.New(4, 1, 0), 1.0, mat3))

	b := hittable.BuildBVH(world)
	c.Render(b, hittable.NewHittableList(0))
}

func checkeredSpheres(c *camera.Camera) {
	world := hittable.NewHittableList(2)
	checker := hittable.NewCheckerboardColors(0.32, vec.New(.2, .3, .1), vec.New(.9, .9, .9))

	world.Add(hittable.NewSphere(vec.New(0, -10, 0), 10, hittable.NewTexturedLambertian(checker)))
	world.Add(hittable.NewSphere(vec.New(0, 10, 0), 10, hittable.NewTexturedLambertian(checker)))
	b := hittable.BuildBVH(world)
	w := hittable.NewHittableList(1)
	w.Add(b)
	c.AspectRatio = float64(16) / float64(9)
	c.Width = 400
	c.SamplesPerPixel = 50
	c.MaxDepth = 50

	c.VerticalFOV = 20
	c.PositionCamera(vec.New(13, 2, 3), vec.New(0, 0, 0), vec.New(0, 1, 0))
	c.Background = vec.New(0.70, 0.80, 1.00)

	c.DefocusAngle = 0
	c.Render(w, hittable.NewHittableList(0))
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
	l := hittable.NewHittableList(1)
	earth_tex := hittable.NewImageTexture("earthmap.jpg")
	earth_surf := hittable.NewTexturedLambertian(earth_tex)
	l.Add(hittable.NewSphere(vec.Empty(), 2, earth_surf))
	c.Render(l, hittable.NewHittableList(0))
}

func perlin(cam *camera.Camera) {
	world := hittable.NewHittableList(2)
	p := hittable.NewNoiseTextureWithType(4, hittable.MARBLE)
	s1 := hittable.NewSphere(vec.New(0, -1000, 0), 1000, hittable.NewTexturedLambertian(p))
	s2 := hittable.NewSphere(vec.New(0, 2, 0), 2, hittable.NewTexturedLambertian(p))
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

	cam.Render(world, hittable.NewHittableList(0))
}

func quads(cam *camera.Camera) {
	world := hittable.NewHittableList(5)

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

	cam.AspectRatio = 1.0
	cam.Width = 400
	cam.SamplesPerPixel = 100
	cam.MaxDepth = 50
	cam.Background = vec.New(0.70, 0.80, 1.00)

	cam.VerticalFOV = 80
	cam.PositionCamera(vec.New(0, 0, 9), vec.New(0, 0, 0), vec.New(0, 1, 0))
	cam.DefocusAngle = 0
	cam.Render(bvh, hittable.NewHittableList(0))
}

func simpleLight(cam *camera.Camera) {
	world := hittable.NewHittableList(4)
	p := hittable.NewNoiseTextureWithType(4, hittable.MARBLE)
	l := hittable.NewDiffuseLight(vec.New(4, 4, 4))

	s1 := hittable.NewSphere(vec.New(0, -1000, 0), 1000, hittable.NewTexturedLambertian(p))
	s2 := hittable.NewSphere(vec.New(0, 2, 0), 2, hittable.NewTexturedLambertian(p))
	q := hittable.NewQuad(vec.New(3, 1, -2), vec.New(2, 0, 0), vec.New(0, 2, 0), l)
	s := hittable.NewSphere(vec.New(0, 7, 0), 2, l)
	world.Add(s1)
	world.Add(s)
	world.Add(q)
	world.Add(s2)

	cam.AspectRatio = 16.0 / 9.0
	cam.Width = 400
	cam.SamplesPerPixel = 100
	cam.MaxDepth = 50
	cam.Background = vec.New(0, 0, 0)

	cam.VerticalFOV = 20
	cam.PositionCamera(vec.New(26, 3, 6), vec.New(0, 2, 0), vec.New(0, 1, 0))

	cam.DefocusAngle = 0

	cam.Render(world, q)

}

func cornellBox(cam *camera.Camera) {
	world := hittable.NewHittableList(8)

	red := hittable.NewLambertian(vec.New(.65, .05, .05))
	white := hittable.NewLambertian(vec.New(.73, .73, .73))
	green := hittable.NewLambertian(vec.New(.12, .45, .15))
	light := hittable.NewDiffuseLight(vec.New(15, 15, 15))

	// walls and light
	world.Add(hittable.NewQuad(vec.New(555, 0, 0), vec.New(0, 555, 0), vec.New(0, 0, 555), green))
	world.Add(hittable.NewQuad(vec.New(0, 0, 0), vec.New(0, 555, 0), vec.New(0, 0, 555), red))
	world.Add(hittable.NewQuad(vec.New(0, 0, 0), vec.New(555, 0, 0), vec.New(0, 0, 555), white))
	world.Add(hittable.NewQuad(vec.New(555, 555, 555), vec.New(-555, 0, 0), vec.New(0, 0, -555), white))
	world.Add(hittable.NewQuad(vec.New(0, 0, 555), vec.New(555, 0, 0), vec.New(0, 555, 0), white))

	// light  source:
	lights := hittable.NewHittableList(2)
	lights.Add(hittable.NewQuad(vec.New(343, 550, 332), vec.New(-130, 0, 0), vec.New(0, 0, -105), light))
	world.Add(lights)
	// boxes
	b1 := hittable.NewBox(vec.New(0, 0, 0), vec.New(165, 330, 165), white)
	b1 = hittable.RotateY(b1, 15)
	b1 = hittable.Translate(b1, vec.New(265, 0, 295))
	world.Add(b1)

	// b2 := hittable.NewBox(vec.New(0, 0, 0), vec.New(165, 165, 165), white)
	// b2 = hittable.RotateY(b2, -18)
	// b2 = hittable.Translate(b2, vec.New(130, 0, 65))
	// world.Add(b2)
	s := hittable.NewSphere(vec.New(190, 90, 190), 90, hittable.NewDielectric(1.5))
	lights.Add(s)
	world.Add(s)

	cam.AspectRatio = 1.0
	cam.Width = 600
	cam.SamplesPerPixel = 1000
	cam.MaxDepth = 50

	cam.Background = vec.Empty()
	cam.VerticalFOV = 40
	cam.PositionCamera(vec.New(278, 278, -800), vec.New(278, 278, 0), vec.New(0, 1, 0))
	cam.DefocusAngle = 0

	cam.Render(hittable.BuildBVH(world), lights)
}
func cornellSmoke(cam *camera.Camera) {
	world := hittable.NewHittableList(10)

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
	cam.SamplesPerPixel = 10
	cam.MaxDepth = 50

	cam.Background = vec.Empty()
	cam.VerticalFOV = 40
	cam.PositionCamera(vec.New(278, 278, -800), vec.New(278, 278, 0), vec.New(0, 1, 0))
	cam.DefocusAngle = 0

	bvh := hittable.BuildBVH(world)
	cam.Render(bvh, hittable.NewHittableList(0))
}
func book2Scene(cam *camera.Camera) {
	boxes1 := hittable.NewHittableList(20 * 20)
	groundColor := hittable.NewLambertian(vec.New(.48, .83, .53))

	// floor
	boxesPerSide := 20
	for i := 0; i < boxesPerSide; i++ {
		for j := 0; j < boxesPerSide; j++ {
			w := 100.0
			x0 := -1000.0 + float64(i)*w
			z0 := -1000.0 + float64(j)*w
			y0 := 0.0
			x1 := x0 + w
			y1 := util.RangeRange(1, 101)
			z1 := z0 + w
			boxes1.Add(hittable.NewBox(vec.New(x0, y0, z0), vec.New(x1, y1, z1), groundColor))
		}
	}
	world := hittable.NewHittableList(12)
	world.Add(hittable.BuildBVH(boxes1))

	// light
	light := hittable.NewDiffuseLight(vec.New(7, 7, 7))
	world.Add(hittable.NewQuad(vec.New(123, 554, 147), vec.New(300, 0, 0), vec.New(0, 0, 265), light))

	// motion blur
	c1 := vec.New(400, 400, 200)
	c2 := c1.Add(vec.New(30, 0, 0))
	sphereMat := hittable.NewLambertian(vec.New(.7, .3, .1))
	world.Add(hittable.NewMotionSphere(c1, c2, 50, sphereMat))

	// glass orb
	world.Add(hittable.NewSphere(vec.New(260, 150, 45), 50, hittable.NewDielectric(1.5)))

	// metal orb
	world.Add(hittable.NewSphere(vec.New(0, 150, 145), 50, hittable.NewMetal(vec.New(0.8, 0.8, 0.9), 1.0)))

	// water orb
	boundary := hittable.NewSphere(vec.New(360, 150, 145), 70, hittable.NewDielectric(1.5))
	world.Add(boundary)
	world.Add(hittable.ConstantMedium(boundary, .2, vec.New(0.2, 0.4, 0.9)))

	// fog
	b2 := hittable.NewSphere(vec.New(0, 0, 0), 5000, hittable.NewDielectric(1.5))
	world.Add(hittable.ConstantMedium(b2, .0001, vec.New(1, 1, 1)))

	// earth
	eMat := hittable.NewTexturedLambertian(hittable.NewImageTexture("earthmap.jpg"))
	world.Add(hittable.NewSphere(vec.New(400, 200, 400), 100, eMat))

	// perlin
	p := hittable.NewTexturedLambertian(hittable.NewNoiseTextureWithType(.2, hittable.MARBLE))
	world.Add(hittable.NewSphere(vec.New(220, 280, 300), 80, p))

	// weird spheres
	boxes2 := hittable.NewHittableList(1000)
	white := hittable.NewLambertian(vec.New(.73, .73, .73))
	ns := 1000
	for i := 0; i < ns; i++ {
		boxes2.Add(hittable.NewSphere(vec.RangeRandom(0, 165), 10, white))
	}
	world.Add(
		hittable.Translate(
			hittable.RotateY(hittable.BuildBVH(boxes2), 15),
			vec.New(-100, 270, 395)),
	)
	cam.AspectRatio = 1.0
	cam.Width = 800
	cam.SamplesPerPixel = 10000
	cam.MaxDepth = 40
	cam.Background = vec.Empty()

	cam.VerticalFOV = 40
	cam.PositionCamera(vec.New(478, 278, -600), vec.New(278, 278, 0), vec.New(0, 1, 0))

	cam.DefocusAngle = 0

	cam.Render(world, hittable.NewHittableList(0))
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
	if err != nil {
		log.Fatal("Error creating output file\n")
	}
	defer file.Close()

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
	cornellBox(&c)
	// cornellSmoke(&c)
	// book2Scene(&c)
	file.Write(outBuf.Bytes())
}
