package nakitu

import (
	"bufio"
	"math"
	"math/rand"
	"os"
	"sync"

	pb "github.com/cheggaaa/pb/v3"
)

type Scene struct {
	Width           int
	Height          int
	Background      Color
	SamplesPerPixel int
	MaxDepth        int
	World           Hittable
	Camera          *Camera

	Output *Image
}

func NewScene(width, height int, world Hittable, camera *Camera) *Scene {
	return &Scene{
		Width:           width,
		Height:          height,
		Background:      NewVec3(0, 0, 0),
		SamplesPerPixel: 10,
		MaxDepth:        6,
		World:           world,
		Camera:          camera,
		Output:          NewImage(width, height),
	}
}

func (s *Scene) WriteToFile(name string) {
	f, _ := os.Create(name)
	buf := bufio.NewWriter(f)
	s.Output.Write(buf)
	buf.Flush()
	f.Close()
}

func (s *Scene) Render() {
	bar := pb.StartNew(s.Height)
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			s.RenderPixel(x, y)
		}
		bar.Increment()
	}
	bar.Finish()
}

func (s *Scene) RenderParallel(numOfCore int) {
	lines := make(chan int)
	go func() {
		for y := 0; y < s.Height; y++ {
			lines <- y
		}
		close(lines)
	}()

	wg := sync.WaitGroup{}
	bar := pb.StartNew(s.Height)

	for i := 0; i < numOfCore; i++ {
		wg.Add(1)
		go func() {
			for y := range lines {
				for x := 0; x < s.Width; x++ {
					s.RenderPixel(x, y)
				}
				bar.Increment()
			}
			wg.Done()
		}()
	}

	wg.Wait()
	bar.Finish()
}

func (s *Scene) RenderPixel(x, y int) {
	sumColor := Zero()

	for i := 0; i < s.SamplesPerPixel; i++ {
		u := (float64(x) + rand.Float64()) / float64(s.Width-1)
		v := (float64(y) + rand.Float64()) / float64(s.Height-1)
		r := s.Camera.GetRay(u, v)
		color := rayColor(r, s.Background, s.World, s.MaxDepth)
		sumColor = sumColor.Add(color)
	}

	rgb := toRGB(sumColor, s.SamplesPerPixel)
	s.Output.SetPixel(s.Width-x-1, s.Height-y-1, rgb)
}

func rayColor(r *Ray, background Color, world Hittable, depth int) Color {
	var rec HitRecord

	if depth <= 0 {
		return Zero()
	}

	if !world.Hit(r, 0.001, math.Inf(1), &rec) {
		return background
	}

	var scattered Ray
	var attenuation Color
	emitted := rec.Mat.Emitted(rec.U, rec.V, rec.Point)

	if !rec.Mat.Scatter(r, &rec, &attenuation, &scattered) {
		return emitted
	}

	c := rayColor(&scattered, background, world, depth-1)
	return emitted.Add(attenuation.Mul(c))
}

func toRGB(color Vec3, samplesPerPixel int) RGB {
	toColor := func(x float64) int {
		scale := 1.0 / float64(samplesPerPixel)
		scaledX := math.Sqrt(x * scale)
		return int(256 * Clamp(scaledX, 0, 0.999))
	}

	return NewRGB(
		toColor(color.X()),
		toColor(color.Y()),
		toColor(color.Z()),
	)
}
