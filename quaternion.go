package sbvision

import "math"

func NewQuaternion(r, i, j, k float64) Quaternion {
	return [4]float64{r, i, j, k}
}

func (q Quaternion) Scale(b float64) {
	q[0] /= b
	q[1] /= b
	q[2] /= b
	q[3] /= b
}

func (q Quaternion) Normalize() {
	q.Scale(math.Sqrt(q.Dot(q)))
}

func (q Quaternion) Dot(b Quaternion) float64 {
	return q[0]*b[0] + q[1]*b[1] + q[2]*b[2] + q[3]*b[3]
}

func (q Quaternion) Multiply(b Quaternion) Quaternion {
	return NewQuaternion(
		q[0]*b[0]-q[1]*b[1]-q[2]*b[2]-q[3]*b[3],
		q[0]*b[1]+q[1]*b[0]+q[2]*b[3]-q[3]*b[2],
		q[0]*b[2]-q[1]*b[3]+q[2]*b[0]+q[3]*b[1],
		q[0]*b[3]+q[1]*b[2]-q[2]*b[1]+q[3]*b[0],
	)
}

func (q Quaternion) Divide(b Quaternion) Quaternion {
	return NewQuaternion(
		b[0]*q[0]+b[1]*q[1]+b[2]*q[2]+b[3]*q[3],
		-b[0]*q[1]+b[1]*q[0]-b[2]*q[3]+b[3]*q[2],
		-b[0]*q[2]+b[1]*q[3]+b[2]*q[0]-b[3]*q[1],
		-b[0]*q[3]-b[1]*q[2]+b[2]*q[1]+b[3]*q[0],
	)
}

func (q Quaternion) Lerp(b Quaternion, t float64) Quaternion {
	dot := q.Dot(b)
	if dot > 0.9995 {
		// lerp them directly
		r := NewQuaternion(
			q[0]+(b[0]-q[0])*t,
			q[1]+(b[1]-q[1])*t,
			q[2]+(b[2]-q[2])*t,
			q[3]+(b[3]-q[3])*t,
		)
		r.Normalize()
		return r
	}
	theta0 := math.Acos(dot)
	theta := theta0 * t
	sinTheta := math.Sin(theta)
	sinTheta0 := math.Sin(theta0)
	s0 := math.Cos(theta) - dot*sinTheta/sinTheta0
	s1 := sinTheta / sinTheta0
	r := NewQuaternion(
		s0*q[0]+s1*b[0],
		s0*q[1]+s1*b[1],
		s0*q[2]+s1*b[2],
		s0*q[3]+s1*b[3],
	)
	r.Normalize()
	return r
}
