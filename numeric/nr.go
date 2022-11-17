package numeric

import "math"

type Function1D func(x float64) float64
type Derivative1D func(fun Function1D, x float64) float64

type NRReturn struct {
	Converged bool
	X         float64
	Iter      int
	Error     float64
}

func GetNumericalDerivative(h float64) Derivative1D {
	return func(fun Function1D, x float64) float64 {
		der := (fun(x+h) - fun(x)) / h
		return der
	}
}

func NewtonRaphson1D(fun Function1D, x0 float64, tol float64, maxIter int, derivative Derivative1D) NRReturn {
	xOld := x0
	xNew := x0
	iter := 0

	err := 1e10

	for err > tol {
		fOld := fun(xOld)
		der := derivative(fun, xOld)
		xNew = xOld - fOld/der

		//err = 0.5 * math.Pow(fun(xNew), 2)
		err = math.Abs(xOld - xNew)

		if iter > maxIter {
			return NRReturn{
				Converged: false,
				X:         xNew,
				Iter:      iter,
				Error:     err,
			}
		}

		iter += 1
		xOld = xNew

	}
	return NRReturn{
		Converged: true,
		X:         xNew,
		Iter:      iter,
		Error:     err,
	}

}
