package interpolate

import (
	"math"

	"github.com/kevinwylder/sbvision"
)

// QuaternionFunction is a map from R to to the unit quaternions
type QuaternionFunction interface {
	At(t float64) sbvision.Quaternion
	Duration() float64
}

type linearFunction []sbvision.Quaternion

func Linear(clip *sbvision.Clip) QuaternionFunction {
	var points linearFunction = linearFunction{
		clip.Rotations[clip.End],
	}
	for i := clip.Start; i <= clip.End; i++ {
		a := points[len(points)-1]
		b := clip.Rotations[i]
		if i == clip.Start {
			b = clip.Rotations[clip.End]
		}
		var options = []sbvision.Quaternion{
			{b[0], b[1], b[2], b[3]},
			{b[3], -b[2], b[2], -b[0]},
			{-b[0], b[1], -b[2], b[3]},
			{-b[0], -b[1], -b[2], -b[3]},
		}
		var largestDot = -2.
		var bestIdx = 0
		for i, o := range options {
			dot := a.Dot(o)
			if dot > largestDot {
				largestDot = dot
				bestIdx = i
			}
		}
		points = append(points, options[bestIdx])
	}
	return points
}

func (l linearFunction) At(t float64) sbvision.Quaternion {
	curr := ((int(t) % len(l)) + len(l)) % len(l)
	return l[curr].Lerp(l[curr+1], t-math.Floor(t))
}

func (l linearFunction) Duration() float64 {
	return float64(len(l)) - 1.
}
