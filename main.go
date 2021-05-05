package main

import (
	"math/rand"
	"time"

	. "github.com/arata-nvm/nakitu/nakitu"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// image
	aspectRatio := 16.0 / 9.0
	imageWidth := 400

	// camera
	lookFrom := NewVec3(13, 2, 3)
	lookAt := NewVec3(0, 0, 0)

	vFOV := 20.0
	vUp := NewVec3(0, 1, 0)
	distToFocus := 10.0
	aperture := 0.0

	// scene
	background := NewVec3(0.7, 0.8, 1.0)
	samplesPerPixel := 100
	maxDepth := 10

	// world
	var world Hittable
	switch 4 {
	case 0:
		world = randomScene()
	case 1:
		world = earth()
	case 2:
		world = simpleLight()
		background = NewVec3(0, 0, 0)
		lookFrom = NewVec3(26, 3, 6)
		lookAt = NewVec3(0, 2, 0)
	case 3:
		world = cornellBox()
		aspectRatio = 1.0
		imageWidth = 600
		background = NewVec3(0, 0, 0)
		lookFrom = NewVec3(278, 278, -800)
		lookAt = NewVec3(278, 278, 0)
		vFOV = 40.0
	case 4:
		world = finalScene()
		aspectRatio = 1.0
		imageWidth = 800
		background = NewVec3(0, 0, 0)
		lookFrom = NewVec3(478, 278, -600)
		lookAt = NewVec3(278, 278, 0)
		vFOV = 40.0
	}

	// camera
	camera := NewCamera(
		lookFrom,
		lookAt,
		vUp,
		vFOV,
		aspectRatio,
		aperture,
		distToFocus,
	)
	camera.Time1 = 1.0

	// render
	imageHeight := int(float64(float64(imageWidth)) / aspectRatio)
	scene := NewScene(imageWidth, imageHeight, world, camera)
	scene.SamplesPerPixel = samplesPerPixel
	scene.MaxDepth = maxDepth
	scene.Background = background
	scene.RenderParallel(8)
	scene.WriteToFile("image.ppm")
}

func randomScene() Hittable {
	world := NewHittableList()

	checker := NewCheckerTexture(
		NewSolidColor(0.2, 0.3, 0.1),
		NewSolidColor(0.9, 0.9, 0.9),
	)
	matGround := NewLambertian(checker)
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
					mat := NewLambertian(NewSolidColor(albedo.X(), albedo.Y(), albedo.Z()))
					center2 := center.Add(NewVec3(0, Random(0, 0.5), 0))
					world.Add(NewMovingSphere(center, center2, 0.0, 1.0, 0.2, mat))
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

	mat2 := NewLambertian(NewSolidColor(0.4, 0.2, 0.1))
	world.Add(NewSphere(NewVec3(-4, 1, 0), 1.0, mat2))

	mat3 := NewMetal(NewVec3(0.7, 0.6, 0.5), 0)
	world.Add(NewSphere(NewVec3(4, 1, 0), 1.0, mat3))

	return world
}

func earth() Hittable {
	world := NewHittableList()

	earthTexture := NewImageTexture("earthmap.jpg")
	earthSurface := NewLambertian(earthTexture)
	globe := NewSphere(NewVec3(0, 0, 0), 2, earthSurface)
	world.Add(globe)

	return world
}

func simpleLight() Hittable {
	world := NewHittableList()

	tex := NewCheckerTexture(
		NewSolidColor(0.2, 0.3, 0.1),
		NewSolidColor(0.9, 0.9, 0.9),
	)
	world.Add(NewSphere(NewVec3(0, -1000, 0), 1000, NewLambertian(tex)))
	world.Add(NewSphere(NewVec3(0, 2, 0), 2, NewLambertian(tex)))

	difflight := NewDiffuseLight(NewSolidColor(4, 4, 4))
	world.Add(NewXYRect(3, 1, 5, 3, -2, difflight))

	return world
}

