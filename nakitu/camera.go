package nakitu

import "math"

type Camera struct {
	Origin          Point3
	LowerLeftCorner Point3
	Horizontal      Vec3
	Vertical        Vec3
	U               Vec3
	V               Vec3
	W               Vec3
	LensRadius      float64
	Time0           float64
	Time1           float64
}

func NewCamera(
	lookFrom Point3,
	lookAt Point3,
	vUp Vec3,
	vfov float64,
	aspectRatio float64,
	aperture float64,
	focusDist float64,
) *Camera {
	theta := Rad(vfov)
	h := math.Tan(theta / 2)
	viewportHeight := 2.0 * h
	viewportWidth := aspectRatio * viewportHeight

	w := lookFrom.Sub(lookAt).Unit()
	u := vUp.Cross(w).Unit()
	v := w.Cross(u)

	origin := lookFrom
	horizontal := u.Mulf(focusDist * viewportWidth)
	vertical := v.Mulf(focusDist * viewportHeight)
	lowerLeftCorner := origin.
		Sub(horizontal.Divf(2)).
		Sub(vertical.Divf(2)).
		Sub(w.Mulf(focusDist))
	lendRadius := aperture / 2

	return &Camera{
		Origin:          origin,
		Horizontal:      horizontal,
		Vertical:        vertical,
		LowerLeftCorner: lowerLeftCorner,
		U:               u,
		V:               v,
		W:               w,
		LensRadius:      lendRadius,
		Time0:           0,
		Time1:           0,
	}
}

func (c *Camera) GetRay(s, t float64) *Ray {
	rd := RandomInUnitDisk().Mulf(c.LensRadius)
	offset := c.U.Mulf(rd.X()).Add(c.V.Mulf(rd.Y()))

	return NewRay(
		c.Origin.Add(offset),
		c.LowerLeftCorner.
			Add(c.Horizontal.Mulf(s)).
			Add(c.Vertical.Mulf(t)).
			Sub(c.Origin).
			Sub(offset),
		Random(c.Time0, c.Time1),
	)

}
