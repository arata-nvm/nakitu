package nakitu

import (
	"image"
	"image/jpeg"
	"log"
	"math"
	"os"
)

type Texture interface {
	Value(u, v float64, p Point3) Color
}

type SolidColor struct {
	ColorValue Color
}

func NewSolidColor(r, g, b float64) *SolidColor {
	return &SolidColor{
		ColorValue: NewVec3(r, g, b),
	}
}

func (s *SolidColor) Value(u, v float64, p Point3) Color {
	return s.ColorValue
}

type CheckerTexture struct {
	Odd  Texture
	Even Texture
}

func NewCheckerTexture(odd Texture, even Texture) *CheckerTexture {
	return &CheckerTexture{
		Odd:  odd,
		Even: even,
	}
}

func (c *CheckerTexture) Value(u, v float64, p Point3) Color {
	sines := math.Sin(10*p.X()) * math.Sin(10*p.Y()) * math.Sin(10*p.Z())
	if sines < 0 {
		return c.Odd.Value(u, v, p)
	} else {
		return c.Even.Value(u, v, p)
	}
}

type ImageTexture struct {
	Image  image.Image
	Width  int
	Height int
}

func NewImageTexture(name string) *ImageTexture {
	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}

	img, err := jpeg.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	return &ImageTexture{
		Image:  img,
		Width:  img.Bounds().Dx(),
		Height: img.Bounds().Dy(),
	}
}

func (t *ImageTexture) Value(u, v float64, p Point3) Color {
	u = Clamp(u, 0, 1)
	v = 1 - Clamp(v, 0, 1)

	i := int(u * float64(t.Width))
	j := int(v * float64(t.Height))

	if i >= t.Width {
		i = t.Width - 1
	}
	if j >= t.Height {
		j = t.Height - 1
	}

	c := t.Image.At(i, j)
	r, g, b, _ := c.RGBA()

	colorScale := 1.0 / 0xffff
	return NewVec3(
		float64(r)*colorScale,
		float64(g)*colorScale,
		float64(b)*colorScale,
	)
}
