package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Trains []Train

type Train struct {
	TrainID            int
	DepartureStationID int
	ArrivalStationID   int
	Price              float32
	ArrivalTime        time.Time
	DepartureTime      time.Time
}

func (t Trains) findBy(depID, arrID int) Trains {
	var found Trains
	for _, train := range t {
		depOk := train.DepartureStationID == depID
		arrOk := train.ArrivalStationID == arrID
		if depOk && arrOk {
			found = append(found, train)
		}
	}
	return found
}

var sortCriterias = []string{
	0: "price",
	1: "arrival-time",
	2: "departure-time",
}

func (t Trains) sortBy(criteria string) Trains {
	sorted := make(Trains, len(t))
	copy(sorted, t)
	switch criteria {
	case sortCriterias[0]:
		sort.SliceStable(sorted, func(i, j int) bool {
			return sorted[i].Price < sorted[j].Price
		})
	case sortCriterias[1]:
		sort.SliceStable(sorted, func(i, j int) bool {
			return sorted[i].ArrivalTime.UnixNano() < sorted[j].ArrivalTime.UnixNano()
		})
	case sortCriterias[2]:
		sort.SliceStable(sorted, func(i, j int) bool {
			return sorted[i].DepartureTime.UnixNano() < sorted[j].DepartureTime.UnixNano()
		})
	default:
	}
	return sorted
}

func main() {
	result, err := FindTrains(getQuery())
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range result {
		fmt.Printf("%+v\n", v)
	}
}

func getQuery() (dep, arr, by string) {
	fmt.Print("Departure station number: ")
	fmt.Scanf("%s", &dep)
	fmt.Print("Arrival station number: ")
	fmt.Scanf("%s", &arr)
	var validBy = strings.Join(sortCriterias, "|")
	fmt.Print("Criteria (" + validBy + "): ")
	fmt.Scanf("%s", &by)
	return
}

func FindTrains(dep, arr, criteria string) (Trains, error) {
	depID, arrID, err := validateQuery(dep, arr, criteria)
	if err != nil {
		return nil, err
	}

	trains, err := getTrains()
	if err != nil {
		return nil, err
	}

	foundTrains := trains.findBy(depID, arrID).sortBy(criteria)
	if len(foundTrains) == 0 {
		return nil, nil
	}
	if len(foundTrains) > 3 {
		foundTrains = foundTrains[:3]
	}
	return foundTrains, nil
}

var ErrEmptyDep = errors.New("empty departure station")
var ErrEmptyArr = errors.New("empty arrival station")

var ErrBadDep = errors.New("bad departure station input")
var ErrBadArr = errors.New("bad arrival station input")

var ErrBadBy = errors.New("unsupported criteria")

func validateQuery(dep, arr, by string) (depID, arrID int, err error) {
	if dep == "" {
		return 0, 0, ErrEmptyDep
	}
	depID, err = strconv.Atoi(dep)
	if err != nil || depID < 0 {
		return 0, 0, ErrBadDep
	}

	if arr == "" {
		return 0, 0, ErrEmptyArr
	}
	arrID, err = strconv.Atoi(arr)
	if err != nil || arrID < 0 {
		return 0, 0, ErrBadArr
	}

	valid := false
	for _, c := range sortCriterias {
		if by == c {
			valid = true
		}
	}
	if !valid {
		return 0, 0, ErrBadBy
	}
	return depID, arrID, nil
}

const ParseTime = "15:04:05"

var dataPath = "data.json"

func getTrains() (Trains, error) {
	marshaled, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, err
	}
	var unmarshaled []map[string]any
	if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
		return nil, err
	}
	var trains Trains
	for _, u := range unmarshaled {
		var t Train

		trainID := u["trainId"].(float64)
		t.TrainID = int(trainID)

		departureStationID := u["departureStationId"].(float64)
		t.DepartureStationID = int(departureStationID)

		arrivalStationID := u["arrivalStationId"].(float64)
		t.ArrivalStationID = int(arrivalStationID)

		price := u["price"].(float64)
		t.Price = float32(price)

		if t.ArrivalTime, err = time.Parse(ParseTime, u["arrivalTime"].(string)); err != nil {
			return nil, err
		}

		if t.DepartureTime, err = time.Parse(ParseTime, u["departureTime"].(string)); err != nil {
			return nil, err
		}
		trains = append(trains, t)
	}
	return trains, nil
}
