package utils

import (
	"at.ourproject/energystore/model"
	"math"
)

func DetermineDirection(meteringPoint string) string {
	// 012345678901234567890123456789012
	// AT0030000000000000000000030032764

	switch meteringPoint[25] {
	case '3':
		return model.PRODUCER_DIRECTION
	default:
		return model.CONSUMER_DIRECTION
	}

}

var Insert = func(orig []float64, index int, value float64) []float64 {
	l := len(orig)
	if index >= l {
		target := make([]float64, index+1)
		copy(target, orig)
		orig = target
	}
	orig[index] = value
	return orig
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
