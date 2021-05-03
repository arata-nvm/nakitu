package main

import (
	"fmt"
	"io"
	"math"
	"math/rand"
)

type Vec3 [3]float64
type Point3 = Vec3
type Color = Vec3

func NewVec3(x, y, z float64) Vec3 {
	return Vec3{x, y, z}
}

func Zero() Vec3 {
	return NewVec3(0, 0, 0)
}

func RandomVec3() Vec3 {
	return NewVec3(
		rand.Float64(),
		rand.Float64(),
		rand.Float64(),
	)
}

func RandomVec3In(min, max float64) Vec3 {
	return NewVec3(
		Random(min, max),
		Random(min, max),
		Random(min, max),
	)
}

func RandomInUnitSphere() Vec3 {
	for {
		p := RandomVec3In(-1, 1)
		if p.Len() >= 1 {
			continue
		}

		return p
	}
}

func RandomUnitVector() Vec3 {
	return RandomInUnitSphere().Unit()
}

func RandomInHemisphere(normal Vec3) Vec3 {
	inUnitSphere := RandomInUnitSphere()
	if inUnitSphere.Dot(normal) > 0 {
		return inUnitSphere
	} else {
		return inUnitSphere.Neg()
	}
}

func RandomInUnitDisk() Vec3 {
	for {
		p := NewVec3(Random(-1, 1), Random(-1, 1), 0)
		if p.LenSquared() >= 1 {
			continue
		}

		return p
	}
}

func (v Vec3) X() float64 {
	return v[0]
}

func (v Vec3) Y() float64 {
	return v[1]
}

func (v Vec3) Z() float64 {
	return v[2]
}

func (v Vec3) Neg() Vec3 {
	return NewVec3(-v[0], -v[1], -v[2])
}

func (v Vec3) Add(other Vec3) Vec3 {
	return NewVec3(
		v[0]+other[0],
		v[1]+other[1],
		v[2]+other[2],
	)
}

func (v Vec3) Sub(other Vec3) Vec3 {
	return v.Add(other.Neg())
}

func (v Vec3) Mul(other Vec3) Vec3 {
	return NewVec3(
		v[0]*other[0],
		v[1]*other[1],
		v[2]*other[2],
	)
}

func (v Vec3) Mulf(f float64) Vec3 {
	return NewVec3(
		v[0]*f,
		v[1]*f,
		v[2]*f,
	)
}

func (v Vec3) Divf(f float64) Vec3 {
	return v.Mulf(1.0 / f)
}

func (v Vec3) Dot(u Vec3) float64 {
	return v[0]*u[0] + v[1]*u[1] + v[2]*u[2]
}

func (v Vec3) Cross(u Vec3) Vec3 {
	return NewVec3(
		u[1]*v[2]-u[2]*v[1],
		u[2]*v[0]-u[0]*v[2],
		u[0]*v[1]-u[1]*v[0],
	)
}

func (v Vec3) Unit() Vec3 {
	return v.Divf(v.Len())
}

func (v Vec3) Len() float64 {
	return math.Sqrt(v.LenSquared())
}

func (v Vec3) LenSquared() float64 {
	return v[0]*v[0] + v[1]*v[1] + v[2]*v[2]
}

func (v Vec3) NearZero() bool {
	s := 1e-8
	return math.Abs(v[0]) < s && math.Abs(v[1]) < s && math.Abs(v[1]) < s
}

func (v Vec3) Reflect(n Vec3) Vec3 {
	return v.Sub(n.Mulf(2 * v.Dot(n)))
}

func (uv Vec3) Refract(n Vec3, etaiOverEtat float64) Vec3 {
	cosTheta := math.Min(uv.Neg().Dot(n), 1)
	rOutPerp := uv.Add(n.Mulf(cosTheta)).Mulf(etaiOverEtat)
	rOutParallel := n.Mulf(-math.Sqrt(math.Abs(1 - rOutPerp.LenSquared())))
	return rOutPerp.Add(rOutParallel)
}

func (v Vec3) Dump(w io.Writer) {
	fmt.Fprintf(w, "%f %f %f", v[0], v[1], v[2])
}
