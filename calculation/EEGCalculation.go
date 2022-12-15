package calculation

import (
	"at.ourproject/energystore/model"
	"at.ourproject/energystore/store"
	"at.ourproject/energystore/store/ebow"
	"at.ourproject/energystore/utils"
	"fmt"
	"time"
)

type CalcHandler func(*store.BowStorage, string) (*model.Matrix, *model.Matrix, float64)

func CalculateEEG(db *store.BowStorage, period string) (*model.Matrix, *model.Matrix, float64) {
	metaMap, err := store.GetConsumerMetaMap(db)
	if err != nil {
		return nil, nil, 0
	}
	iter := db.GetLinePrefix(fmt.Sprintf("CP/%s", period))
	defer iter.Close()

	var _line model.RawSourceLine

	var rAlloc *model.Matrix
	var rCons *model.Matrix
	var pSum float64 = 0.0
	defaultConsumerLen := len(metaMap)

	for iter.Next(&_line) {
		//line := transformConsumer(&_line)
		line := _line.Copy(defaultConsumerLen)
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

func CalculateMonthlyDash(db *store.BowStorage, year string, calc CalcHandler) error {
	mounth := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

	metaMap, err := store.GetConsumerMetaMap(db)
	if err != nil {
		return nil
	}
	var annual_am, annual_cm *model.Matrix = &model.Matrix{}, &model.Matrix{}
	var annual_ps float64 = 0.0

	defaultMatrix := model.NewMatrix(len(metaMap), 1)
	verifyResult := func(matrix *model.Matrix) *model.Matrix {
		if matrix == nil {
			return defaultMatrix
		} else {
			return matrix
		}
	}
	for _, m := range mounth {
		am, cm, ps := calc(db, fmt.Sprintf("%s/%.2d/", year, m))
		am = verifyResult(am)
		cm = verifyResult(cm)

		if err := db.SetReport(
			&model.EnergyReport{
				Id:        fmt.Sprintf("MRP/%s/%.2d", year, m),
				Consumed:  cm.Elements,
				Allocated: am.Elements,
				Produced:  ps}); err != nil {
			return err
		}

		_ = annual_am.Add(am)
		_ = annual_cm.Add(cm)

		annual_ps += ps
		fmt.Printf("AllocMonth (%d) %f - %+v\n", m, ps, am)
	}

	//fmt.Printf("Annual Alloc-Report: %+v\n", annual_am)

	if err := db.SetReport(
		&model.EnergyReport{
			Id:        fmt.Sprintf("YRP/%s", year),
			Consumed:  annual_cm.Elements,
			Allocated: annual_am.Elements,
			Produced:  annual_ps}); err != nil {
		return err
	}

	return nil
}

func CalculateWeeklyReport(db *store.BowStorage, year, month int, calc CalcHandler) ([]*model.EnergyReport, *model.EnergyReport, error) {
	var report_am, report_cm *model.Matrix = &model.Matrix{}, &model.Matrix{}
	var report_ps float64 = 0.0

	metaMap, err := store.GetConsumerMetaMap(db)
	if err != nil {
		return nil, nil, err
	}
	defaultMatrix := model.NewMatrix(len(metaMap), 1)
	verifyResult := func(matrix *model.Matrix) *model.Matrix {
		if matrix == nil {
			return defaultMatrix
		} else {
			return matrix
		}
	}
	dayReports := []*model.EnergyReport{}

	t := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC)
	for day := 1; day <= t.Day(); day++ {
		am, cm, ps := calc(db, fmt.Sprintf("%d/%.2d/%.2d", year, month, day))
		//if am == nil || cm == nil {
		//	continue
		//}

		am = verifyResult(am)
		cm = verifyResult(cm)

		dayReports = append(dayReports, &model.EnergyReport{
			Id:        fmt.Sprintf("WRP/%d/%.2d/%.2d", year, month, day),
			Consumed:  cm.Elements,
			Allocated: am.Elements,
			Produced:  ps,
		})

		_ = report_am.Add(am)
		_ = report_cm.Add(cm)

		report_ps += ps
	}

	return dayReports, &model.EnergyReport{
		Id:        fmt.Sprintf("MRP/%d/%.2d", year, month),
		Consumed:  report_cm.Elements,
		Allocated: report_am.Elements,
		Produced:  report_ps}, nil
}

func CalculateYearlyReport(db *store.BowStorage, year int, calc CalcHandler) ([]*model.EnergyReport, *model.EnergyReport, error) {
	monthReports := []*model.EnergyReport{}
	resAnnual := &model.EnergyReport{}

	var line model.EnergyReport = model.EnergyReport{}

	iter := db.GetLinePrefix(fmt.Sprintf("MRP/%d/", year))
	for iter.Next(&line) {
		monthReports = append(monthReports, &model.EnergyReport{
			Id:        line.Id,
			Consumed:  line.Consumed,
			Allocated: line.Allocated,
			Produced:  line.Produced,
		})
		line = model.EnergyReport{}
	}

	if annual, err := db.GetReport(fmt.Sprintf("YRP/%d", year)); err != nil {
		if err != ebow.ErrNotFound {
			return nil, nil, err
		}
	} else {
		resAnnual = annual
	}
	return monthReports, resAnnual, nil
}

func transformConsumer(line *model.RawSourceLine) *model.RawSourceLine {
	result := &model.RawSourceLine{Id: line.Id, Consumers: []float64{}, Producers: line.Producers}
	for i := 0; i < len(line.Consumers); i += 3 {
		result.Consumers = append(result.Consumers, line.Consumers[i])
	}
	return result
}

func CalculatePeriodReport(db *store.BowStorage, from, to time.Time, calc CalcHandler) (*model.EnergyReport, map[string]*model.CounterPointMeta, error) {

	var report_am, report_cm *model.Matrix = &model.Matrix{}, &model.Matrix{}
	var report_ps float64 = 0

	year, months := utils.GetMonthDuration(from, to)
	verifyResult := func(matrix *model.Matrix, defaultMatrix *model.Matrix) *model.Matrix {
		if matrix == nil {
			return defaultMatrix
		} else {
			return matrix
		}
	}

	periodFrom := int(from.Month())
	totalMonth := months + periodFrom
	var periodTo int
	var metaMap map[string]*model.CounterPointMeta
	for {
		if totalMonth > 12 {
			periodTo = 12
		} else {
			periodTo = totalMonth
		}

		cm, am, prodSum, meta, err := CalculatePeriodWithinYearReport(db, year, periodFrom, periodTo, calc)
		if err != nil {
			return nil, nil, err
		}

		if len(meta) > 0 {

			defaultMatrix := model.NewMatrix(len(metaMap), 1)

			am = verifyResult(am, defaultMatrix)
			cm = verifyResult(cm, defaultMatrix)

			if am.Rows > 0 && am.Cols == defaultMatrix.Cols {
				_ = report_am.Add(am)
			}
			if cm.Rows > 0 && cm.Cols == defaultMatrix.Cols {
				_ = report_cm.Add(cm)
			}

			report_ps += prodSum
			metaMap = meta
		}

		totalMonth -= 12
		if totalMonth <= 0 {
			break
		}

		year += 1
		periodFrom = 1
	}

	return &model.EnergyReport{
		Id:        fmt.Sprintf("PRP/%d/%.2d", year, int(from.Month())),
		Consumed:  report_cm.Elements,
		Allocated: report_am.Elements,
		Produced:  report_ps}, metaMap, nil
}

func CalculatePeriodWithinYearReport(db *store.BowStorage, year, from, to int, calc CalcHandler) (*model.Matrix, *model.Matrix, float64, map[string]*model.CounterPointMeta, error) {
	var report_am, report_cm *model.Matrix = &model.Matrix{}, &model.Matrix{}
	var report_ps float64 = 0.0

	metaMap, err := store.GetConsumerMetaMap(db)
	if err != nil {
		return nil, nil, 0, nil, err
	}
	defaultMatrix := model.NewMatrix(len(metaMap), 1)
	verifyResult := func(matrix *model.Matrix) *model.Matrix {
		if matrix == nil {
			return defaultMatrix
		} else {
			return matrix
		}
	}

	for m := from; m <= to; m++ {
		am, cm, ps := calc(db, fmt.Sprintf("%d/%.2d/", year, m))

		am = verifyResult(am)
		cm = verifyResult(cm)

		_ = report_am.Add(am)
		_ = report_cm.Add(cm)

		report_ps += ps
	}

	return report_cm, report_am, report_ps, metaMap, nil
}
