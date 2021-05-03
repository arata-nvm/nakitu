package main

type Ray struct {
	Origin Point3
	Dir    Vec3
}

func NewRay(origin Point3, dir Vec3) *Ray {
	return &Ray{
		Origin: origin,
		Dir:    dir,
	}
}

func (r *Ray) At(t float64) Point3 {
	return r.Origin.Add(r.Dir.Mulf(t))
}
