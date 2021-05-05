package nakitu

import "math"

type AABB struct {
	Min Point3
	Max Point3
}

func NewAABB(min, max Point3) *AABB {
	return &AABB{Min: min, Max: max}
}

func (a *AABB) Hit(r *Ray, tMin, tMax float64) bool {
	for i := 0; i < 3; i++ {
		invD := 1.0 / r.Dir[i]
		t0 := (a.Min[i] - r.Origin[i]) * invD
		t1 := (a.Max[i] - r.Origin[i]) * invD
		if invD < 0 {
			t0, t1 = t1, t0
		}
		if t0 > tMin {
			tMin = t0
		}
		if t1 < tMax {
			tMax = t1
		}
		if tMax <= tMin {
			return false
		}
	}

	return true
}

func SurroundingBox(box0, box1 *AABB) *AABB {
	small := NewVec3(
		math.Min(box0.Min[0], box1.Min[0]),
		math.Min(box0.Min[1], box1.Min[1]),
		math.Min(box0.Min[2], box1.Min[2]),
	)
	big := NewVec3(
		math.Max(box0.Max[0], box1.Max[0]),
		math.Max(box0.Max[1], box1.Max[1]),
		math.Max(box0.Max[2], box1.Max[2]),
	)
	return NewAABB(small, big)
}
