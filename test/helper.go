package test

import (
	"at.ourproject/energystore/excel"
	"at.ourproject/energystore/store"
	"github.com/stretchr/testify/require"
	"testing"
)

func ImportTestContent(t *testing.T, db *store.BowStorage) (yearSet []int) {
	excelFile, err := excel.OpenExceFile("./zaehlpunkte-beispieldatei.xlsx")
	require.NoError(t, err)
	defer excelFile.Close()

	yearSet, err = excel.ImportExcelEnergyFile(excelFile, "ConsumptionDataReport", db)
	require.NoError(t, err)

	//for _, k := range yearSet {
	//	err = calculation.CalculateMonthlyDash(db, fmt.Sprintf("%d", k), calculation.CalculateEEG)
	//	require.NoError(t, err)
	//}
	return
}
