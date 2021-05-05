package nakitu

type Ray struct {
	Origin Point3
	Dir    Vec3
	Time   float64
}

func NewRay(origin Point3, dir Vec3, time float64) *Ray {
	return &Ray{
		Origin: origin,
		Dir:    dir,
		Time:   time,
	}
}

func (r *Ray) At(t float64) Point3 {
	return r.Origin.Add(r.Dir.Mulf(t))
}
