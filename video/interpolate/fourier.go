package interpolate

import (
	"math"

	"github.com/kevinwylder/sbvision"
)

type lowPassFilter []sbvision.Quaternion

func LowPassFilter(clip *sbvision.Clip) QuaternionFunction {
	f := Linear(clip)
	samples := 500
	freqs := 21
	harmonics := make([]sbvision.Quaternion, freqs)
	for i := 0; i < freqs; i++ {
		h := i - freqs/2
		// compute avg(q * exp(2 pi sqrt(-1) h)
		for j := 0; j < samples; j++ {
			t := float64(j) / float64(samples)
			theta := 2 * float64(h) * math.Pi * t
			convolve := f.At(f.Duration() * t).Multiply(sbvision.NewQuaternion(
				math.Cos(theta),
				math.Sin(theta),
				0,
				0,
			))
			harmonics[i][0] += convolve[0]
			harmonics[i][1] += convolve[1]
			harmonics[i][2] += convolve[2]
			harmonics[i][3] += convolve[3]
		}
		harmonics[i][0] /= float64(samples)
		harmonics[i][1] /= float64(samples)
		harmonics[i][2] /= float64(samples)
		harmonics[i][3] /= float64(samples)
	}

	return lowPassFilter(harmonics)
}

func (harmonics lowPassFilter) At(t float64) sbvision.Quaternion {
	freqs := len(harmonics)
	var r sbvision.Quaternion
	for i := 0; i < freqs; i++ {
		h := i - freqs/2
		theta := 2 * float64(h) * math.Pi * t
		convolve := harmonics[i].Multiply(sbvision.NewQuaternion(
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
	r = r.Normalize()
	return r
}

func (lowPassFilter) Duration() float64 {
	return 1.
}
