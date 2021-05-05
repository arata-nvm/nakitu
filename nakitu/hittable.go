package nakitu

import (
	"math"
	"math/rand"
)

type HitRecord struct {
	Point     Point3
	Normal    Vec3
	Mat       Material
	T         float64
	U         float64
	V         float64
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

func (hl *HittableList) BoundingBox(time0, time1 float64, outputBox *AABB) bool {
	if len(hl.Objects) == 0 {
		return false
	}

	var tempBox AABB
	for _, object := range hl.Objects {
		if !object.BoundingBox(time0, time1, &tempBox) {
			return false
		}
		outputBox = SurroundingBox(outputBox, &tempBox)
	}

	return true
}

type Hittable interface {
	Hit(r *Ray, tMin, tMax float64, rec *HitRecord) bool
	BoundingBox(time0, time1 float64, outputBox *AABB) bool
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
	getSphereUV(outwardNormal, &rec.U, &rec.V)
	rec.Mat = s.Mat

	return true
}

func (s *Sphere) BoundingBox(time0, time1 float64, outputBox *AABB) bool {
	*outputBox = *NewAABB(
		s.Center.Sub(NewVec3(s.Radius, s.Radius, s.Radius)),
		s.Center.Add(NewVec3(s.Radius, s.Radius, s.Radius)),
	)
	return true
}

func getSphereUV(p Point3, u, v *float64) {
	theta := math.Acos(-p.Y())
	phi := math.Atan2(-p.Z(), p.X()) + math.Pi

	*u = phi / (2 * math.Pi)
	*v = theta / math.Pi
}

type MovingSphere struct {
	Center0 Point3
	Center1 Point3
	Time0   float64
	Time1   float64
	Radius  float64
	Mat     Material
}

func NewMovingSphere(center0, center1 Point3, time0, time1 float64, radius float64, mat Material) *MovingSphere {
	return &MovingSphere{
		Center0: center0,
		Center1: center1,
		Time0:   time0,
		Time1:   time1,
		Radius:  radius,
		Mat:     mat,
	}
}

func (s *MovingSphere) Hit(r *Ray, tMin, tMax float64, rec *HitRecord) bool {
	oc := r.Origin.Sub(s.Center(r.Time))
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
	outwardNormal := rec.Point.Sub(s.Center(r.Time)).Divf(s.Radius)
	rec.SetFaceNormal(r, outwardNormal)
	rec.Mat = s.Mat
	return true
}

func (s *MovingSphere) BoundingBox(time0, time1 float64, outputBox *AABB) bool {
	box0 := NewAABB(
		s.Center(time0).Sub(NewVec3(s.Radius, s.Radius, s.Radius)),
		s.Center(time0).Add(NewVec3(s.Radius, s.Radius, s.Radius)),
	)
	box1 := NewAABB(
		s.Center(time1).Sub(NewVec3(s.Radius, s.Radius, s.Radius)),
		s.Center(time1).Add(NewVec3(s.Radius, s.Radius, s.Radius)),
	)
	*outputBox = *SurroundingBox(box0, box1)
	return true
}

func (s *MovingSphere) Center(time float64) Point3 {
	return s.Center0.Add(s.Center1.Sub(s.Center0).Mulf((time - s.Time0) / (s.Time1 - s.Time0)))
}

type XYRect struct {
	X0, Y0 float64
	X1, Y1 float64
	K      float64
	Mat    Material
}

func NewXYRect(x0, y0, x1, y1, k float64, mat Material) *XYRect {
	return &XYRect{
		X0:  x0,
		Y0:  y0,
		X1:  x1,
		Y1:  y1,
		K:   k,
		Mat: mat,
	}
}

func (s *XYRect) Hit(r *Ray, tMin, tMax float64, rec *HitRecord) bool {
	t := (s.K - r.Origin.Z()) / r.Dir.Z()
	if t < tMin || t > tMax {
		return false
	}
	x := r.Origin.X() + t*r.Dir.X()
	y := r.Origin.Y() + t*r.Dir.Y()
	if x < s.X0 || x > s.X1 || y < s.Y0 || y > s.Y1 {
		return false
	}
	rec.U = (x - s.X0) / (s.X1 - s.X0)
	rec.V = (y - s.Y0) / (s.Y1 - s.Y0)
	rec.T = t
	outwardNormal := NewVec3(0, 0, 1)
	rec.SetFaceNormal(r, outwardNormal)
	rec.Mat = s.Mat
	rec.Point = r.At(t)
	return true
}

func (s *XYRect) BoundingBox(time0, time1 float64, outputBox *AABB) bool {
	*outputBox = *NewAABB(
		NewVec3(s.X0, s.Y0, s.K-0.0001),
		NewVec3(s.X1, s.Y1, s.K+0.0001),
	)
	return true
}

type XZRect struct {
	X0, Z0 float64
	X1, Z1 float64
	K      float64
	Mat    Material
}

func NewXZRect(x0, z0, x1, z1, k float64, mat Material) *XZRect {
	return &XZRect{
		X0:  x0,
		Z0:  z0,
		X1:  x1,
		Z1:  z1,
		K:   k,
		Mat: mat,
	}
}

func (s *XZRect) Hit(r *Ray, tMin, tMax float64, rec *HitRecord) bool {
	t := (s.K - r.Origin.Y()) / r.Dir.Y()
	if t < tMin || t > tMax {
		return false
	}
	x := r.Origin.X() + t*r.Dir.X()
	z := r.Origin.Z() + t*r.Dir.Z()
	if x < s.X0 || x > s.X1 || z < s.Z0 || z > s.Z1 {
		return false
	}
	rec.U = (x - s.X0) / (s.X1 - s.X0)
	rec.V = (z - s.Z0) / (s.Z1 - s.Z0)
	rec.T = t
	outwardNormal := NewVec3(0, 1, 0)
	rec.SetFaceNormal(r, outwardNormal)
	rec.Mat = s.Mat
	rec.Point = r.At(t)
	return true
}

func (s *XZRect) BoundingBox(time0, time1 float64, outputBox *AABB) bool {
	*outputBox = *NewAABB(
		NewVec3(s.X0, s.K-0.0001, s.Z0),
		NewVec3(s.X1, s.K+0.0001, s.Z1),
	)
	return true
}

type YZRect struct {
	Y0, Z0 float64
	Y1, Z1 float64
	K      float64
	Mat    Material
}

func NewYZRect(y0, z0, y1, z1, k float64, mat Material) *YZRect {
	return &YZRect{
		Y0:  y0,
		Z0:  z0,
		Y1:  y1,
		Z1:  z1,
		K:   k,
		Mat: mat,
	}
}

func (s *YZRect) Hit(r *Ray, tMin, tMax float64, rec *HitRecord) bool {
	t := (s.K - r.Origin.X()) / r.Dir.X()
	if t < tMin || t > tMax {
		return false
	}
	y := r.Origin.Y() + t*r.Dir.Y()
	z := r.Origin.Z() + t*r.Dir.Z()
	if y < s.Y0 || y > s.Y1 || z < s.Z0 || z > s.Z1 {
		return false
	}
	rec.U = (y - s.Y0) / (s.Y1 - s.Y0)
	rec.V = (z - s.Z0) / (s.Z1 - s.Z0)
	rec.T = t
	outwardNormal := NewVec3(1, 0, 0)
	rec.SetFaceNormal(r, outwardNormal)
	rec.Mat = s.Mat
	rec.Point = r.At(t)
	return true
}

func (s *YZRect) BoundingBox(time0, time1 float64, outputBox *AABB) bool {
	*outputBox = *NewAABB(
		NewVec3(s.K-0.0001, s.Y0, s.Z0),
		NewVec3(s.K+0.0001, s.Y1, s.Z1),
	)
	return true
}

type Box struct {
	Min   Point3
	Max   Point3
	Sides *HittableList
}

func NewBox(p0, p1 Point3, mat Material) *Box {
	box := &Box{
		Min:   p0,
		Max:   p1,
		Sides: NewHittableList(),
	}

	box.Sides.Add(NewXYRect(p0.X(), p0.Y(), p1.X(), p1.Y(), p1.Z(), mat))
	box.Sides.Add(NewXYRect(p0.X(), p0.Y(), p1.X(), p1.Y(), p0.Z(), mat))

	box.Sides.Add(NewXZRect(p0.X(), p0.Z(), p1.X(), p1.Z(), p1.Y(), mat))
	box.Sides.Add(NewXZRect(p0.X(), p0.Z(), p1.X(), p1.Z(), p0.Y(), mat))

	box.Sides.Add(NewYZRect(p0.Y(), p0.Z(), p1.Y(), p1.Z(), p1.X(), mat))
	box.Sides.Add(NewYZRect(p0.Y(), p0.Z(), p1.Y(), p1.Z(), p0.X(), mat))

	return box
}

func (b *Box) Hit(r *Ray, tMin, tMax float64, rec *HitRecord) bool {
	return b.Sides.Hit(r, tMin, tMax, rec)
}

func (b *Box) BoundingBox(time0, time1 float64, outputBox *AABB) bool {
	*outputBox = *NewAABB(
		b.Min,
		b.Max,
	)
	return true
}

type Translate struct {
	Obj    Hittable
	Offset Vec3
}

func NewTranslate(obj Hittable, disp Vec3) *Translate {
	return &Translate{
		Obj:    obj,
		Offset: disp,
	}
}

func (t *Translate) Hit(r *Ray, tMin, tMax float64, rec *HitRecord) bool {
	movedR := NewRay(r.Origin.Sub(t.Offset), r.Dir, r.Time)
	if !t.Obj.Hit(movedR, tMin, tMax, rec) {
		return false
	}

	rec.Point = rec.Point.Add(t.Offset)
	rec.SetFaceNormal(movedR, rec.Normal)

	return true
}

func (t *Translate) BoundingBox(time0, time1 float64, outputBox *AABB) bool {
	if t.Obj.BoundingBox(time0, time1, outputBox) {
		return false
	}

	*outputBox = *NewAABB(
		outputBox.Min.Add(t.Offset),
		outputBox.Max.Add(t.Offset),
	)

	return true
}

type RotateY struct {
	Obj                Hittable
	SinTheta, CosTheta float64
	hasBox             bool
	Box                *AABB
}

func NewRotateY(obj Hittable, angle float64) *RotateY {
	r := &RotateY{Obj: obj, Box: new(AABB)}

	radians := Rad(angle)
	r.SinTheta = math.Sin(radians)
	r.CosTheta = math.Cos(radians)
	r.hasBox = obj.BoundingBox(0, 1, r.Box)

	posInf := math.Inf(1)
	negInf := math.Inf(-1)
	min := NewVec3(posInf, posInf, posInf)
	max := NewVec3(negInf, negInf, negInf)

	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				x := float64(i)*r.Box.Max.X() + float64(1-i)*r.Box.Min.X()
				y := float64(j)*r.Box.Max.Y() + float64(1-j)*r.Box.Min.Y()
				z := float64(k)*r.Box.Max.Z() + float64(1-k)*r.Box.Min.Z()

				newX := r.CosTheta*x + r.SinTheta*z
				newZ := -r.SinTheta*x + r.CosTheta*z

				tester := NewVec3(newX, y, newZ)

				for c := 0; c < 3; c++ {
					min[c] = math.Min(min[c], tester[c])
					max[c] = math.Max(max[c], tester[c])
				}
			}
		}
	}

	r.Box = NewAABB(min, max)

	return r
}

