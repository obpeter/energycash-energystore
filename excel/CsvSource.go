package excel

import (
	"at.ourproject/energystore/model"
	"at.ourproject/energystore/store"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type CsvSource struct {
	csvReader   *csv.Reader
	consumerIdx []int
	producerIdx []int
	cpNames     []string
}

func OpenCsvFile(name string) (*CsvSource, *os.File, error) {
	c := &CsvSource{}
	var err error
	var csvFile *os.File
	csvFile, err = os.Open(name)
	if err != nil {
		return nil, nil, err
	}
	c.csvReader = csv.NewReader(csvFile)
	c.csvReader.Comma = ';'
	return c, csvFile, nil
}

func OpenCsvReader(r io.Reader) (*CsvSource, error) {
	c := &CsvSource{}
	c.csvReader = csv.NewReader(r)
	c.csvReader.Comma = ';'
	return c, nil
}

//func (e *CsvSource) init() error {
//	rec, err := e.readLine()
//	if err == io.EOF {
//		return nil
//	}
//	if err != nil {
//		return err
//	}
//
//	e.cpNames = rec[1:]
//
//	rec, err = e.readLine()
//	if err == io.EOF {
//		return nil
//	}
//	if err != nil {
//		return err
//	}
//
//	for i, s := range rec {
//		if s == "CONSUMPTION" {
//			e.consumerIdx = append(e.consumerIdx, i)
//			continue
//		}
//		if s == "GENERATION" {
//			e.producerIdx = append(e.producerIdx, i)
//		}
//	}
//	return nil
//}

func (e *CsvSource) Next() ([]string, error) {
	rec, err := e.csvReader.Read()
	//if err == io.EOF {
	//	return []string{}, nil
	//}
	if err != nil {
		return []string{}, err
	}
	return rec, err
}

//func (e *CsvSource) GetAllocationLine() {
//	rec, err := e.readLine()
//	if err != nil {
//		return
//	}
//
//	var line *model.EnergyAllocationLine = &model.EnergyAllocationLine{}
//
//	cpIdxSet := make(map[int]bool, len(e.consumerIdx))
//	for _, i := range e.consumerIdx {
//		cpIdxSet[i] = true
//	}
//
//	prodIdxSet := make(map[int]bool, len(e.producerIdx))
//	for _, i := range e.producerIdx {
//		prodIdxSet[i] = true
//	}
//
//	for i, s := range rec {
//		if _, ok := cpIdxSet[i]; ok {
//			//line.CounterPoints = append(line.CounterPoints, )
//			println(s, line)
//		}
//	}
//}

func ImportCsvFile(f *CsvSource, db *store.BowStorage) (map[int]bool, error) {

	var err error
	var cpMeta map[int]*model.CounterPointMeta
	//var rawData *model.RawSourceLine
	var rIdx int = 1
	var countConsumption int = 0
	var countGeneration int = 0
	var rawDatas []*model.RawSourceLine = []*model.RawSourceLine{}
	var yearSet map[int]bool = make(map[int]bool)

	t := time.Now()
	for {
		var rows []string
		if rows, err = f.Next(); err == io.EOF {
			fmt.Printf("Error Reading CSV row %+v\n", err)
			break
		}
		if err == nil && len(rows) > 0 {
			switch rows[0] {
			case "MeteringpointID":
				cpMeta = make(map[int]*model.CounterPointMeta, len(rows)-1)
				for i, c := range rows[1:] {
					id := fmt.Sprintf("%.3d", i)
					cpMeta[i] = &model.CounterPointMeta{ID: id, Name: c, Idx: i}
				}
			case "Energy direction":
				for i, c := range rows[1:] {
					switch cpMeta[i].Dir = c; cpMeta[i].Dir {
					case "CONSUMPTION":
						countConsumption += 1
					case "GENERATION":
						countGeneration += 1
					}
				}
			default:
				switch {
				case dateLine.MatchString(rows[0]):
					var y, m, d, hh, mm, ss int
					rawData := &model.RawSourceLine{Consumers: []float64{}, Producers: []float64{}}
					if _, err := fmt.Sscanf(rows[0], "%d.%d.%d %d:%d:%d", &d, &m, &y, &hh, &mm, &ss); err == nil {
						rawData.Id = fmt.Sprintf("CP/%d/%.2d/%.2d/%.2d/%.2d/%.2d", y, m, d, hh, mm, ss)
						yearSet[y] = true
					} else {
						fmt.Printf("Error Time parsing: %s (%s)", err, rows[0])
					}
					for i, c := range rows[1:] {
						if len(c) == 0 {
							c = "0"
						}
						switch cpMeta[i].Dir {
						case "CONSUMPTION":
							if d, err := strconv.ParseFloat(strings.Replace(c, ",", ".", -1), 32); err == nil {
								rawData.Consumers = append(rawData.Consumers, float64(d))
								cpMeta[i].Idx = len(rawData.Consumers) - 1
								cpMeta[i].Count += 1
							}
						case "GENERATION":
							if d, err := strconv.ParseFloat(strings.Replace(c, ",", ".", -1), 32); err == nil {
								rawData.Producers = append(rawData.Producers, float64(d))
								cpMeta[i].Idx = len(rawData.Producers) - 1
								cpMeta[i].Count += 1
							}
						}
					}
					rawDatas = append(rawDatas, rawData)
					//if err = db.SetLine(rawData); err != nil {
					//	_ = fmt.Errorf("%v\n", err)
					//}

					//fmt.Printf("(%d) RawData: %+v\n", rawData, rIdx)
					//fmt.Printf("RawData: %d\n", rIdx)
					rIdx += 1
				}
			}
		}
	}
	fmt.Printf("Time taken via read file: %v\n", time.Since(t))
	if err := db.SetLines(rawDatas); err != nil {
		return map[int]bool{}, err
	}

	rawMeta := &model.RawSourceMeta{Id: fmt.Sprintf("cpmeta/%d", 0), CounterPoints: []*model.CounterPointMeta{}, NumberOfMetering: rIdx}
	for _, v := range cpMeta {
		rawMeta.CounterPoints = append(rawMeta.CounterPoints, v)
	}
	fmt.Printf("MetaData: %+v\n", rawMeta)
	err = db.SetMeta(rawMeta)
	fmt.Printf("Time taken via write batch: %v\n", time.Since(t))
	return yearSet, err
}
