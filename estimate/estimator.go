package estimate

type Estimator interface {
	Predict(x float64) float64
}