func (ry *RotateY) Hit(r *Ray, tMin, tMax float64, rec *HitRecord) bool {
	origin := r.Origin
	dir := r.Dir

	origin[0] = ry.CosTheta*r.Origin[0] - ry.SinTheta*r.Origin[2]
	origin[2] = ry.SinTheta*r.Origin[0] + ry.CosTheta*r.Origin[2]

	dir[0] = ry.CosTheta*r.Dir[0] - ry.SinTheta*r.Dir[2]
	dir[2] = ry.SinTheta*r.Dir[0] + ry.CosTheta*r.Dir[2]

	rotatedR := NewRay(origin, dir, r.Time)

	if !ry.Obj.Hit(rotatedR, tMin, tMax, rec) {
		return false
	}

	p := rec.Point
	normal := rec.Normal

	p[0] = ry.CosTheta*rec.Point[0] + ry.SinTheta*rec.Point[2]
	p[2] = -ry.SinTheta*rec.Point[0] + ry.CosTheta*rec.Point[2]

	normal[0] = ry.CosTheta*rec.Normal[0] + ry.SinTheta*rec.Normal[2]
	normal[2] = -ry.SinTheta*rec.Normal[0] + ry.CosTheta*rec.Normal[2]

	rec.Point = p
	rec.SetFaceNormal(rotatedR, normal)

	return true
}

