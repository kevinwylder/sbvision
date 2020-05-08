package interpolate

import (
	"math"

	"github.com/kevinwylder/sbvision"
)

type QuaternionFunction interface {
	At(t float64) sbvision.Quaternion
	Duration() float64
}

type lowPassFilter struct {
	control   []sbvision.Quaternion
	resampled []sbvision.Quaternion
	harmonics []sbvision.Quaternion

	totalSweep float64
	dsweep     []float64
	sweep      []float64
}

// Inerpolate will return a quaternion interpolation function
func Clip(clip *sbvision.Clip) QuaternionFunction {
	interp := lowPassFilter{}

	// compute the sweep angles between neighbor points (dsweep) and since the start (sweep)
	for i := clip.Start; i <= clip.End; i++ {
		j := i + 1
		if i == clip.End {
			j = clip.Start
		}
		diff := clip.Rotations[i].Divide(clip.Rotations[j])
		diff.Normalize()
		dtheta := math.Acos(diff[0])
		interp.dsweep = append(interp.dsweep, dtheta)
		interp.totalSweep += dtheta
		interp.sweep = append(interp.sweep, interp.totalSweep)
		interp.control = append(interp.control, clip.Rotations[i])
	}

	interp.resample(500)
	interp.fourierTransform(11)
	return &interp
}

// resample follows the path defined by the input clip, but gives a constant distance between points
func (interp *lowPassFilter) resample(size int) {
	interp.resampled = make([]sbvision.Quaternion, size)
	k := 0 // the knot index
	for i := 0; i < size; i++ {
		t := interp.totalSweep*float64(i)/float64(size) - interp.sweep[k]
		for t > 1 { // if we would be slerping past the end of this pair of knots, then we move on to the next knot
			t -= interp.dsweep[k]
			k++
		}
		nextK := (k + 1) % len(interp.control)
		interp.resampled[i] = interp.control[k].Lerp(interp.control[nextK], t)
	}
}

func (interp *lowPassFilter) fourierTransform(harmonics int) {
	interp.harmonics = make([]sbvision.Quaternion, harmonics)
	n := float64(len(interp.resampled))
	for h := -harmonics / 2; h < (harmonics+1)/2; h++ {
		// compute avg(q * exp(2 pi sqrt(-1) h)
		for i := 0; i < len(interp.resampled); i++ {
			theta := float64(h) * math.Pi * float64(i) / n
			convolve := interp.resampled[i].Multiply(sbvision.NewQuaternion(
				math.Cos(theta),
				math.Sin(theta),
				0,
				0,
			))
			interp.harmonics[h][0] += convolve[0]
			interp.harmonics[h][1] += convolve[1]
			interp.harmonics[h][2] += convolve[2]
			interp.harmonics[h][3] += convolve[3]
		}
		interp.harmonics[h].Scale(1 / n)
	}
}

func (interp *lowPassFilter) At(time float64) sbvision.Quaternion {
	t := interp.sweep[int(time)]
	t += (time - math.Floor(time)) * interp.dsweep[int(time)]
	harmonics := len(interp.harmonics)
	var r sbvision.Quaternion
	for h := -harmonics / 2; h < (harmonics+1)/2; h++ {
		theta := float64(h) * math.Pi * t
		convolve := interp.harmonics[h+harmonics/2].Multiply(sbvision.NewQuaternion(
			math.Cos(-theta),
			math.Sin(-theta),
			0,
			0,
		))
		r[0] += convolve[0]
		r[1] += convolve[1]
		r[2] += convolve[2]
		r[3] += convolve[3]
	}
	r.Normalize()
	return r
}

func (interp *lowPassFilter) Duration() float64 {
	return float64(len(interp.control))
}
