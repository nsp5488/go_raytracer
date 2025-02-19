package hittable

import (
	"sort"

	"github.com/nsp5488/go_raytracer/internal/aabb"
	"github.com/nsp5488/go_raytracer/internal/interval"
	"github.com/nsp5488/go_raytracer/internal/ray"
)

// Bounded Volume Hierarchy
type BVHNode struct {
	left  Hittable
	right Hittable

	bbox *aabb.AABB
}

func BuildBVH(list *HittableList) *BVHNode {
	return bvhHelper(list, 0, len(list.objects))
}

func boxCompare(a, b Hittable, axis int) bool {
	aAxis := a.BBox().AxisInterval(axis)
	bAxis := b.BBox().AxisInterval(axis)
	if aAxis.Min != bAxis.Min {
		return aAxis.Min < bAxis.Min
	}
	return aAxis.Max < bAxis.Max
}

func bvhHelper(list *HittableList, start, end int) *BVHNode {
	bbox := aabb.EmptyBBox()
	for i := start; i < end; i++ {
		bbox = aabb.FromBBoxes(bbox, list.objects[i].BBox())
	}
	axis := bbox.LongestAxis()
	objSpan := end - start
	var l, r Hittable
	if objSpan == 1 {
		l, r = list.objects[start], list.objects[start]
	} else if objSpan == 2 {
		l, r = list.objects[start], list.objects[start+1]
	} else {
		subslice := list.objects[start:end]
		sort.Slice(subslice, func(i, j int) bool {
			return boxCompare(subslice[i], subslice[j], axis)
		})
		mid := start + objSpan/2
		l = bvhHelper(list, start, mid)
		r = bvhHelper(list, mid, end)
	}

	return &BVHNode{left: l, right: r, bbox: bbox}
}

func (bvh *BVHNode) BBox() *aabb.AABB {
	return bvh.bbox
}

func (bvh *BVHNode) Hit(r *ray.Ray, rayT interval.Interval, record *HitRecord) bool {
	if !bvh.bbox.Hit(r, rayT) {
		return false
	}
	hitLeft := bvh.left.Hit(r, rayT, record)

	if hitLeft {
		rayT.Max = record.t
	}

	hitRight := bvh.right.Hit(r, rayT, record)

	return hitRight || hitLeft
}
