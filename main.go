package main

import (
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// image
	aspectRatio := 3.0 / 2.0
	imageWidth := 400
	imageHeight := int(float64(float64(imageWidth)) / aspectRatio)

	// world
	world := randomScene()

	// camera
	lookFrom := NewVec3(13, 2, 3)
	lookAt := NewVec3(0, 0, 0)
	vUp := NewVec3(0, 1, 0)
	distToFocus := 10.0
	aperture := 0.1
	camera := NewCamera(
		lookFrom,
		lookAt,
		vUp,
		20,
		aspectRatio,
		aperture,
		distToFocus,
	)

	// render
	scene := NewScene(imageWidth, imageHeight, world, camera)
	scene.SamplesPerPixel = 1
	scene.MaxDepth = 6
	scene.RenderParallel(8)
	scene.WriteToFile("image.ppm")
}

func randomScene() Hittable {
	world := NewHittableList()

	matGround := NewLambertian(NewVec3(0.5, 0.5, 0.5))
	world.Add(NewSphere(NewVec3(0, -1000, 0), 1000, matGround))

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			r := rand.Float64()
			center := NewVec3(
				float64(a)+0.9*rand.Float64(),
				0.2,
				float64(b)+0.9*rand.Float64(),
			)

			if center.Sub(NewVec3(4, 0.2, 0)).Len() > 0.9 {
				switch {
				case r < 0.8:
					albedo := RandomVec3().Mul(RandomVec3())
					mat := NewLambertian(albedo)
					world.Add(NewSphere(center, 0.2, mat))
				case r < 0.95:
					albedo := RandomVec3In(0.5, 1)
					fuzz := Random(0, 0.5)
					mat := NewMetal(albedo, fuzz)
					world.Add(NewSphere(center, 0.2, mat))
				default:
					mat := NewDielectric(1.5)
					world.Add(NewSphere(center, 0.2, mat))
				}

			}
		}
	}

	mat1 := NewDielectric(1.5)
	world.Add(NewSphere(NewVec3(0, 1, 0), 1, mat1))

	mat2 := NewLambertian(NewVec3(0.4, 0.2, 0.1))
	world.Add(NewSphere(NewVec3(-4, 1, 0), 1.0, mat2))

	mat3 := NewMetal(NewVec3(0.7, 0.6, 0.5), 0)
	world.Add(NewSphere(NewVec3(4, 1, 0), 1.0, mat3))

	return world
}
