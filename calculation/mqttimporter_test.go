package calculation

import (
	"at.ourproject/energystore/model"
	"at.ourproject/energystore/store"
	"at.ourproject/energystore/utils"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

//var (
//	db *store.BowStorage
//)

//func TestMain(m *testing.M) {
//
//	setUp()
//
//	retCode := m.Run()
//
//	setDown()
//
//	os.Exit(retCode)
//}

//func setUp() {
//	viper.Set("persistence.path", "../../../rawdata/test")
//
//}
//
//func setDown() {
//	os.RemoveAll("../../../rawdata/test/importer")
//	fmt.Println("Dataset (rawdata/test/importer) removed!")
//}

func TestNewMqttEnergyImporter(t *testing.T) {
	timeV1, err := utils.ParseTime("24.10.2022 00:00:00", time.Now().UnixMilli())
	timeV2, err := utils.ParseTime("24.10.2022 00:15:00", time.Now().UnixMilli())
	require.NoError(t, err)
	tests := []struct {
		name     string
		energy   *model.MqttEnergyMessage
		expected func(t *testing.T, l *model.RawSourceLine)
	}{
		{
			name: "Insert New Energy Allocated",
			energy: &model.MqttEnergyMessage{
				EcId: "ecIdTest",
				Meter: model.EnergyMeter{
					MeteringPoint: "AT0030000000000000000000000000001",
					Direction:     "",
				},
				Energy: model.MqttEnergy{
					Start: timeV1.UnixMilli(),
					End:   timeV2.UnixMilli(),
					Data: []model.MqttEnergyData{
						model.MqttEnergyData{
							MeterCode: "1-1:1.9.0 G.01",
							Value: []model.MqttEnergyValue{
								model.MqttEnergyValue{
									From:   timeV1.UnixMilli(),
									To:     timeV2.UnixMilli(),
									Method: "",
									Value:  1.11,
								},
							},
						},
					},
				},
			},
			expected: func(t *testing.T, l *model.RawSourceLine) {
				require.Equal(t, 1, len(l.Consumers))
				assert.Equal(t, 1.11, l.Consumers[0])
			},
		},
		{
			name: "Second Energy Consumer",
			energy: &model.MqttEnergyMessage{
				EcId: "ecIdTest",
				Meter: model.EnergyMeter{
					MeteringPoint: "AT0030000000000000000000000000002",
					Direction:     "",
				},
				Energy: model.MqttEnergy{
					Start: timeV1.UnixMilli(),
					End:   timeV2.UnixMilli(),
					Data: []model.MqttEnergyData{
						model.MqttEnergyData{
							MeterCode: "1-1:1.9.0 G.01",
							Value: []model.MqttEnergyValue{
								model.MqttEnergyValue{
									From:   timeV1.UnixMilli(),
									To:     timeV2.UnixMilli(),
									Method: "",
									Value:  0.11,
								},
							},
						},
					},
				},
			},
			expected: func(t *testing.T, l *model.RawSourceLine) {
				require.Equal(t, 4, len(l.Consumers))
				assert.Equal(t, 1.11, l.Consumers[0])
				assert.Equal(t, 0.11, l.Consumers[3])
			},
		},
		{
			name: "Insert Generator energy values",
			energy: &model.MqttEnergyMessage{
				EcId: "ecIdTest",
				Meter: model.EnergyMeter{
					MeteringPoint: "AT0030000000000000000000030000011",
					Direction:     "",
				},
				Energy: model.MqttEnergy{
					Start: timeV1.UnixMilli(),
					End:   timeV2.UnixMilli(),
					Data: []model.MqttEnergyData{
						model.MqttEnergyData{
							MeterCode: "1-1:2.9.0 P.01",
							Value: []model.MqttEnergyValue{
								model.MqttEnergyValue{
									From:   timeV1.UnixMilli(),
									To:     timeV2.UnixMilli(),
									Method: "",
									Value:  0.11,
								},
							},
						},
						model.MqttEnergyData{
							MeterCode: "1-1:1.9.0 G.01",
							Value: []model.MqttEnergyValue{
								model.MqttEnergyValue{
									From:   timeV1.UnixMilli(),
									To:     timeV2.UnixMilli(),
									Method: "",
									Value:  10.1,
								},
							},
						},
					},
				},
			},
			expected: func(t *testing.T, l *model.RawSourceLine) {
				require.Equal(t, 2, len(l.Producers))
				assert.Equal(t, 10.1, l.Producers[0])
			},
		},
		{
			name: "Insert second Generator Allocated",
			energy: &model.MqttEnergyMessage{
				EcId: "ecIdTest",
				Meter: model.EnergyMeter{
					MeteringPoint: "AT0030000000000000000000030000010",
					Direction:     "",
				},
				Energy: model.MqttEnergy{
					Start: timeV1.UnixMilli(),
					End:   timeV2.UnixMilli(),
					Data: []model.MqttEnergyData{
						model.MqttEnergyData{
							MeterCode: model.CODE_PLUS,
							Value: []model.MqttEnergyValue{
								model.MqttEnergyValue{
									From:   timeV1.UnixMilli(),
									To:     timeV2.UnixMilli(),
									Method: "",
									Value:  20.1,
								},
							},
						},
						model.MqttEnergyData{
							MeterCode: model.CODE_GEN,
							Value: []model.MqttEnergyValue{
								model.MqttEnergyValue{
									From:   timeV1.UnixMilli(),
									To:     timeV2.UnixMilli(),
									Method: "",
									Value:  21.1,
								},
							},
						},
					},
				},
			},
			expected: func(t *testing.T, l *model.RawSourceLine) {
				require.Equal(t, 4, len(l.Producers))
				assert.Equal(t, 10.1, l.Producers[0])
				assert.Equal(t, 21.1, l.Producers[2])
				assert.Equal(t, 20.1, l.Producers[3])
			},
		},
		{
			name: "Insert Generator - summarize energy values",
			energy: &model.MqttEnergyMessage{
				EcId: "ecIdTest",
				Meter: model.EnergyMeter{
					MeteringPoint: "AT0030000000000000000000030000010",
					Direction:     "",
				},
				Energy: model.MqttEnergy{
					Start: timeV1.UnixMilli(),
					End:   timeV2.UnixMilli(),
					Data: []model.MqttEnergyData{
						model.MqttEnergyData{
							MeterCode: model.CODE_PLUS,
							Value: []model.MqttEnergyValue{
								model.MqttEnergyValue{
									From:   timeV1.UnixMilli(),
									To:     timeV2.UnixMilli(),
									Method: "L1",
									Value:  20.1,
								},
								model.MqttEnergyValue{
									From:   timeV1.UnixMilli(),
									To:     timeV2.UnixMilli(),
									Method: "L1",
									Value:  10.1,
								},
							},
						},
						model.MqttEnergyData{
							MeterCode: model.CODE_GEN,
							Value: []model.MqttEnergyValue{
								model.MqttEnergyValue{
									From:   timeV1.UnixMilli(),
									To:     timeV2.UnixMilli(),
									Method: "L1",
									Value:  5.1,
								},
								model.MqttEnergyValue{
									From:   timeV1.UnixMilli(),
									To:     timeV2.UnixMilli(),
									Method: "L1",
									Value:  2.2,
								},
							},
						},
					},
				},
			},
			expected: func(t *testing.T, l *model.RawSourceLine) {
				fmt.Printf("Producer Line: %+v\n", l)
				require.Equal(t, 4, len(l.Producers))
				assert.Equal(t, 7.3, utils.RoundToFixed(l.Producers[2], 1))
				assert.Equal(t, 1, l.QoVProducers[2])
				assert.Equal(t, 30.2, utils.RoundToFixed(l.Producers[3], 1))
				assert.Equal(t, 1, l.QoVProducers[3])
			},
		},
	}

	viper.Set("persistence.path", "../test/rawdata")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = importEnergyV2("importer", "ecid", tt.energy)
			require.NoError(t, err)

			db, err := store.OpenStorageTest("importer", "ecid", "../test/rawdata")
			require.NoError(t, err)
			it := db.GetLinePrefix(fmt.Sprintf("CP/%s", "2022/10/24"))
			defer it.Close()
			defer db.Close()

			var _line model.RawSourceLine

			r := it.Next(&_line)
			assert.Equal(t, true, r)
			assert.Equal(t, "CP/2022/10/24/00/00/00", _line.Id)
			tt.expected(t, &_line)
		})
	}

	os.RemoveAll("../test/rawdata/importer")
}

