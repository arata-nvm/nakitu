package nakitu

import (
	"log"
	"sort"
)

type BVHNode struct {
	Left  Hittable
	Right Hittable
	Box   *AABB
}

func NewBVHNode(list *HittableList, time0, time1 float64) *BVHNode {
	return newBVHNode(list.Objects, 0, len(list.Objects), time0, time1)
}

func newBVHNode(srcObjects []Hittable, start, end int, time0, time1 float64) *BVHNode {
	node := &BVHNode{}

	objects := make([]Hittable, len(srcObjects))
	copy(objects, srcObjects)

	axis := RandomInt(0, 3)
	objectSpan := end - start

	switch objectSpan {
	case 1:
		node.Left = objects[start]
		node.Right = objects[start]
	case 2:
		if boxCompare(objects[start], objects[start+1], axis) {
			node.Left = objects[start]
			node.Right = objects[start+1]
		} else {
			node.Left = objects[start+1]
			node.Right = objects[start]
		}
	default:
		sort.Slice(objects, func(i, j int) bool {
			return boxCompare(objects[i], objects[j], axis)
		})

		mid := start + objectSpan/2
		node.Left = newBVHNode(objects, start, mid, time0, time1)
		node.Right = newBVHNode(objects, mid, end, time0, time1)
	}

	var boxLeft, boxRight AABB
	if !node.Left.BoundingBox(time0, time1, &boxLeft) || !node.Right.BoundingBox(time0, time1, &boxRight) {
		log.Fatalln("No bounding box in BVHNode constructor.")
	}

	node.Box = SurroundingBox(&boxLeft, &boxRight)

	return node
}

func boxCompare(a, b Hittable, axis int) bool {
	var boxA, boxB AABB

	if !a.BoundingBox(0, 0, &boxA) || !b.BoundingBox(0, 0, &boxB) {
		log.Fatalln("No bounding box in BVHNode constructor.")
	}

	return boxA.Min[axis] < boxB.Min[axis]
}

func (b *BVHNode) Hit(r *Ray, tMin, tMax float64, rec *HitRecord) bool {
	if !b.Box.Hit(r, tMin, tMax) {
		return false
	}

	hitLeft := b.Left.Hit(r, tMin, tMax, rec)
	hitRight := b.Right.Hit(r, tMin, tMax, rec)

	return hitLeft || hitRight
}

func (b *BVHNode) BoundingBox(time0, time1 float64, outputBox *AABB) bool {
	*outputBox = *b.Box
	return true
}
