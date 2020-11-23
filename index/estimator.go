package index

type Estimator interface {
	Predict(x float64) float64
}