func TestImportRawdataStore(t *testing.T) {

	viper.Set("persistence.path", "../test/rawdata")

	jsonRaw, err := os.ReadFile("../test/energy-response-text.json")
	require.NoError(t, err)

	rawData := decodeMessage(jsonRaw)
	require.NotNil(t, rawData)

	err = importEnergyV2("te100190", "ecid", rawData)
	require.NoError(t, err)

	rawData.Meter.MeteringPoint = "AT0030000000000000000000000381702"
	err = importEnergyV2("te100190", "ecid", rawData)
	require.NoError(t, err)

	db, err := store.OpenStorageTest("te100190", "ecid", "../test/rawdata")
	require.NoError(t, err)

	meta, err := db.GetMeta("cpmeta/0")
	for i, v := range meta.CounterPoints {
		fmt.Printf("[%d]: %+v\n", i, v)
	}

	it := db.GetLinePrefix("CP/")

	line := model.RawSourceLine{}
	lines := []*model.RawSourceLine{}
	for it.Next(&line) {
		_line := line.Copy(len(line.Consumers))
		lines = append(lines, &_line)
	}
	it.Close()
	db.Close()

	require.Equal(t, 23*4, len(lines)) // one hour is missing from the test source file

	energy, err := EnergyReport("te100190", 2023, 3, "YM")
	require.NoError(t, err)

	response, err := json.Marshal(energy)
	require.NoError(t, err)

	require.Equal(t, 3, len(energy.Report.Allocated))
	require.Equal(t, 1.088021, energy.Report.Allocated[0])
	require.Equal(t, 3, len(energy.Report.Consumed))
	require.Equal(t, 5.388, energy.Report.Consumed[0])

	fmt.Printf("META_DATA: %+v\n", string(response))

	os.RemoveAll("../test/rawdata/rc100190")
}

