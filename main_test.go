package main

import (
	"path"
	"reflect"
	"testing"
	"time"
)

type WantFindTrains struct {
	Trains
	error
}

func TestFindTrains(t *testing.T) {
	dep1arr2 := Trains{
		Train{
			TrainID:            1,
			DepartureStationID: 1,
			ArrivalStationID:   2,
			Price:              1,
			ArrivalTime:        parseTime("00:00:01"),
			DepartureTime:      parseTime("00:00:01"),
		},
		Train{
			TrainID:            6,
			DepartureStationID: 1,
			ArrivalStationID:   2,
			Price:              6,
			ArrivalTime:        parseTime("00:00:06"),
			DepartureTime:      parseTime("00:00:06"),
		},
		Train{
			TrainID:            7,
			DepartureStationID: 1,
			ArrivalStationID:   2,
			Price:              7,
			ArrivalTime:        parseTime("00:00:07"),
			DepartureTime:      parseTime("00:00:07"),
		},
	}
	tests := []struct {
		input []string
		want  WantFindTrains
	}{
		{[]string{"", "1", "price"}, WantFindTrains{nil, ErrEmptyDep}},
		{[]string{"1", "", "price"}, WantFindTrains{nil, ErrEmptyArr}},
		{[]string{"bad", "1", "price"}, WantFindTrains{nil, ErrBadDep}},
		{[]string{"1", "bad", "price"}, WantFindTrains{nil, ErrBadArr}},
		{[]string{"1", "1", "bad"}, WantFindTrains{nil, ErrBadBy}},
		{[]string{"1", "1", "price"}, WantFindTrains{nil, ErrNoTrains}},
		{[]string{"1", "2", "price"}, WantFindTrains{dep1arr2, nil}},
		{[]string{"1", "2", "arrival-time"}, WantFindTrains{dep1arr2, nil}},
		{[]string{"1", "2", "departure-time"}, WantFindTrains{dep1arr2, nil}},
		{
			[]string{"2", "1", "departure-time"},
			WantFindTrains{
				Trains{
					Train{
						TrainID:            9,
						DepartureStationID: 2,
						ArrivalStationID:   1,
						Price:              9,
						ArrivalTime:        parseTime("00:00:09"),
						DepartureTime:      parseTime("00:00:02"),
					},
					Train{
						TrainID:            2,
						DepartureStationID: 2,
						ArrivalStationID:   1,
						Price:              2,
						ArrivalTime:        parseTime("00:00:02"),
						DepartureTime:      parseTime("09:00:00"),
					},
				},
				nil,
			},
		},
		{
			[]string{"2", "1", "price"},
			WantFindTrains{
				Trains{
					Train{
						TrainID:            2,
						DepartureStationID: 2,
						ArrivalStationID:   1,
						Price:              2,
						ArrivalTime:        parseTime("00:00:02"),
						DepartureTime:      parseTime("09:00:00"),
					},
					Train{
						TrainID:            9,
						DepartureStationID: 2,
						ArrivalStationID:   1,
						Price:              9,
						ArrivalTime:        parseTime("00:00:09"),
						DepartureTime:      parseTime("00:00:02"),
					},
				},
				nil,
			},
		},
	}
	dataPath = path.Join("testdata", "data.json")

	for _, tt := range tests {
		args := tt.input
		want := tt.want
		gotTrains, gotErr := FindTrains(args[0], args[1], args[2])
		if !reflect.DeepEqual(gotTrains, want.Trains) {
			t.Fatalf("got: %v, want: %v", gotTrains, want.Trains)
		}
		if gotErr != tt.want.error {
			t.Fatalf("got: %v, want: %v", gotErr, want.error)
		}
	}
}

func parseTime(value string) time.Time {
	parsed, _ := time.Parse(ParseTime, value)
	return parsed
}
