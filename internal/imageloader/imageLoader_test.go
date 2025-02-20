package ImageLoader_test

import (
	"testing"

	ImageLoader "github.com/nsp5488/go_raytracer/internal/imageloader"
)

func TestReadPNG(t *testing.T) {
	// test.png was generated as a simple 5x5 image in ppm format then converted to png

	rti := ImageLoader.LoadImage("test.png")
	exp := "5, 5, png, 25"
	act := rti.String()
	if exp != act {
		t.Errorf("Expected %s, but got %s", exp, act)
	}
	testImageData(rti, IMG_DATA, t)

}
func TestReadJPG(t *testing.T) {
	// test.jpg was generated as a simple 5x5 image in ppm format then converted to JPEG
	rti := ImageLoader.LoadImage("test.jpg")
	exp := "5, 5, jpeg, 25"
	act := rti.String()
	if exp != act {
		t.Errorf("Expected %s, but got %s", exp, act)
	}
	testImageData(rti, LOSSY_IMG_DATA, t)
}

// Slightly lossy img data from conversion to JPEG
var LOSSY_IMG_DATA = [25][3]uint8{
	{216, 226, 255},
	{160, 170, 208},
	{133, 143, 180},
	{132, 142, 179},
	{226, 237, 255},
	{114, 123, 162},
	{60, 69, 108},
	{85, 95, 131},
	{47, 57, 93},
	{136, 147, 177},
	{90, 99, 138},
	{0, 9, 48},
	{26, 37, 71},
	{76, 86, 121},
	{99, 110, 140},
	{130, 140, 178},
	{3, 13, 52},
	{9, 20, 54},
	{124, 134, 169},
	{147, 158, 188},
	{218, 228, 255},
	{108, 117, 156},
	{102, 112, 147},
	{134, 144, 179},
	{222, 234, 255},
}

// Accurate image data from a png conversion
var IMG_DATA = [25][3]uint8{
	{209, 226, 249},
	{161, 176, 189},
	{126, 146, 139},
	{146, 162, 191},
	{214, 230, 254},
	{124, 132, 172},
	{53, 64, 122},
	{63, 80, 105},
	{26, 34, 112},
	{154, 165, 198},
	{83, 88, 143},
	{4, 13, 116},
	{19, 35, 120},
	{93, 119, 64},
	{94, 108, 131},
	{138, 146, 181},
	{0, 0, 113},
	{0, 12, 114},
	{110, 122, 112},
	{140, 149, 164},
	{220, 231, 249},
	{111, 122, 166},
	{109, 122, 166},
	{140, 152, 175},
	{215, 227, 246},
}

func testImageData(rti *ImageLoader.RTImage, expData [25][3]uint8, t *testing.T) {
	expHeight := 5
	expWidth := 5
	template := "Expected %d, but got %d"
	if expHeight != rti.Height {
		t.Errorf(template, expHeight, rti.Height)
	}
	if expWidth != rti.Width {
		t.Errorf(template, expWidth, rti.Width)
	}
	for x := 0; x < expWidth; x++ {
		for y := 0; y < expHeight; y++ {
			x, y = min(max(x, 0), rti.Width), min(max(y, 0), rti.Height)
			idx := y*rti.Width + x
			actual := rti.PixelData(x, y)
			for i := 0; i < 3; i++ {
				if expData[idx][i] != actual.Data[i] {
					t.Errorf(template, expData[idx][i], actual.Data[i])
				}
			}
		}
	}
}
