package controllers

type Controllers struct {
	CacheBasedPolicy *CacheBasedPolicy
	Estimator        *Estimator
}

func NewControllers() *Controllers {
	return &Controllers{
		CacheBasedPolicy: &CacheBasedPolicy{},
		Estimator:        NewEstimator(),
	}
}
