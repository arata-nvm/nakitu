package main

import (
	"math"
	"math/rand"
)

type Material interface {
	Scatter(rIn *Ray, rec *HitRecord, attenuation *Color, scattered *Ray) bool
}

type Lambertian struct {
	Albedo Color
}

func NewLambertian(a Color) *Lambertian {
	return &Lambertian{Albedo: a}
}

func (l *Lambertian) Scatter(rIn *Ray, rec *HitRecord, attenuation *Color, scattered *Ray) bool {
	scatterDir := rec.Normal.Add(RandomUnitVector())

	if scatterDir.NearZero() {
		scatterDir = rec.Normal
	}

	*scattered = *NewRay(rec.Point, scatterDir)
	*attenuation = l.Albedo
	return true
}

type Metal struct {
	Albedo Color
	Fuzz   float64
}

func NewMetal(a Color, f float64) *Metal {
	return &Metal{Albedo: a, Fuzz: f}
}

func (m *Metal) Scatter(rIn *Ray, rec *HitRecord, attenuation *Color, scattered *Ray) bool {
	reflected := rIn.Dir.Unit().Reflect(rec.Normal)
	*scattered = *NewRay(rec.Point, reflected.Add(RandomInUnitSphere().Mulf(m.Fuzz)))
	*attenuation = m.Albedo
	return scattered.Dir.Dot(rec.Normal) > 0
}

type Dielectric struct {
	Ir float64
}

func NewDielectric(indexOfRefraction float64) *Dielectric {
	return &Dielectric{Ir: indexOfRefraction}
}

func (d *Dielectric) Scatter(rIn *Ray, rec *HitRecord, attenuation *Color, scattered *Ray) bool {
	*attenuation = NewVec3(1, 1, 1)
	var refractionRatio float64
	if rec.frontFace {
		refractionRatio = 1.0 / d.Ir
	} else {
		refractionRatio = d.Ir
	}

	unitDir := rIn.Dir.Unit()
	cosTheta := math.Min(unitDir.Neg().Dot(rec.Normal), 1)
	sinTheta := math.Sqrt(1 - cosTheta*cosTheta)

	cannotRefract := refractionRatio*sinTheta > 1
	var dir Vec3

	if cannotRefract || reflectance(cosTheta, refractionRatio) > rand.Float64() {
		dir = unitDir.Reflect(rec.Normal)
	} else {
		dir = unitDir.Refract(rec.Normal, refractionRatio)
	}

	*scattered = *NewRay(rec.Point, dir)
	return true
}

func reflectance(cosine, refIdx float64) float64 {
	r0 := (1 - refIdx) / (1 + refIdx)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}
