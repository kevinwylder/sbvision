package interpolate

import (
	"math"

	"github.com/kevinwylder/sbvision"
)

type constantSweep struct {
	sweep      []float64
	dsweep     []float64
	control    []sbvision.Quaternion
	totalSweep float64
}

func LinearConstantRate(clip *sbvision.Clip) QuaternionFunction {
	var c constantSweep
	// compute the sweep angles between neighbor points (dsweep) and since the start (sweep)
	for i := clip.Start; i <= clip.End; i++ {
		j := i + 1
		if i == clip.End {
			j = clip.Start
		}
		diff := clip.Rotations[i].Divide(clip.Rotations[j])
		diff.Normalize()
		dtheta := math.Acos(diff[0])
		c.dsweep = append(c.dsweep, dtheta)
		c.totalSweep += dtheta
		c.sweep = append(c.sweep, c.totalSweep)
		c.control = append(c.control, clip.Rotations[i])
	}

	return c
}

func (l constantSweep) At(t float64) sbvision.Quaternion {
	t = math.Mod(math.Mod(t, l.Duration())+l.Duration(), l.Duration())
	for k := 0; k < len(l.control); k++ {
		p := (t - l.sweep[k]) / l.dsweep[k]
		if p < 1 {
			return l.control[k].Lerp(l.control[(k+1)%len(l.control)], p)
		}
	}
	return sbvision.NewQuaternion(1, 0, 0, 0)
}

func (l constantSweep) Duration() float64 {
	return l.totalSweep
}
