package model

type RawSourceLine struct {
	Id        string    `bow:"key"`
	Consumers []float64 `bow:"consumers"`
	Producers []float64 `bow:"producers"`
}

func (c RawSourceLine) Copy(cLength int) RawSourceLine {
	r := RawSourceLine{Id: c.Id, Consumers: make([]float64, cLength), Producers: make([]float64, len(c.Producers))}
	copy(r.Consumers[:], c.Consumers[:])
	copy(r.Producers[:], c.Producers[:])
	return r
}

type EegEnergy struct {
	Report  *EnergyReport       `json:"report"`
	Results []*EnergyReport     `json:"intermediateReportResults"`
	Meta    []*CounterPointMeta `json:"meta"`
}

type RawSourceMeta struct {
	Id               string              `bow:"key"`
	CounterPoints    []*CounterPointMeta `bow:"counterpoints"`
	NumberOfMetering int                 `bow:"numberOfMetering"`
}

func (rsm RawSourceMeta) Copy() RawSourceMeta {
	r := RawSourceMeta{Id: rsm.Id, CounterPoints: []*CounterPointMeta{}, NumberOfMetering: rsm.NumberOfMetering}
	for i := 0; i < len(rsm.CounterPoints); i++ {
		r.CounterPoints = append(r.CounterPoints,
			&CounterPointMeta{ID: rsm.CounterPoints[i].ID,
				Name:        rsm.CounterPoints[i].Name,
				Idx:         rsm.CounterPoints[i].Idx,
				Dir:         rsm.CounterPoints[i].Dir,
				Count:       rsm.CounterPoints[i].Count,
				PeriodStart: rsm.CounterPoints[i].PeriodStart,
				PeriodEnd:   rsm.CounterPoints[i].PeriodEnd,
			})
	}
	return r
}

/*
CounterPointMeta describe the raw data source in db. Be aware, the meta data are just applicable for one year.
*/
type CounterPointMeta struct {
	ID          string `bow:"key" json:"id"`                   // Usually, the Id represents Table-Type and Year (e.g. cpmeta/2021)
	Name        string `bow:"name" json:"name"`                // Counterpoint name (e.g. ZPxxxxxxxxZaehlpunktxx)
	Idx         int    `bow:"idx" json:"idx"`                  // Index of rawdata array in excel (Just needed for manual excel/csv import)
	SourceIdx   int    `bow:"srcIdx" json:"sourceIdx"`         // Index of rawdata array in db
	Dir         string `bow:"dir" json:"dir"`                  // Direction of consumption (GENERATOR or. CONSUMER)
	Count       uint16 `bow:"count" json:"count"`              // Number of measurements
	PeriodStart string `bow:"periodstart" json:"period_start"` // Period Start
	PeriodEnd   string `bow:"periodend" json:"period_end"`     // Period End
}

type CounterPointMetaInfo struct {
	ConsumerCount  int // Max. number of consumer in the CounterPointMeta data structure
	ProducerCount  int // Max. number of producer in the CounterPointMeta data structure
	MaxConsumerIdx int // Highest index of consumer in the CounterPointMeta data structure
	MaxProducerIdx int // Highest index of producer in the CounterPointMeta data structure
}

type ByReportDate []EnergyReport

func (a ByReportDate) Len() int           { return len(a) }
func (a ByReportDate) Less(i, j int) bool { return a[i].Id < a[j].Id }
func (a ByReportDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type EnergyReport struct {
	Id            string    `bow:"key" json:"id"`
	Allocated     []float64 `bow:"values" json:"allocated"`
	Consumed      []float64 `bow:"consumed" json:"consumed"`
	Produced      []float64 `bow:"produced" json:"produced"`
	TotalProduced float64   `bow:"totalProduced" json:"total_produced"`
}
