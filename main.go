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
	"github.com/nsp5488/go_raytracer/internal/objLoader"
	"github.com/nsp5488/go_raytracer/internal/util"
	"github.com/nsp5488/go_raytracer/internal/vec"
)

// Creates the world from the cover of Ray Tracing in One Weekend with some additional modifications to showcase later features.
func book1Scene(c *camera.Camera) {
	c.AspectRatio = float64(16) / float64(9)
	c.Width = 400
	c.SamplesPerPixel = 100
	c.MaxDepth = 50

	c.VerticalFOV = 20
	c.PositionCamera(vec.New(13, 2, 3), vec.New(0, 0, 0), vec.New(0, 1, 0))

	c.DefocusAngle = 0.6
	c.FocusDistance = 10.0
	c.Background = vec.New(0.70, 0.80, 1.00)

	world := hittable.NewHittableList(4 + 22*21)
	lights := hittable.NewHittableList(1)

	glass := hittable.NewDielectric(1.5)
	checker := hittable.NewCheckerboardColors(0.32, vec.New(.2, .3, .1), vec.New(.9, .9, .9))
	world.Add(hittable.NewSphere(vec.New(0, -1000, 0), 1000, hittable.NewTexturedLambertian(checker)))
	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			mat := rand.Float64()
			center := vec.New(float64(a)+0.9*rand.Float64(), 0.2, float64(b)+0.9*rand.Float64())

			if center.Add(vec.New(4, 0.2, 0).Negate()).Length() > 0.9 {
				var material hittable.Material

				if mat < 0.6 {
					// matte solid color orbs
					albedo := vec.Random().Multiply(vec.Random())
					material = hittable.NewLambertian(albedo)
					world.Add(hittable.NewMotionSphere(center, center.Add(vec.New(0, util.RangeRange(0, 0.5), 0)), 0.2, material))

				} else if mat < 0.8 {
					// perlin orbs
					if mat < .65 {
						material = hittable.NewTexturedLambertian(hittable.NewNoiseTextureWithType(float64(rand.Intn(10)), hittable.MARBLE))
					} else if mat < .7 {
						material = hittable.NewTexturedLambertian(hittable.NewNoiseTextureWithType(float64(rand.Intn(10)), hittable.TURBULENT))
					} else {
						material = hittable.NewTexturedLambertian(hittable.NewNoiseTextureWithType(float64(rand.Intn(10)), hittable.PERLIN))
					}
				} else if mat < 0.95 {
					// Reflective metallic orbs
					albedo := vec.RangeRandom(0.5, 1.0)
					fuzz := rand.Float64()
					material = hittable.NewMetal(albedo, fuzz)
					world.Add(hittable.NewSphere(center, 0.2, material))

				} else {
					// Glass orbs
					s := hittable.NewSphere(center, 0.2, glass)
					world.Add(s)
				}
			}
		}
	}

	// Big central spheres
	world.Add(hittable.NewSphere(vec.New(0, 1, 0), 1.0, glass))
	mat2 := hittable.NewLambertian(vec.New(0.4, 0.2, 0.1))
	world.Add(hittable.NewSphere(vec.New(-4, 1, 0), 1.0, mat2))
	mat3 := hittable.NewMetal(vec.New(.7, .6, .5), 0)
	world.Add(hittable.NewSphere(vec.New(4, 1, 0), 1.0, mat3))

	// A large light source "sun"
	sun := hittable.NewSphere(vec.New(0, 100, 0), 50, hittable.NewDiffuseLight(vec.New(5, 5, 5)))
	world.Add(sun)
	lights.Add(sun)

	b := hittable.BuildBVH(world)
	c.Render(b, lights)
}

