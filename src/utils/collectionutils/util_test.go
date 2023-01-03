package collectionutils

import (
	"fmt"
	"reflect"
	"testing"
)

type MockCar struct {
	Model, Brand, Color string
}

func TestGetChunkWhenChunkIsLessThanLengthOfSlice(t *testing.T) {
	mockCars := &[]MockCar{
		{
			Model: "Egea",
			Brand: "Fiat",
			Color: "Blue",
		},
		{
			Model: "A3",
			Brand: "Audi",
			Color: "White",
		},
		{
			Model: "A4",
			Brand: "Audi",
			Color: "White",
		},
		{
			Model: "A5",
			Brand: "Audi",
			Color: "White",
		},
		{
			Model: "A6",
			Brand: "Audi",
			Color: "White",
		},
	}
	chunkSize := 2
	firstIterationResult := (*mockCars)[0:2]
	secondIterationResult := (*mockCars)[2:4]
	thirdIterationResult := (*mockCars)[4:5]

	var tests = []struct {
		collection     *[]MockCar
		iterationIndex int
		want           *[]MockCar
	}{
		{mockCars, 0, &firstIterationResult},
		{mockCars, 1, &secondIterationResult},
		{mockCars, 2, &thirdIterationResult},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("For %d iteration index", tt.iterationIndex)
		t.Run(testName, func(t *testing.T) {
			got := GetChunk(tt.collection, chunkSize, tt.iterationIndex)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestGetChunkWhenChunkIsGreaterThanLengthOfSlice(t *testing.T) {
	mockCars := &[]MockCar{
		{
			Model: "Egea",
			Brand: "Fiat",
			Color: "Blue",
		},
		{
			Model: "A3",
			Brand: "Audi",
			Color: "White",
		},
		{
			Model: "A4",
			Brand: "Audi",
			Color: "White",
		},
		{
			Model: "A5",
			Brand: "Audi",
			Color: "White",
		},
		{
			Model: "A6",
			Brand: "Audi",
			Color: "White",
		},
	}
	chunkSize := 6
	firstIterationResult := (*mockCars)[0:len(*mockCars)]

	testName := fmt.Sprintf("For %d iteration index", 0)
	t.Run(testName, func(t *testing.T) {
		got := GetChunk(mockCars, chunkSize, 0)
		if !reflect.DeepEqual(got, &firstIterationResult) {
			t.Errorf("got %s, want %s", got, &firstIterationResult)
		}
	})
}
