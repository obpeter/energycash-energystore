package excel

import (
	"at.ourproject/energystore/calculation"
	"at.ourproject/energystore/model"
	"at.ourproject/energystore/store"
	"fmt"
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
//
//func setUp() {
//	println("Setup TEST!!!!!!!!!!!")
//	var err error
//	db, err = store.OpenStorageTest("excelsource", "../../../rawdata/test")
//	if err != nil {
//		panic(err)
//	}
//}
//
//func setDown() {
//	_ = db.Close()
//	os.RemoveAll("../../../rawdata/test/excelsource")
//
//}

func TestImportExcelEnergyFile(t *testing.T) {
	db, err := store.OpenStorageTest("dashboard", "../test/rawdata")
	require.Nil(t, err)
	defer func() {
		_ = db.Close()
		os.RemoveAll("../test/rawdata/dashboard")
	}()

	excelFile, err := OpenExceFile("../test/zaehlpunkte-beispieldatei.xlsx")
	require.NoError(t, err)
	defer excelFile.Close()

	yearSet, err := ImportExcelEnergyFile(excelFile, "ConsumptionDataReport", db)
	require.NoError(t, err)
	for _, k := range yearSet {
		err = calculation.CalculateMonthlyDash(db, fmt.Sprintf("%d", k), calculation.CalculateEEG)
		require.NoError(t, err)
	}
}

func TestBuildMatixMetaStruct(t *testing.T) {

	db, err := store.OpenStorageTest("excelsource", "../test/rawdata")
	require.NoError(t, err)
	defer func() {
		db.Close()
		os.RemoveAll("../test/rawdata/excelsource")
	}()

	t.Run("Initiate Meta Struckt", func(t *testing.T) {
		header := excelHeader{
			meteringPointId: map[int]string{
				0:  "AT0030000000000000000000000000001",
				1:  "AT0030000000000000000000000000001",
				2:  "AT0030000000000000000000000000001",
				3:  "AT0030000000000000000000000000002",
				4:  "AT0030000000000000000000000000002",
				5:  "AT0030000000000000000000000000002",
				6:  "AT0030000000000000000000000000003",
				7:  "AT0030000000000000000000000000003",
				8:  "AT0030000000000000000000000000003",
				9:  "AT0030000000000000000000000000004",
				10: "AT0030000000000000000000000000005",
				11: "AT0030000000000000000000000000006",
				12: "AT0030000000000000000000000000007",
				13: "AT0030000000000000000000000000008",
				14: "AT0030000000000000000000000000009",
				15: "AT0030000000000000000000000000010",
				16: "AT0030000000000000000000000000011",
				17: "AT0030000000000000000000000000012",
				18: "AT0030000000000000000000000000013",
				19: "TOTAL",
				20: "TOTAL",
				21: "TOTAL",
			},
			energyDirection: map[int]string{
				0:  "CONSUMPTION",
				1:  "CONSUMPTION",
				2:  "CONSUMPTION",
				3:  "CONSUMPTION",
				4:  "CONSUMPTION",
				5:  "CONSUMPTION",
				6:  "CONSUMPTION",
				7:  "CONSUMPTION",
				8:  "CONSUMPTION",
				9:  "CONSUMPTION",
				10: "CONSUMPTION",
				11: "CONSUMPTION",
				12: "CONSUMPTION",
				13: "CONSUMPTION",
				14: "CONSUMPTION",
				15: "CONSUMPTION",
				16: "CONSUMPTION",
				17: "CONSUMPTION",
				18: "GENERATION",
			},
			meterCode: map[int]MeterCodeType{
				0:  Total,
				1:  Share,
				2:  Coverage,
				3:  Total,
				4:  Share,
				5:  Coverage,
				6:  Total,
				7:  Share,
				8:  Coverage,
				9:  Total,
				10: Total,
				11: Total,
				12: Total,
				13: Total,
				14: Total,
				15: Total,
				16: Total,
				17: Total,
				18: Total,
				19: Total,
				20: Total,
				21: Total,
			},
			periodStart: map[int]string{
				0:  "24.10.2022 00:00:00",
				1:  "24.10.2022 00:00:00",
				2:  "24.10.2022 00:00:00",
				3:  "24.10.2022 00:00:00",
				4:  "24.10.2022 00:00:00",
				5:  "24.10.2022 00:00:00",
				6:  "24.10.2022 00:00:00",
				7:  "24.10.2022 00:00:00",
				8:  "24.10.2022 00:00:00",
				9:  "24.10.2022 00:00:00",
				10: "24.10.2022 00:00:00",
				11: "24.10.2022 00:00:00",
				12: "24.10.2022 00:00:00",
				13: "24.10.2022 00:00:00",
				14: "24.10.2022 00:00:00",
				15: "24.10.2022 00:00:00",
				16: "24.10.2022 00:00:00",
				17: "24.10.2022 00:00:00",
				18: "24.10.2022 00:00:00",
			},
			periodEnd: map[int]string{
				0:  "24.10.2022 00:00:00",
				1:  "24.10.2022 00:00:00",
				2:  "24.10.2022 00:00:00",
				3:  "24.10.2022 00:00:00",
				4:  "24.10.2022 00:00:00",
				5:  "24.10.2022 00:00:00",
				6:  "24.10.2022 00:00:00",
				7:  "24.10.2022 00:00:00",
				8:  "24.10.2022 00:00:00",
				9:  "24.10.2022 00:00:00",
				10: "24.10.2022 00:00:00",
				11: "24.10.2022 00:00:00",
				12: "24.10.2022 00:00:00",
				13: "24.10.2022 00:00:00",
				14: "24.10.2022 00:00:00",
				15: "24.10.2022 00:00:00",
				16: "24.10.2022 00:00:00",
				17: "24.10.2022 00:00:00",
				18: "24.10.2022 00:00:00",
			},
		}

		excelCpMeta, cpMeta, err := buildMatrixMetaStruct(db, header)
		require.NoError(t, err)
		require.NotNil(t, cpMeta)
		require.Equal(t, 13, len(cpMeta))
		require.NotNil(t, excelCpMeta)
		require.Equal(t, 13, len(excelCpMeta))

		cConsuption := 0
		cGeneration := 0

		//expectedIdx := []int{0, 3, 6, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18}
		//expectedSourceIdx := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 0}
		currentIdx := []int{}
		currentSourceIdx := []int{}
		for j := 0; j < len(cpMeta); j++ {
			v := cpMeta[j]
			if v.Dir == "CONSUMPTION" {
				cConsuption += 1
			} else if v.Dir == "GENERATION" {
				cGeneration += 1
			}
			currentIdx = append(currentIdx, v.Idx)
			currentSourceIdx = append(currentSourceIdx, v.SourceIdx)
			fmt.Printf("Meta: %+v (%+v)\n", v, excelCpMeta[j])
		}

		require.Equal(t, 1, cGeneration)
		require.Equal(t, 12, cConsuption)

		//sort.Ints(currentIdx)
		//require.ElementsMatch(t, expectedIdx, currentIdx)
		//require.ElementsMatch(t, expectedSourceIdx, currentSourceIdx)

		rawMeta := &model.RawSourceMeta{Id: fmt.Sprintf("cpmeta/%d", 0), CounterPoints: cpMeta, NumberOfMetering: 13}
		err = db.SetMeta(rawMeta)
	})

	t.Run("Initiate Meta Struckt", func(t *testing.T) {

		header := excelHeader{
			meteringPointId: map[int]string{
				0:  "AT0030000000000000000000000000015", // 0
				1:  "AT0030000000000000000000000000015",
				2:  "AT0030000000000000000000000000015",
				3:  "AT0030000000000000000000000000014", // 3
				4:  "AT0030000000000000000000000000014",
				5:  "AT0030000000000000000000000000014",
				6:  "AT0030000000000000000000000000002", // 6
				7:  "AT0030000000000000000000000000002",
				8:  "AT0030000000000000000000000000002",
				9:  "AT0030000000000000000000000000004", // 9
				10: "AT0030000000000000000000000000005",
				11: "AT0030000000000000000000000000006",
				12: "AT0030000000000000000000000000007",
				13: "AT0030000000000000000000000000008",
				14: "AT0030000000000000000000000000009",
				15: "AT0030000000000000000000000000010",
				16: "AT0030000000000000000000000000011",
				17: "AT0030000000000000000000000000012",
				18: "AT0030000000000000000000000000013",
				19: "AT0030000000000000000000000000013",
				20: "TOTAL",
				21: "TOTAL",
				22: "TOTAL",
			},
			energyDirection: map[int]string{
				0:  "CONSUMPTION",
				1:  "CONSUMPTION",
				2:  "CONSUMPTION",
				3:  "CONSUMPTION",
				4:  "CONSUMPTION",
				5:  "CONSUMPTION",
				6:  "CONSUMPTION",
				7:  "CONSUMPTION",
				8:  "CONSUMPTION",
				9:  "CONSUMPTION",
				10: "CONSUMPTION",
				11: "CONSUMPTION",
				12: "CONSUMPTION",
				13: "CONSUMPTION",
				14: "CONSUMPTION",
				15: "CONSUMPTION",
				16: "CONSUMPTION",
				17: "CONSUMPTION",
				18: "GENERATION",
				19: "GENERATION",
			},
			meterCode: map[int]MeterCodeType{
				0:  Total,
				1:  Share,
				2:  Coverage,
				3:  Total,
				4:  Share,
				5:  Coverage,
				6:  Total,
				7:  Share,
				8:  Coverage,
				9:  Total,
				10: Total,
				11: Total,
				12: Total,
				13: Total,
				14: Total,
				15: Total,
				16: Total,
				17: Total,
				18: Bad,
				19: Total,
			},
			periodStart: map[int]string{
				0:  "25.10.2022 00:00:00",
				1:  "25.10.2022 00:00:00",
				2:  "25.10.2022 00:00:00",
				3:  "25.10.2022 00:00:00",
				4:  "25.10.2022 00:00:00",
				5:  "25.10.2022 00:00:00",
				6:  "24.10.2022 00:00:00",
				7:  "24.10.2022 00:00:00",
				8:  "24.10.2022 00:00:00",
				9:  "24.10.2022 00:00:00",
				10: "24.10.2022 00:00:00",
				11: "24.10.2022 00:00:00",
				12: "24.10.2022 00:00:00",
				13: "24.10.2022 00:00:00",
				14: "24.10.2022 00:00:00",
				15: "24.10.2022 00:00:00",
				16: "24.10.2022 00:00:00",
				17: "24.10.2022 00:00:00",
				18: "",
				19: "24.10.2022 00:00:00",
			},
			periodEnd: map[int]string{
				0:  "25.10.2022 00:00:00",
				1:  "25.10.2022 00:00:00",
				2:  "25.10.2022 00:00:00",
				3:  "25.10.2022 00:00:00",
				4:  "25.10.2022 00:00:00",
				5:  "25.10.2022 00:00:00",
				6:  "24.10.2022 00:00:00",
				7:  "24.10.2022 00:00:00",
				8:  "24.10.2022 00:00:00",
				9:  "24.10.2022 00:00:00",
				10: "24.10.2022 00:00:00",
				11: "25.10.2022 00:00:00",
				12: "24.10.2022 00:00:00",
				13: "24.10.2022 00:00:00",
				14: "25.10.2022 00:00:00",
				15: "24.10.2022 00:00:00",
				16: "24.10.2022 00:00:00",
				17: "24.10.2022 00:00:00",
				18: "24.10.2022 00:00:00",
			},
		}

		excelCpMeta, cpMeta, err := buildMatrixMetaStruct(db, header)
		require.NoError(t, err)
		require.NotNil(t, cpMeta)
		require.Equal(t, 15, len(cpMeta))
		require.NotNil(t, excelCpMeta)
		require.Equal(t, 13, len(excelCpMeta))

		cConsuption := 0
		cGeneration := 0
		var currentIdx []int
		currentSourceIdx := []int{}
		for j := 0; j < len(excelCpMeta); j++ {
			v := excelCpMeta[j]
			if v.Dir == "CONSUMPTION" {
				cConsuption += 1
			} else if v.Dir == "GENERATION" {
				cGeneration += 1
			}
			currentIdx = append(currentIdx, v.Idx)
			currentSourceIdx = append(currentSourceIdx, v.SourceIdx)
			fmt.Printf("Meta: %+v\n", v)
		}

		require.Equal(t, 1, cGeneration)
		require.Equal(t, 12, cConsuption)

		fmt.Printf("currentIdx: (%+v)\n", currentIdx)
		fmt.Printf("currentSourceIdx: (%+v)\n", currentSourceIdx)

		expectedIdx := []int{0, 3, 6, 9, 10, 11, 12, 13, 14, 15, 16, 17, 19}
		require.ElementsMatch(t, expectedIdx, currentIdx)
		expectedSourceIdx := []int{12, 13, 1, 3, 4, 5, 6, 7, 8, 9, 10, 11, 0}
		require.ElementsMatch(t, expectedSourceIdx, currentSourceIdx)
	})
}
