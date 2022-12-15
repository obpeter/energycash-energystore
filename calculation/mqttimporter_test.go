package calculation

import (
	"at.ourproject/energystore/model"
	"at.ourproject/energystore/store"
	"at.ourproject/energystore/utils"
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
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
	timeV1, err := utils.ParseTime("24.10.2022 00:00:00")
	timeV2, err := utils.ParseTime("24.10.2022 00:15:00")
	require.NoError(t, err)
	tests := []struct {
		name     string
		energy   *model.MqttEnergyResponse
		expected func(t *testing.T, l *model.RawSourceLine)
	}{
		{
			name: "Insert New Energy Allocated",
			energy: &model.MqttEnergyResponse{
				Message: model.MqttEnergyMessage{
					Meter: model.EnergyMeter{
						MeteringPoint: "AT0030000000000000000000000000001",
						Direction:     "",
					},
					Energy: model.MqttEnergy{
						Start: timeV1.UnixMilli(),
						End:   timeV2.UnixMilli(),
						Data: []model.MqttEnergyData{
							model.MqttEnergyData{
								MeterCode: "G1",
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
			},
			expected: func(t *testing.T, l *model.RawSourceLine) {
				require.Equal(t, 1, len(l.Consumers))
				assert.Equal(t, 1.11, l.Consumers[0])
			},
		},
		{
			name: "Second Energy Consumer",
			energy: &model.MqttEnergyResponse{
				Message: model.MqttEnergyMessage{
					Meter: model.EnergyMeter{
						MeteringPoint: "AT0030000000000000000000000000002",
						Direction:     "",
					},
					Energy: model.MqttEnergy{
						Start: timeV1.UnixMilli(),
						End:   timeV2.UnixMilli(),
						Data: []model.MqttEnergyData{
							model.MqttEnergyData{
								MeterCode: "G1",
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
			},
			expected: func(t *testing.T, l *model.RawSourceLine) {
				require.Equal(t, 2, len(l.Consumers))
				assert.Equal(t, 0.11, l.Consumers[1])
			},
		},
		{
			name: "Insert Generator energy values",
			energy: &model.MqttEnergyResponse{
				Message: model.MqttEnergyMessage{
					Meter: model.EnergyMeter{
						MeteringPoint: "AT0030000000000000000000030000011",
						Direction:     "",
					},
					Energy: model.MqttEnergy{
						Start: timeV1.UnixMilli(),
						End:   timeV2.UnixMilli(),
						Data: []model.MqttEnergyData{
							model.MqttEnergyData{
								MeterCode: "G1",
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
			},
			expected: func(t *testing.T, l *model.RawSourceLine) {
				require.Equal(t, 2, len(l.Producers))
				assert.Equal(t, 10.1, l.Producers[0])
			},
		},
		{
			name: "Insert second Generator Allocated",
			energy: &model.MqttEnergyResponse{
				Message: model.MqttEnergyMessage{
					Meter: model.EnergyMeter{
						MeteringPoint: "AT0030000000000000000000030000010",
						Direction:     "",
					},
					Energy: model.MqttEnergy{
						Start: timeV1.UnixMilli(),
						End:   timeV2.UnixMilli(),
						Data: []model.MqttEnergyData{
							model.MqttEnergyData{
								MeterCode: "G1",
								Value: []model.MqttEnergyValue{
									model.MqttEnergyValue{
										From:   timeV1.UnixMilli(),
										To:     timeV2.UnixMilli(),
										Method: "",
										Value:  20.1,
									},
								},
							},
						},
					},
				},
			},
			expected: func(t *testing.T, l *model.RawSourceLine) {
				require.Equal(t, 2, len(l.Producers))
				assert.Equal(t, 20.1, l.Producers[1])
			},
		},
	}

	viper.Set("persistence.path", "../test/rawdata")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = importEnergy("importer", tt.energy)
			require.NoError(t, err)

			db, err := store.OpenStorageTest("importer", "../test/rawdata")
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
