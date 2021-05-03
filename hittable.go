package main

import "math"

type HitRecord struct {
	Point     Point3
	Normal    Vec3
	Mat       Material
	T         float64
	frontFace bool
}

func (hr *HitRecord) SetFaceNormal(r *Ray, outwardNormal Vec3) {
	hr.frontFace = r.Dir.Dot(outwardNormal) < 0
	if hr.frontFace {
		hr.Normal = outwardNormal
	} else {
		hr.Normal = outwardNormal.Neg()
	}
}

type HittableList struct {
	Objects []Hittable
}

func NewHittableList() *HittableList {
	return &HittableList{}
}

func (hl *HittableList) Clear() {
	hl.Objects = hl.Objects[:0]
}

func (hl *HittableList) Add(object Hittable) {
	hl.Objects = append(hl.Objects, object)
}

func (hl *HittableList) Hit(r *Ray, tMin, tMax float64, rec *HitRecord) bool {
	var tempRec HitRecord
	hitAnything := false
	closestSoFar := tMax

	for _, object := range hl.Objects {
		if object.Hit(r, tMin, closestSoFar, &tempRec) {
			hitAnything = true
			closestSoFar = tempRec.T
			*rec = tempRec
		}
	}

	return hitAnything
}

type Hittable interface {
	Hit(r *Ray, tMin, tMax float64, rec *HitRecord) bool
}

type Sphere struct {
	Center Point3
	Radius float64
	Mat    Material
}

func NewSphere(center Point3, radius float64, mat Material) *Sphere {
	return &Sphere{Center: center, Radius: radius, Mat: mat}
}

func (s *Sphere) Hit(r *Ray, tMin, tMax float64, rec *HitRecord) bool {
	oc := r.Origin.Sub(s.Center)
	a := r.Dir.LenSquared()
	halfB := oc.Dot(r.Dir)
	c := oc.LenSquared() - s.Radius*s.Radius

	d := halfB*halfB - a*c
	if d < 0 {
		return false
	}

	sqrtD := math.Sqrt(d)
	root := (-halfB - sqrtD) / a
	if root < tMin || tMax < root {
		root = (-halfB + sqrtD) / a

		if root < tMin || tMax < root {
			return false
		}
	}

	rec.T = root
	rec.Point = r.At(rec.T)
	outwardNormal := rec.Point.Sub(s.Center).Divf(s.Radius)
	rec.SetFaceNormal(r, outwardNormal)
	rec.Mat = s.Mat

	return true
}
