package ImageLoader

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

// Holds the RGB values of a pixel in 3 [0-255] values
type pixel struct {
	Data [3]uint8
}

// Raytrace Image, used to import images as textures
type RTImage struct {
	img    image.Image
	format string
	Width  int
	Height int

	bdata []*pixel
}

// Loads the specified filename as an RTImage
func LoadImage(filename string) *RTImage {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Could not open %s", filename)
	}
	img, format, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Error while decoding %s: %s", filename, err.Error())
	}
	i := &RTImage{img: img, format: format}
	i.Width = i.img.Bounds().Dx()
	i.Height = i.img.Bounds().Dy()

	i.bdata = make([]*pixel, i.Width*i.Height, i.Width*i.Height)
	i.convertToBytes()

	return i
}

// Magenta to return whenever there's no valid data
var magenta = &pixel{[3]uint8{255, 0, 255}}

// Returns the pixel data for the wrapped image at (x,y)
func (rti *RTImage) PixelData(x, y int) *pixel {
	x, y = min(max(x, 0), rti.Width), min(max(y, 0), rti.Height)
	idx := y*rti.Width + x

	if idx >= len(rti.bdata) || rti.bdata[idx] == nil {
		fmt.Println(idx)
		return magenta
	}
	return rti.bdata[idx]

}

// Scales an RGBA uint16 to a uint8 by right shifting 8 times.
// NOTE this only works because Go's RGBA color.Color stores uint16 as uint32,
// otherwise we would need to shift 24 bits to the right.
func scaleUint32ToUint8(v uint32) uint8 {
	return uint8(v >> 8)
}

// Converts a color.Color pixel to an internal slice of uint8
func (rti *RTImage) rgbToByte(pData color.Color, idx int) {
	r, g, b, _ := pData.RGBA()
	rti.bdata[idx] = &pixel{Data: [3]uint8{scaleUint32ToUint8(r), scaleUint32ToUint8(g), scaleUint32ToUint8(b)}}
}

// Private helper method to build the internal representation of pixels in the image.
func (rti *RTImage) convertToBytes() {
	for x := 0; x < rti.Width; x++ {
		for y := 0; y < rti.Height; y++ {
			rti.rgbToByte(rti.img.At(x, y), y*rti.Width+x)
		}
	}
}

func (rti *RTImage) String() string {
	return fmt.Sprintf("%d, %d, %s, %d", rti.Width, rti.Height, rti.format, len(rti.bdata))
}
