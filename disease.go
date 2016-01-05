package main

type Disease struct {
	Virality float64 // The probability of getting infected if exposed.
	Timer    []int   // Times used to show the progression of the illness.
}

func NewDisease(virality float64, timer []int) *Disease {
	disease := new(Disease)
	disease.Virality = virality
	disease.Timer = timer
	return disease
}
