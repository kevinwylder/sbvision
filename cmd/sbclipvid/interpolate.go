package main

import (
	"math"

	"github.com/kevinwylder/sbvision"
)

func posmod(x, y int64) int64 {
	return ((x % y) + y) % y
}

func getQuaternion(clip *sbvision.Clip, time float64) quaternion {
	return clip.Rotations[clip.Start+posmod(int64(time)-clip.Start, clip.End-clip.Start+1)]
}

func interpolateBound(clip *sbvision.Clip, time float64) sbvision.Bound {
	return clip.Bounds[clip.Start+posmod(int64(time)-clip.Start, clip.End-clip.Start+1)]
}

type quaternion [4]float64

func interpolateRotation(clip *sbvision.Clip, time float64) quaternion {
	p0 := getQuaternion(clip, time-1)
	p1 := p0.nearestReflection(getQuaternion(clip, time))
	p2 := p1.nearestReflection(getQuaternion(clip, time+1))
	p3 := p2.nearestReflection(getQuaternion(clip, time+2))
	pa := qSlerp([4]float64{1, 0, 0, 0}, p0.diff(p2), 0.166666666).mult(p1)
	pb := qSlerp([4]float64{1, 0, 0, 0}, p3.diff(p1), 0.166666).mult(p2)
	return qBezier(p1, pa, pb, p2, time-math.Floor(time))
}

func (a quaternion) versor() (float64, float64, float64) {
	l := math.Sqrt(a[1]*a[1] + a[2]*a[2] + a[3]*a[3])
	return a[1] / l, a[2] / l, a[3] / l
}

func (a quaternion) nearestReflection(b quaternion) quaternion {
	var options = []quaternion{
		{b[0], b[1], b[2], b[3]},
		{b[3], -b[2], b[2], -b[0]},
		{-b[0], b[1], -b[2], b[3]},
		{-b[0], -b[1], -b[2], -b[3]},
	}
	var largestDot float64 = 0
	var bestIdx = 0
	for i, o := range options {
		dot := o[0]*a[0] + o[1]*a[1] + o[2]*a[2] + o[3]*a[3]
		if dot > largestDot {
			largestDot = dot
			bestIdx = i
		}
	}
	return options[bestIdx]

}

func qBezier(a0, a1, a2, a3 quaternion, t float64) quaternion {
	b1 := qSlerp(a1, a2, t)
	return qSlerp(
		qSlerp(
			qSlerp(
				a0,
				a1,
				t),
			b1,
			t),
		qSlerp(
			b1,
			qSlerp(
				a2,
				a3,
				t),
			t),
		t)
}

func qSlerp(a, b quaternion, t float64) quaternion {
	dot := a[0]*b[0] + a[1]*b[1] + a[2]*b[2] + a[3]*b[3]
	if dot > 0.9995 {
		// lerp them directly
		return quaternion([4]float64{
			a[0] + (b[0]-a[0])*t,
			a[1] + (b[1]-a[1])*t,
			a[2] + (b[2]-a[2])*t,
			a[3] + (b[3]-a[3])*t,
		}).normalize()
	}
	theta0 := math.Acos(dot)
	theta := theta0 * t
	sinTheta := math.Sin(theta)
	sinTheta0 := math.Sin(theta0)

	s0 := math.Cos(theta) - dot*sinTheta/sinTheta0
	s1 := sinTheta / sinTheta0
	return quaternion([4]float64{
		s0*a[0] + s1*b[0],
		s0*a[1] + s1*b[1],
		s0*a[2] + s1*b[2],
		s0*a[3] + s1*b[3],
	}).normalize()
}

func (a quaternion) normalize() quaternion {
	norm := math.Sqrt(a[0]*a[0] + a[1]*a[1] + a[2]*a[2] + a[3]*a[3])
	a[0] /= norm
	a[1] /= norm
	a[2] /= norm
	a[3] /= norm
	return a
}

func (a quaternion) mult(b quaternion) quaternion {
	return quaternion([4]float64{
		a[0]*b[0] - a[1]*b[1] - a[2]*b[2] - a[3]*b[3],
		a[0]*b[1] + a[1]*b[0] + a[2]*b[3] - a[3]*b[2],
		a[0]*b[2] - a[1]*b[3] + a[2]*b[0] + a[3]*b[1],
		a[0]*b[3] + a[1]*b[2] - a[2]*b[1] + a[3]*b[0],
	})
}

func (a quaternion) diff(b quaternion) quaternion {
	return quaternion([4]float64{
		b[0]*a[0] + b[1]*a[1] + b[2]*a[2] + b[3]*a[3],
		-b[0]*a[1] + b[1]*a[0] - b[2]*a[3] + b[3]*a[2],
		-b[0]*a[2] + b[1]*a[3] + b[2]*a[0] - b[3]*a[1],
		-b[0]*a[3] - b[1]*a[2] + b[2]*a[1] + b[3]*a[0],
	})
}
