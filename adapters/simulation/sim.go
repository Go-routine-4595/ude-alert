package simulationpackage

import (
	"math"
)

type SimServer struct {
	Frequency int `yaml:"frequency"`
	MaxPeak   int `yaml:"maxPeak"`
}

type Sim struct {
	freq    int
	maxPeak int
	count   int
}

func NewSim(freq int, maxPeak int) *Sim {
	var (
		s Sim
	)
	s = Sim{
		freq:    freq,
		maxPeak: maxPeak,
	}

	return &s
}

func (s *Sim) readData() int {
	var (
		v int
	)
	s.count++
	v = s.simCalc(s.count)

	return v
}

func (s *Sim) simCalc(t int) int {
	var (
		time int
		d    float64
		r    float64
	)

	time = (300 * t) % s.freq
	d = (float64(time) * 6) / float64(s.freq)
	r = math.Pow(math.E, (d-3)) / (math.Pow(1+math.Pow(math.E, (2*(d-4))), 2))

	return int(math.Ceil(float64(s.maxPeak) * r))
}
