package calculation

import (
	"at.ourproject/energystore/model"
	"at.ourproject/energystore/store"
	"at.ourproject/energystore/utils"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/glog"
	"sort"
	"time"
)

type MqttEnergyImporter struct {
	tenant string
}

func NewMqttEnergyImporter(tenant string) *MqttEnergyImporter {
	return &MqttEnergyImporter{tenant}
}

func (mw *MqttEnergyImporter) Execute(msg mqtt.Message) {
	data := decodeMessage(msg.Payload())
	if data == nil {
		return
	}
	fmt.Printf("Execute Energy Data Message for Topic (%v)\n", mw.tenant)
	err := importEnergy(mw.tenant, data)
	if err != nil {
		glog.Error(err)
	}
}

func decodeMessage(msg []byte) *model.MqttEnergyResponse {
	m := model.MqttEnergyResponse{}
	err := json.Unmarshal(msg, &m)
	if err != nil {
		println(err.Error())
		return nil
	}
	return &m
}

func importEnergy(tenant string, data *model.MqttEnergyResponse) error {
	// GetMetaData from tenant

	db, err := store.OpenStorage(tenant)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	defaultDirection := utils.DetermineDirection(data.Message.Meter.MeteringPoint)

	var consumerCount int
	var producerCount int
	var metaCP *model.CounterPointMeta

	determineMeta := func() error {
		meta, info, err := store.GetMetaInfoMap(db, data.Message.Meter.MeteringPoint, defaultDirection)
		if err != nil {
			return err
		}

		consumerCount = info.ConsumerCount
		producerCount = info.ProducerCount

		metaCP = meta[data.Message.Meter.MeteringPoint]
		return nil
	}

	//// GetRawDataStructur from Period xxxx -> yyyy
	var _line model.RawSourceLine
	var resources map[string]*model.RawSourceLine = map[string]*model.RawSourceLine{}

	begin := time.UnixMilli(data.Message.Energy.Start)
	end := time.UnixMilli(data.Message.Energy.End)

	year, duration := utils.GetMonthDuration(begin, end)
	month := int(begin.Month())
	duration += 1

	if err := determineMeta(); err != nil {
		return err
	}

	y := year
	years := []int{year}
	for i := 0; i < duration; i++ {
		key := fmt.Sprintf("CP/%.4d/%.2d", y, ((month+i-1)%12)+1)
		iter := db.GetLinePrefix(key)
		for iter.Next(&_line) {
			l := _line.Copy(len(_line.Consumers))
			resources[_line.Id] = &l
		}
		if (month+i)%12 == 0 {
			y += 1
			years = append(years, y)
		}
	}
	// Update RawDataStructure
	glog.Infof("Len Loaded Resources %d", len(resources))

	sort.Slice(data.Message.Energy.Data[0].Value, func(i, j int) bool {
		a := time.UnixMilli(data.Message.Energy.Data[0].Value[i].From)
		b := time.UnixMilli(data.Message.Energy.Data[0].Value[j].From)
		return a.Unix() < b.Unix()
	})
	updated := []*model.RawSourceLine{}
	for _, v := range data.Message.Energy.Data[0].Value {
		t := time.UnixMilli(v.From)
		if t.Year() != year {
			year = t.Year()
			if err = determineMeta(); err != nil {
				return err
			}
		}

		id, err := utils.ConvertUnixTimeToRowId("CP/", time.UnixMilli(v.From))
		if err != nil {
			return err
		}
		_, ok := resources[id]
		if !ok {
			resources[id] = &model.RawSourceLine{Id: id, Consumers: make([]float64, consumerCount), Producers: make([]float64, producerCount)}
		}

		switch metaCP.Dir {
		case model.CONSUMER_DIRECTION:
			resources[id].Consumers = utils.Insert(resources[id].Consumers, metaCP.SourceIdx, v.Value)
		case model.PRODUCER_DIRECTION:
			resources[id].Producers = utils.Insert(resources[id].Producers, metaCP.SourceIdx, v.Value)
		}
		updated = append(updated, resources[id])
	}

	// Store updated RawDataStructure
	glog.Infof("Update CP %s energy values (%d) from %s to %s",
		data.Message.Meter.MeteringPoint,
		len(updated),
		time.UnixMilli(data.Message.Energy.Start).Format(time.RFC822),
		time.UnixMilli(data.Message.Energy.End).Format(time.RFC822))
	err = db.SetLines(updated)
	if err != nil {
		return err
	}

	for _, y := range years {
		if err = CalculateMonthlyDash(db, fmt.Sprintf("%d", y), CalculateEEG); err != nil {
			return err
		}
	}
	return nil
}