// Creates the scene on the cover of Ray Tracing: The Next Week by Peter Shirley
func book2Scene(cam *camera.Camera) {
	boxes1 := hittable.NewHittableList(20 * 20)
	groundColor := hittable.NewLambertian(vec.New(.48, .83, .53))

	// floor
	boxesPerSide := 20
	for i := range boxesPerSide {
		for j := range boxesPerSide {
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
	lights := hittable.NewHittableList(1)

	// light
	light := hittable.NewQuad(vec.New(123, 554, 147), vec.New(300, 0, 0), vec.New(0, 0, 265), hittable.NewDiffuseLight(vec.New(7, 7, 7)))
	world.Add(light)
	lights.Add(light)

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
	for range ns {
		boxes2.Add(hittable.NewSphere(vec.RangeRandom(0, 165), 10, white))
	}
	world.Add(
		hittable.Translate(
			hittable.RotateY(hittable.BuildBVH(boxes2), 15),
			vec.New(-100, 270, 395)),
	)
	cam.AspectRatio = 1.0
	cam.Width = 800
	cam.SamplesPerPixel = 100
	cam.MaxDepth = 40
	cam.Background = vec.Empty()

	cam.VerticalFOV = 40
	cam.PositionCamera(vec.New(478, 278, -600), vec.New(278, 278, 0), vec.New(0, 1, 0))

	cam.DefocusAngle = 0

	cam.Render(world, lights)
}

// Creates the scene on thee cover of Ray Tracing: The Rest of Your Life by Peter Shirley
func book3Scene(cam *camera.Camera) {
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

	s := hittable.NewSphere(vec.New(190, 90, 190), 90, hittable.NewDielectric(1.5))
	lights.Add(s)
	world.Add(s)

	cam.AspectRatio = 1.0
	cam.Width = 600
	cam.SamplesPerPixel = 10
	cam.MaxDepth = 50

	cam.Background = vec.Empty()
	cam.VerticalFOV = 40
	cam.PositionCamera(vec.New(278, 278, -800), vec.New(278, 278, 0), vec.New(0, 1, 0))
	cam.DefocusAngle = 0

	cam.Render(hittable.BuildBVH(world), lights)
}

func quads(cam *camera.Camera) {
	world := hittable.NewHittableList(5)
	lights := hittable.NewHittableList(1)
	leftEarth := hittable.NewTexturedLambertian(hittable.NewImageTexture("earthmap.jpg"))
	backLight := hittable.NewDiffuseLight(vec.New(3, 3, 3))
	rightPerlin := hittable.NewTexturedLambertian(hittable.NewNoiseTextureWithType(5, hittable.MARBLE))
	upperMetal := hittable.NewMetal(vec.New(0.8, 0.6, 0.2), 0)
	lowerTeal := hittable.NewLambertian(vec.New(0.2, 0.8, 0.8))

	world.Add(hittable.NewQuad(vec.New(-3, -2, 5), vec.New(0, 0, -4), vec.New(0, 4, 0), leftEarth))
	light := hittable.NewQuad(vec.New(-2, -2, 0), vec.New(4, 0, 0), vec.New(0, 4, 0), backLight)
	world.Add(light)
	world.Add(hittable.NewQuad(vec.New(3, -2, 1), vec.New(0, 0, 4), vec.New(0, 4, 0), rightPerlin))
	world.Add(hittable.NewQuad(vec.New(-2, 3, 1), vec.New(4, 0, 0), vec.New(0, 0, 4), upperMetal))
	world.Add(hittable.NewQuad(vec.New(-2, -3, 5), vec.New(4, 0, 0), vec.New(0, 0, -4), lowerTeal))
	bvh := hittable.BuildBVH(world)
	lights.Add(light)
	cam.AspectRatio = 1.0
	cam.Width = 400
	cam.SamplesPerPixel = 100
	cam.MaxDepth = 50
	cam.Background = vec.New(0.70, 0.80, 1.00)

	cam.VerticalFOV = 80
	cam.PositionCamera(vec.New(0, 0, 9), vec.New(0, 0, 0), vec.New(0, 1, 0))
	cam.DefocusAngle = 0
	cam.Render(bvh, lights)
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

// A cornell box
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

	// light source:
	lights := hittable.NewHittableList(2)
	lights.Add(hittable.NewQuad(vec.New(343, 550, 332), vec.New(-130, 0, 0), vec.New(0, 0, -105), light))
	world.Add(lights)

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
	cam.SamplesPerPixel = 100
	cam.MaxDepth = 50

	cam.Background = vec.Empty()
	cam.VerticalFOV = 40
	cam.PositionCamera(vec.New(278, 278, -800), vec.New(278, 278, 0), vec.New(0, 1, 0))
	cam.DefocusAngle = 0

	cam.Render(hittable.BuildBVH(world), lights)
}

// A cornell box scene with the boxes replaced by boxes of smoke
func cornellSmoke(cam *camera.Camera) {
	world := hittable.NewHittableList(10)
	lights := hittable.NewHittableList(1)

	red := hittable.NewLambertian(vec.New(.65, .05, .05))
	white := hittable.NewLambertian(vec.New(.73, .73, .73))
	green := hittable.NewLambertian(vec.New(.12, .45, .15))
	light := hittable.NewDiffuseLight(vec.New(15, 15, 15))

	// walls and light
	world.Add(hittable.NewQuad(vec.New(555, 0, 0), vec.New(0, 555, 0), vec.New(0, 0, 555), green))
	world.Add(hittable.NewQuad(vec.New(0, 0, 0), vec.New(0, 555, 0), vec.New(0, 0, 555), red))
	lightQuad := hittable.NewQuad(vec.New(343, 550, 332), vec.New(-130, 0, 0), vec.New(0, 0, -105), light)
	world.Add(lightQuad)
	lights.Add(lightQuad)
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
	cam.Render(bvh, lights)
}

// This function won't work without an external .obj / .mtl file. It's currently setup to use the MTL file located here:
// https://casual-effects.com/data/index.html under "Chinese Dragon"
func modelExample(cam *camera.Camera) {
	world := hittable.NewHittableList(3)
	ground := hittable.NewSphere(vec.New(0, -1000, 0), 1000, hittable.NewLambertian(vec.New(.4, .4, .4)))
	world.Add(ground)
	// Load the model
	// Check the objectLoader file to find additional options for loading the model such as pre-positioning, and finding dielectric materials for sampling
	opt := objLoader.DefaultLoadOptions()
	opt.ScaleFactor = 5 // Scale the model up or down in size
	opt.Center = true
	opt.Position = vec.New(0, 1.8, 0)                                                  // hint: usee  the debug output to find the minimum y-value in the model
	opt.Debug = true                                                                   // change this to true to see information about the model as it's being loaded.
	opt.DefaultMaterial = hittable.NewMetal(vec.New(255.0/255.0, 215.0/255.0, 0), 0.5) // Solid gold dragon statue
	model, lights := objLoader.LoadObjWithOptions("dragon.obj", opt)
	world.Add(hittable.RotateY(model, 180))

	// Add a separate light source to the scene. I think of this as a "sun"
	light := hittable.NewSphere(vec.New(7, 13, 7), 5, hittable.NewDiffuseLight(vec.New(4, 4, 4)))
	world.Add(light)
	if hl, ok := lights.(*hittable.HittableList); ok {
		hl.Add(light)
		lights = hl
	}

	cam.AspectRatio = 16.0 / 9.0
	cam.Width = 600

	cam.SamplesPerPixel = 250
	cam.MaxDepth = 50

	cam.Background = vec.New(0, 0, 0)

	cam.VerticalFOV = 40
	cam.MaxContribution = 2.0
	cam.PositionCamera(vec.New(10, 5, 10), vec.New(0, 0, 0), vec.New(0, 1, 0))

	cam.DefocusAngle = .1

	cam.Render(world, lights)
}

// The default scene that will render when no scene is specified.
func defaultScene(c *camera.Camera) {

}

func main() {
	// parse CLI arguments
	cpuprofile := flag.String("cpuprofile", "", "Write cpu profile to file")
	outFile := flag.String("o", "image.ppm", "Specify a custom output file")
	coreCount := flag.Int("N", 1, "Set the number of cores to allocate to rendering")
	scene := flag.Int("S", -1, "Set the scene to render, default will render a custom scene function")

	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// Attempt to create the output file.
	file, err := os.Create(*outFile)
	if err != nil {
		log.Fatal("Error creating output file\n")
	}
	defer file.Close()

	// Initialize an output buffer.
	outBuf := bytes.Buffer{}

	// Initialize the camera.
	c := camera.Camera{}
	c.Out = &outBuf
	c.MaxThreads = *coreCount

	switch *scene {
	case 1:
		book1Scene(&c)
		break
	case 2:
		book2Scene(&c)
		break
	case 3:
		book3Scene(&c)
		break
	case 4:
		simpleLight(&c)
		break
	case 5:
		quads(&c)
		break
	case 6:
		cornellBox(&c)
		break
	case 7:
		cornellSmoke(&c)
		break
	case 8:
		modelExample(&c)
		break
	default:
		defaultScene(&c)
	}

	// Write the image to the output file.
	file.Write(outBuf.Bytes())
}
