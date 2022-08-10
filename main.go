package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
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
		if train.DepartureStationID == depID && train.ArrivalStationID == arrID {
			found = append(found, train)
		}
	}

	return found
}

const (
	sortCriteriaPrice = "price"
	sortCriteriaArr   = "arrival-time"
	sortCriteriaDep   = "departure-time"
)

func (t Trains) sortBy(criteria string) (Trains, error) {
	sorted := make(Trains, len(t))
	copy(sorted, t)

	switch criteria {
	case sortCriteriaPrice:
		sort.SliceStable(sorted, func(i, j int) bool {
			return sorted[i].Price < sorted[j].Price
		})
	case sortCriteriaArr:
		sort.SliceStable(sorted, func(i, j int) bool {
			return sorted[i].ArrivalTime.UnixNano() < sorted[j].ArrivalTime.UnixNano()
		})
	case sortCriteriaDep:
		sort.SliceStable(sorted, func(i, j int) bool {
			return sorted[i].DepartureTime.UnixNano() < sorted[j].DepartureTime.UnixNano()
		})
	default:
		return nil, errBadCriteria
	}

	return sorted, nil
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

func getQuery() (dep, arr, criteria string) {
	fmt.Print("Departure station number: ")
	fmt.Scanf("%s", &dep)

	fmt.Print("Arrival station number: ")
	fmt.Scanf("%s", &arr)

	fmt.Print("Criteria: ")
	fmt.Scanf("%s", &criteria)

	return
}

var (
	errEmptyDep = errors.New("empty departure station")
	errEmptyArr = errors.New("empty arrival station")

	errBadDep = errors.New("bad departure station input")
	errBadArr = errors.New("bad arrival station input")

	errBadCriteria = errors.New("unsupported criteria")
)

func FindTrains(dep, arr, criteria string) (Trains, error) {
	if dep == "" {
		return nil, errEmptyDep
	}
	depID, err := strconv.Atoi(dep)
	if err != nil || depID < 0 {
		return nil, errBadDep
	}

	if arr == "" {
		return nil, errEmptyArr
	}
	arrID, err := strconv.Atoi(arr)
	if err != nil || arrID < 0 {
		return nil, errBadArr
	}

	switch criteria {
	case sortCriteriaPrice, sortCriteriaArr, sortCriteriaDep:
	default:
		return nil, errBadCriteria
	}

	trains, err := getTrains()
	if err != nil {
		return nil, err
	}

	foundTrains, err := trains.findBy(depID, arrID).sortBy(criteria)
	if err != nil {
		return nil, err
	}

	if len(foundTrains) == 0 {
		return nil, nil
	}
	if len(foundTrains) > 3 {
		foundTrains = foundTrains[:3]
	}

	return foundTrains, nil
}

const (
	parseTime = "15:04:05"
	dataPath  = "data.json"
)

var errTypeAssertion = errors.New("invalid type assertion")

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

		trainID, ok := u["trainId"].(float64)
		if !ok {
			return nil, errTypeAssertion
		}
		t.TrainID = int(trainID)

		departureStationID, ok := u["departureStationId"].(float64)
		if !ok {
			return nil, errTypeAssertion
		}
		t.DepartureStationID = int(departureStationID)

		arrivalStationID, ok := u["arrivalStationId"].(float64)
		if !ok {
			return nil, errTypeAssertion
		}
		t.ArrivalStationID = int(arrivalStationID)

		price, ok := u["price"].(float64)
		if !ok {
			return nil, errTypeAssertion
		}
		t.Price = float32(price)

		arrTime, ok := u["arrivalTime"].(string)
		if !ok {
			return nil, errTypeAssertion
		}
		if t.ArrivalTime, err = time.Parse(parseTime, arrTime); err != nil {
			return nil, err
		}

		depTime, ok := u["departureTime"].(string)
		if !ok {
			return nil, errTypeAssertion
		}
		if t.DepartureTime, err = time.Parse(parseTime, depTime); err != nil {
			return nil, err
		}

		trains = append(trains, t)
	}

	return trains, nil
}
