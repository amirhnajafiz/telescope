package models

// Estimane is a structure that holds the bandwidth estimates for a client
type Estimane struct {
	Uncached float64 // Tn
	Cached   float64 // Tg
	CurBW    float64 // Tc
	Alpha    float64 // smoothing factor
}