func Test_updateMetaCP(t *testing.T) {
	type args struct {
		metaCP *model.CounterPointMeta
		begin  time.Time
		end    time.Time
	}
	tests := []struct {
		name string
		args args
		test func(t *testing.T, args args)
	}{
		{
			name: "Adjust meta period end time",
			args: args{
				metaCP: &model.CounterPointMeta{
					ID:          "000",
					Name:        "IV0000999222222222221",
					SourceIdx:   0,
					Dir:         "CONSUMPTION",
					Count:       0,
					PeriodStart: "30.12.2023 15:00:0000",
					PeriodEnd:   "30.12.2023 15:00:0000",
				},
				begin: time.Date(2023, 12, 30, 15, 1, 0, 0, time.Local),
				end:   time.Date(2023, 12, 30, 15, 15, 0, 0, time.Local),
			},
			test: func(t *testing.T, args args) {
				result := updateMetaCP(args.metaCP, args.begin, args.end)
				assert.Equalf(t, true, result, "updateMetaCP(%v, %v, %v)", args.metaCP, args.begin, args.end)
				assert.Equal(t, "30.12.2023 15:00:0000", args.metaCP.PeriodStart)
				assert.Equal(t, "30.12.2023 15:15:0000", args.metaCP.PeriodEnd)
			},
		},
		{
			name: "Adjust meta period start time",
			args: args{
				metaCP: &model.CounterPointMeta{
					ID:          "000",
					Name:        "IV0000999222222222221",
					SourceIdx:   0,
					Dir:         "CONSUMPTION",
					Count:       0,
					PeriodStart: "30.12.2023 15:00:0000",
					PeriodEnd:   "30.12.2023 15:15:0000",
				},
				begin: time.Date(2023, 12, 30, 14, 0, 0, 0, time.Local),
				end:   time.Date(2023, 12, 30, 15, 15, 0, 0, time.Local),
			},
			test: func(t *testing.T, args args) {
				result := updateMetaCP(args.metaCP, args.begin, args.end)
				assert.Equalf(t, true, result, "updateMetaCP(%v, %v, %v)", args.metaCP, args.begin, args.end)
				assert.Equal(t, "30.12.2023 14:00:0000", args.metaCP.PeriodStart)
				assert.Equal(t, "30.12.2023 15:15:0000", args.metaCP.PeriodEnd)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t, tt.args)
		})
	}
}
