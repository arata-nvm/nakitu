package nakitu

import (
	"math"
	"math/rand"
)

type Material interface {
	Scatter(rIn *Ray, rec *HitRecord, attenuation *Color, scattered *Ray) bool
	Emitted(u, v float64, p Point3) Color
}

type DefaultEmitter struct{}

func (d *DefaultEmitter) Emitted(u, v float64, p Point3) Color {
	return NewVec3(0, 0, 0)
}

type Lambertian struct {
	Albedo Texture
	DefaultEmitter
}

func NewLambertian(a Texture) *Lambertian {
	return &Lambertian{Albedo: a}
}

func (l *Lambertian) Scatter(rIn *Ray, rec *HitRecord, attenuation *Color, scattered *Ray) bool {
	scatterDir := rec.Normal.Add(RandomUnitVector())

	if scatterDir.NearZero() {
		scatterDir = rec.Normal
	}

	*scattered = *NewRay(rec.Point, scatterDir, rIn.Time)
	*attenuation = l.Albedo.Value(rec.U, rec.V, rec.Point)
	return true
}

type Metal struct {
	Albedo Color
	Fuzz   float64
	DefaultEmitter
}

func NewMetal(a Color, f float64) *Metal {
	return &Metal{Albedo: a, Fuzz: f}
}

func (m *Metal) Scatter(rIn *Ray, rec *HitRecord, attenuation *Color, scattered *Ray) bool {
	reflected := rIn.Dir.Unit().Reflect(rec.Normal)
	*scattered = *NewRay(rec.Point, reflected.Add(RandomInUnitSphere().Mulf(m.Fuzz)), rIn.Time)
	*attenuation = m.Albedo
	return scattered.Dir.Dot(rec.Normal) > 0
}

type Dielectric struct {
	Ir float64
	DefaultEmitter
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

	*scattered = *NewRay(rec.Point, dir, rIn.Time)
	return true
}

func reflectance(cosine, refIdx float64) float64 {
	r0 := (1 - refIdx) / (1 + refIdx)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}

type DiffuseLight struct {
	Emit Texture
}

func NewDiffuseLight(t Texture) *DiffuseLight {
	return &DiffuseLight{
		Emit: t,
	}
}

func (d *DiffuseLight) Scatter(rIn *Ray, rec *HitRecord, attenuation *Color, scattered *Ray) bool {
	return false
}

func (d *DiffuseLight) Emitted(u, v float64, p Point3) Color {
	return d.Emit.Value(u, v, p)
}

type Isotropic struct {
	Albedo Texture
	DefaultEmitter
}

func NewIsotropic(tex Texture) *Isotropic {
	return &Isotropic{
		Albedo: tex,
	}
}

func (i *Isotropic) Scatter(rIn *Ray, rec *HitRecord, attenuation *Color, scattered *Ray) bool {
	*scattered = *NewRay(rec.Point, RandomInUnitSphere(), rIn.Time)
	*attenuation = i.Albedo.Value(rec.U, rec.V, rec.Point)
	return true
}
