package main

type Disease struct {
	Virality float64     // The probability of getting infected if exposed.
	Timer    []int       // Times used to show the progression of the illness.
	Stages   [][]float64 // An adjacency matrix to describe the probability of moving to a particular stage
}

func NewDisease(virality float64, timer []int, matrix [][]float64) *Disease {
	disease := new(Disease)
	disease.Virality = virality
	disease.Timer = timer
	disease.Stages = matrix
	return disease
}

var SEIRMatrix = [][]float64{[]float64{1, 0, 0, 0}, []float64{0, 0, 1, 0}, []float64{0, 0, 0, 1}, []float64{0, 0, 0, 1}}
