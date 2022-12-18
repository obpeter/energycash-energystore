package calculation

import (
	"at.ourproject/energystore/model"
	"at.ourproject/energystore/utils"
	"math"
)

func AllocDynamic(line *model.RawSourceLine) *model.Matrix {

	resultArray := make([]float64, len(line.Consumers)*len(line.Producers))
	lineResult := model.MakeMatrix(resultArray, len(line.Consumers), len(line.Producers))

	consumerSum := utils.Sum(line.Consumers)
	producerSum := utils.Sum(line.Producers)

	//consumerMean := consumerSum / float32(len(line.Consumers))
	// Mean value distributed over
	producerMean := producerSum / float64(len(line.Consumers))

	valueAboveMean := make([]float64, len(line.Consumers))
	for i, l := range line.Consumers {
		valueAboveMean[i] = float64(math.Max(0, float64(l-producerMean)))
	}

	diffProdCons := utils.Sum(valueAboveMean)
	var allocation float64 = float64(0)
	if diffProdCons > 0 {
		allocation = math.Max(0, float64((consumerSum-producerSum)/diffProdCons))
	}

	for i, l := range line.Consumers {
		gridValue := valueAboveMean[i] * allocation
		for j, pl := range line.Producers {
			var prod_factor float64 = 0
			if producerSum > float64(0) {
				prod_factor = pl / producerSum
			}
			lineResult.SetElm(i, j, float64((l-gridValue)*prod_factor))
		}
	}
	return lineResult
}

func AllocDynamic1(line *model.RawSourceLine) *model.Matrix {

	lenConsumers := int(math.Max(float64(len(line.Consumers)), 1))
	lenProducers := int(math.Max(float64(len(line.Producers)), 1))

	resultArray := make([]float64, lenConsumers*lenProducers)
	lineResult := model.MakeMatrix(resultArray, lenConsumers, lenProducers)

	consumerSum := utils.Sum(line.Consumers)
	producerSum := utils.Sum(line.Producers)

	var alloc_prod_to_cons_factor = float64(0)
	if producerSum > float64(0) && consumerSum > float64(0) {
		alloc_prod_to_cons_factor = consumerSum / producerSum
	}

	for i, l := range line.Consumers {
		greenValue := float64(0)
		if alloc_prod_to_cons_factor > float64(0) {
			greenValue = l / alloc_prod_to_cons_factor
		}
		for j, pl := range line.Producers {
			var prod_factor float64 = float64(0)
			if producerSum > float64(0) {
				prod_factor = pl / producerSum
			}
			lineResult.SetElm(i, j, math.Min(float64(l), float64(greenValue*prod_factor)))
		}
	}
	return model.Multiply(lineResult, model.NewUniformMatrix(lineResult.Cols, 1))
}