func (r *RotateY) BoundingBox(time0, time1 float64, outputBox *AABB) bool {
	*outputBox = *r.Box
	return r.hasBox
}

type ConstantMedium struct {
	Boundary      Hittable
	PhaseFunction Material
	negInvDensity float64
}

func NewConstantMedium(b Hittable, d float64, tex Texture) *ConstantMedium {
	return &ConstantMedium{
		Boundary:      b,
		PhaseFunction: NewIsotropic(tex),
		negInvDensity: -1 / d,
	}
}

func (cm *ConstantMedium) Hit(r *Ray, tMin, tMax float64, rec *HitRecord) bool {
	var rec1, rec2 HitRecord
	if !cm.Boundary.Hit(r, math.Inf(-1), math.Inf(1), &rec1) {
		return false
	}
	if !cm.Boundary.Hit(r, rec1.T+0.0001, math.Inf(1), &rec2) {
		return false
	}

	if rec1.T < tMin {
		rec1.T = tMin
	}
	if rec2.T > tMax {
		rec2.T = tMax
	}

	if rec1.T >= rec2.T {
		return false
	}

	if rec1.T < 0 {
		rec1.T = 0
	}

	rayLength := r.Dir.Len()
	distanceInsideBoundary := (rec2.T - rec1.T) * rayLength
	hitDistance := cm.negInvDensity * math.Log(rand.Float64())

	if hitDistance > distanceInsideBoundary {
		return false
	}

	rec.T = rec1.T + hitDistance/rayLength
	rec.Point = r.At(rec.T)

	rec.Normal = NewVec3(1, 0, 0)
	rec.frontFace = true
	rec.Mat = cm.PhaseFunction

	return true
}

func (cm *ConstantMedium) BoundingBox(time0, time1 float64, outputBox *AABB) bool {
	return cm.Boundary.BoundingBox(time0, time1, outputBox)
}
