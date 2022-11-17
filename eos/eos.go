package eos

import "math"

type EOS interface {
	V(TK float64, P float64) float64
}

type CubicEOS struct {
	Epsilon float64
	Sigma   float64
	Alpha   func(float64, float64) float64
}

func (eos CubicEOS) V(TK float64, P float64) float64 {
	return 1
}

func NewIdealEOS() CubicEOS {
	return CubicEOS{
		Epsilon: 0,
		Sigma:   0,
		Alpha: func(Tr float64, wacc float64) float64 {
			return 1
		},
	}
}

func NewPengRobinson() CubicEOS {
	return CubicEOS{
		Epsilon: 1 - math.Sqrt(2),
		Sigma:   1 + math.Sqrt(2),
		Alpha: func(Tr float64, wacc float64) float64 {
			return math.Pow(1+(0.037464+1.54226*wacc-0.26992*wacc*wacc)*(1-math.Sqrt(Tr)), 2)
		},
	}
}
