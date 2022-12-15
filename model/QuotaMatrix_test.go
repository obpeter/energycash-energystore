package model

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

/*
	 Matrix Quota

	       |-P1|-P2|-P3|
		ZP1|0.4|0.5|0.1|
		ZP2|0.1|0.3|0.4|
		ZP3|0.5|0.2|0.3|
*/
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestMakeQuotaMatrix(t *testing.T) {
	cpNames := []string{"ZP1", "ZP2", "ZP3"}
	prodNames := []string{"P1", "P2", "P3"}

	qoutamatrix := NewQuotaMatrix(cpNames, prodNames)

	qoutamatrix.Add("ZP1", "P1", 0.4)
	qoutamatrix.Add("ZP1", "P2", 0.5)
	qoutamatrix.Add("ZP1", "P3", 0.1)
	qoutamatrix.Add("ZP2", "P1", 0.1)
	qoutamatrix.Add("ZP2", "P2", 0.3)
	qoutamatrix.Add("ZP2", "P3", 0.45)
	qoutamatrix.Add("ZP3", "P1", 0.5)
	qoutamatrix.Add("ZP3", "P2", 0.2)
	qoutamatrix.Add("ZP3", "P3", 0.3)

	v := qoutamatrix.GetQuota("ZP2", "P3")
	println(v)
	assert.Equal(t, float64(0.45), v)

	v = qoutamatrix.GetAllocQuota("ZP2", "P3")
	println(v)
	assert.Equal(t, float64(0.5294117647058822), v)

	ts := [7]int{1}

	fmt.Printf("Test Slice: %+v\n", ts)
	qoutamatrix.Validate()
}

func TestMerge(t *testing.T) {
	m1 := MakeMatrix([]float64{1, 2, 5, 6}, 2, 2)
	m2 := MakeMatrix([]float64{3, 4, 7, 8}, 2, 2)

	m3 := Merge(m1, m2)

	fmt.Printf("TM-Matrix: %+v\n", m3)
	assert.Equal(t, 4, m3.CountCols())
	assert.Equal(t, 2, m3.CountRows())

	assert.Equal(t, float64(4), m3.GetElm(0, 3))
}

func TestMultiply(t *testing.T) {
	m1 := MakeMatrix([]float64{1, 2, 5, 6, 7, 8}, 3, 2)
	m2 := NewUniformMatrix(2, 1)

	m3 := Multiply(m1, m2)

	assert.Equal(t, 1, m3.CountCols())
	assert.Equal(t, 3, m3.CountRows())

	assert.Equal(t, float64(3), m3.GetElm(0, 0))
	assert.Equal(t, float64(11), m3.GetElm(1, 0))
	assert.Equal(t, float64(15), m3.GetElm(2, 0))
}

func TestMatrix_Add_DifferentSize(t *testing.T) {

	tests := []struct {
		name      string
		matrixA   *Matrix
		matrixB   *Matrix
		exprected *Matrix
	}{
		{
			name:      "Different Size A > B",
			matrixA:   MakeMatrix([]float64{0, 0, 1, 0, 1, 0}, 6, 1),
			matrixB:   MakeMatrix([]float64{3, 4, 7, 8}, 4, 1),
			exprected: MakeMatrix([]float64{3, 4, 8, 8, 1, 0}, 6, 1),
		},
		{
			name:      "Different Size A < B",
			matrixA:   MakeMatrix([]float64{0, 0, 1, 0}, 4, 1),
			matrixB:   MakeMatrix([]float64{3, 4, 7, 8, 9, 1}, 6, 1),
			exprected: MakeMatrix([]float64{3, 4, 8, 8, 9, 1}, 6, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.matrixA.Add(tt.matrixB)
			assert.NoError(t, err)

			fmt.Printf("Result: %+v\n", tt.matrixA)
			assert.Equal(t, tt.exprected.CountRows(), tt.matrixA.CountRows())
			assert.Equal(t, tt.exprected.Cols, tt.matrixA.Cols)
			assert.Equal(t, tt.exprected.Elements, tt.matrixA.Elements)
		})
	}
	//m1 := MakeMatrix([]float64{0, 0, 1, 0, 1, 0}, 6, 1)
	//m2 := MakeMatrix([]float64{3, 4, 7, 8}, 4, 1)
	//
	//err := m1.Add(m2)
	//assert.NoError(t, err)
	//
	//fmt.Printf("Result: %+v\n", m1)
	//assert.Equal(t, 6, m1.CountRows())
	//assert.Equal(t, []float64{3, 4, 8, 8, 1, 0}, m1.Elements)
}
