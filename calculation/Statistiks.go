package calculation

import (
	"at.ourproject/energystore/model"
	"at.ourproject/energystore/store"
	"at.ourproject/energystore/utils"
	"fmt"
)

func GetCalcFunc(id string) CalcHandler {
	switch id {
	case "CalcWhenProduced":
		return CalcWhenProduced
	}
	return nil
}

func CalcWhenProduced(db *store.BowStorage, period string) (*model.Matrix, *model.Matrix, float64) {
	iter := db.GetLinePrefix(fmt.Sprintf("CP/%s", period))
	defer iter.Close()

	var _line model.RawSourceLine

	var rAlloc *model.Matrix
	var rCons *model.Matrix
	var pSum float64 = 0.0

	for iter.Next(&_line) {

		producerSum := utils.Sum(_line.Producers)
		if producerSum == 0 {
			continue
		}
		line := _line.Copy(len(_line.Consumers))
		m := AllocDynamic1(&line)

		if rCons == nil {
			rCons = model.MakeMatrix(line.Consumers, len(line.Consumers), 1)
		} else {
			rCons.Add(model.MakeMatrix(line.Consumers, len(line.Consumers), 1))
		}

		if rAlloc == nil {
			rAlloc = model.MakeMatrix(m.Elements, m.CountRows(), m.CountCols())
		} else {
			rAlloc.Add(m)
		}
		pSum += utils.Sum(line.Producers)
	}
	return rAlloc, rCons, pSum
}