func cornellBox() Hittable {
	world := NewHittableList()

	red := NewLambertian(NewSolidColor(0.65, 0.05, 0.05))
	white := NewLambertian(NewSolidColor(0.73, 0.73, 0.73))
	green := NewLambertian(NewSolidColor(0.12, 0.45, 0.15))
	light := NewDiffuseLight(NewSolidColor(15, 15, 15))

	world.Add(NewYZRect(0, 0, 555, 555, 555, green))
	world.Add(NewYZRect(0, 0, 555, 555, 0, red))
	world.Add(NewXZRect(213, 227, 343, 332, 554, light))
	world.Add(NewXZRect(0, 0, 555, 555, 0, white))
	world.Add(NewXZRect(0, 0, 555, 555, 555, white))
	world.Add(NewXYRect(0, 0, 555, 555, 555, white))

	var box1 Hittable
	box1 = NewBox(NewVec3(0, 0, 0), NewVec3(165, 330, 165), white)
	box1 = NewRotateY(box1, 15)
	box1 = NewTranslate(box1, NewVec3(265, 0, 295))

	var box2 Hittable
	box2 = NewBox(NewVec3(0, 0, 0), NewVec3(165, 165, 165), white)
	box2 = NewRotateY(box2, -18)
	box2 = NewTranslate(box2, NewVec3(130, 0, 65))

	world.Add(NewConstantMedium(box1, 0.01, NewSolidColor(0, 0, 0)))
	world.Add(NewConstantMedium(box2, 0.01, NewSolidColor(1, 1, 1)))

	return world
}

func finalScene() Hittable {
	boxes1 := NewHittableList()

	matGround := NewLambertian(NewSolidColor(0.48, 0.83, 0.53))

	boxesPerSide := 20
	for i := 0; i < boxesPerSide; i++ {
		for j := 0; j < boxesPerSide; j++ {
			w := 100.0
			x0 := -1000.0 + float64(i)*w
			x1 := x0 + w

			z0 := -1000.0 + float64(j)*w
			z1 := z0 + w

			y0 := 0.0
			y1 := Random(0, 101)

			boxes1.Add(NewBox(NewVec3(x0, y0, z0), NewVec3(x1, y1, z1), matGround))
		}
	}

	world := NewHittableList()

	world.Add(NewBVHNode(boxes1, 0, 1))

	light := NewDiffuseLight(NewSolidColor(7, 7, 7))
	world.Add(NewXZRect(123, 147, 423, 412, 554, light))

	center1 := NewVec3(400, 400, 200)
	center2 := center1.Add(NewVec3(30, 0, 0))
	matMovingSphere := NewLambertian(NewSolidColor(0.7, 0.3, 0.1))
	world.Add(NewMovingSphere(center1, center2, 0, 1, 50, matMovingSphere))

	world.Add(NewSphere(NewVec3(260, 150, 45), 50, NewDielectric(1.5)))
	world.Add(NewSphere(NewVec3(0, 150, 145), 50, NewMetal(NewVec3(0.8, 0.8, 0.9), 1.0)))

	boundary := NewSphere(NewVec3(360, 150, 145), 70, NewDielectric(1.5))
	world.Add(boundary)
	world.Add(NewConstantMedium(boundary, 0.2, NewSolidColor(0.2, 0.4, 0.9)))
	boundary = NewSphere(NewVec3(0, 0, 0), 5000, NewDielectric(1.5))
	world.Add(NewConstantMedium(boundary, 0.0001, NewSolidColor(1, 1, 1)))

	texEarth := NewImageTexture("earthmap.jpg")
	world.Add(NewSphere(NewVec3(400, 200, 400), 100, NewLambertian(texEarth)))
	// noise

	boxes2 := NewHittableList()
	white := NewLambertian(NewSolidColor(0.73, 0.73, 0.73))
	ns := 1000
	for j := 0; j < ns; j++ {
		boxes2.Add(NewSphere(RandomVec3In(0, 165), 10, white))
	}

	world.Add(NewTranslate(NewRotateY(NewBVHNode(boxes2, 0, 1), 15), NewVec3(-100, 270, 395)))

	return world
}
