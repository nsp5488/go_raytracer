package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/nsp5488/go_raytracer/internal/progress"
	"github.com/nsp5488/go_raytracer/internal/ray"
	"github.com/nsp5488/go_raytracer/internal/vec"
)

func main() {
	file, err := os.Create("image.ppm")
	if err != nil {
		log.Fatal("Error creating output file\n")
	}
	aspectRatio := float64(16) / float64(9)
	imageWidth := 400
	// calculate image height given aspect ratio, clamped to >=1
	imageHeight := max(1, int(float64(imageWidth)/aspectRatio))

	// define camera information
	focalLength := 1.0
	viewportHeight := 2.0
	viewportWidth := viewportHeight * (float64(imageWidth) / float64(imageHeight))
	cameraCenter := vec.Empty()

	// Calculate the vectors across the horizontal and down the vertical viewport edges.
	viewportU := vec.New(viewportWidth, 0, 0)
	viewportV := vec.New(0, -viewportHeight, 0)

	// Calculate the horizontal and vertical delta vectors from pixel to pixel.
	pixelDeltaU := viewportU.Scale(1.0 / float64(imageWidth))
	pixelDeltaV := viewportV.Scale(1.0 / float64(imageHeight))

	// Calculate the location of the upper left pixel.
	viewportTopLeft := cameraCenter.Add(vec.New(0, 0, focalLength).Negate()).Add(viewportU.Scale(0.5).Negate()).Add(viewportV.Scale(0.5).Negate())
	pixel00Loc := viewportTopLeft.Add(pixelDeltaU.Add(pixelDeltaV).Scale(0.5))

	io.WriteString(file, fmt.Sprintf("P3\n%d %d\n255\n", imageWidth, imageHeight))
	p := progress.InitBar(imageHeight)

	// Run the processing in a separate goroutine
	go func() {
		for i := range imageHeight {
			for j := range imageWidth {
				pixelCenter := pixel00Loc.Add(pixelDeltaV.Scale(float64(i))).Add(pixelDeltaU.Scale(float64(j)))
				rayDirection := pixelCenter.Add(cameraCenter.Negate())
				r := ray.New(cameraCenter, rayDirection)

				pixelColor := r.Color()

				pixelColor.PrintColor(file)
			}
			p.Send(i + 1)
		}
	}()

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		return
	}
}
