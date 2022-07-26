package travel

import (
	"fmt"
	"math"
	"time"

	"stop-checker.com/db"
	"stop-checker.com/db/model"
)

type PlannedLeg struct {
	Origin      string // origin stop id
	Destination string // destination stop id
	Walk        bool   // if we walk between the two stops
	TripId      string // trip id

	Departure time.Time     // when do we depart from the origin
	Duration  time.Duration // duration between arriving at the destination and leaving the origin
}

func (pl *PlannedLeg) String() string {
	if pl.Walk {
		return fmt.Sprintf("{origin:%s, destination:%s, walk:%t, departure:%s, duration:%s}",
			pl.Origin, pl.Destination, pl.Walk, pl.Departure, pl.Duration)
	}
	return fmt.Sprintf("{origin:%s, destination:%s, walk:%t, departure:%s, duration:%s, trip:%s}",
		pl.Origin, pl.Destination, pl.Walk, pl.Departure, pl.Duration, pl.TripId)
}

/* FixedLeg
- leg of a travel plan without assigned times
*/
type FixedLeg struct {
	Origin      string // stop id
	Destination string // stop id
	RouteId     string // stop id
	Walk        bool   // if true then: route id is empty
}

func (fl *FixedLeg) String() string {
	if fl.Walk {
		return fmt.Sprintf("walk{origin:%s, destination:%s}", fl.Origin, fl.Destination)
	}
	return fmt.Sprintf("transit{origin:%s, destination:%s}", fl.Origin, fl.Destination)
}

const ARRIVE_BY = 0
const DEPART_AT = false

type FixedPlanner struct {
	StopIndex         *db.Index[model.Stop]
	StopTimesFromTrip *db.InvertedIndex[model.StopTime]
	ScheduleIndex     *db.ScheduleIndex
}

func (p *FixedPlanner) Depart(at time.Time, fixed []*FixedLeg) ([]*PlannedLeg, error) {
	planned := []*PlannedLeg{}
	acc := at

	for _, leg := range fixed {
		plan, err := p.planDepart(acc, leg)
		if err != nil {
			return nil, err
		}
		acc = plan.Departure.Add(plan.Duration)
		planned = append(planned, plan)
	}

	return planned, nil
}

func (p *FixedPlanner) Arrive(by time.Time, fixed []*FixedLeg) ([]*PlannedLeg, error) {
	planned := []*PlannedLeg{}
	acc := by

	for i := len(fixed) - 1; i >= 0; i-- {
		leg := fixed[i]
		plan, err := p.planArrive(acc, leg)
		if err != nil {
			return nil, err
		}

		acc = plan.Departure
		planned = append(planned, plan)
	}

	order := []*PlannedLeg{}
	for i := 0; i < len(planned); i++ {
		order = append(order, planned[len(planned)-(i+1)])
	}

	return order, nil
}

func (p *FixedPlanner) planDepart(acc time.Time, fixed *FixedLeg) (*PlannedLeg, error) {
	// get origin and destination stops
	origin, destination, err := p.stops(fixed)
	if err != nil {
		return nil, err
	}

	if fixed.Walk {
		// planned leg by walking
		distance := origin.Distance(destination.Location)
		duration := time.Duration(math.Round(distance*1.4/60)) * time.Minute

		return &PlannedLeg{
			Origin:      fixed.Origin,
			Destination: fixed.Destination,
			Walk:        true,
			TripId:      "",
			Departure:   acc,
			Duration:    duration,
		}, nil
	}

	// planned leg by transit
	next, err := p.ScheduleIndex.Get(fixed.Origin, fixed.RouteId).Next(acc)
	if err != nil {
		return nil, err
	}

	// all stop times next trip
	all, _ := p.StopTimesFromTrip.Get(next.TripId)

	// origin stop times
	originArrival, err := p.stopTime(fixed.Origin, all)
	if err != nil {
		return nil, err
	}

	// destination stop times
	destinationArrival, err := p.stopTime(fixed.Destination, all)
	if err != nil {
		return nil, err
	}

	wait := stopTimeDiffDuration(acc, originArrival.Time)
	duration := stopTimeDiffDuration(originArrival.Time, destinationArrival.Time)

	// planned leg
	return &PlannedLeg{
		Origin:      fixed.Origin,
		Destination: fixed.Destination,
		Walk:        false,
		TripId:      next.TripId,
		Departure:   acc.Add(wait),
		Duration:    duration,
	}, nil
}

func (p *FixedPlanner) planArrive(acc time.Time, fixed *FixedLeg) (*PlannedLeg, error) {
	// get origin and destination stops
	origin, destination, err := p.stops(fixed)
	if err != nil {
		return nil, err
	}

	if fixed.Walk {
		// planned leg by walking
		distance := origin.Distance(destination.Location)
		duration := time.Duration(math.Round(distance*1.4/60)) * time.Minute

		return &PlannedLeg{
			Origin:      fixed.Origin,
			Destination: fixed.Destination,
			Walk:        true,
			TripId:      "",
			Departure:   acc.Add(-duration),
			Duration:    duration,
		}, nil
	}

	// planned leg by transit
	previous, err := p.ScheduleIndex.Get(fixed.Origin, fixed.RouteId).Previous(acc)
	if err != nil {
		return nil, err
	}

	// all stop times next trip
	all, _ := p.StopTimesFromTrip.Get(previous.TripId)

	// origin stop times
	originArrival, err := p.stopTime(fixed.Origin, all)
	if err != nil {
		return nil, err
	}

	// destination stop times
	destinationArrival, err := p.stopTime(fixed.Destination, all)
	if err != nil {
		return nil, err
	}

	diff := stopTimeDiffDuration(originArrival.Time, acc)
	departure := acc.Add(-diff)
	duration := stopTimeDiffDuration(originArrival.Time, destinationArrival.Time)

	// planned leg
	return &PlannedLeg{
		Origin:      fixed.Origin,
		Destination: fixed.Destination,
		Walk:        false,
		TripId:      previous.TripId,
		Departure:   departure,
		Duration:    duration,
	}, nil
}

func (p *FixedPlanner) stops(fixed *FixedLeg) (model.Stop, model.Stop, error) {
	empty := model.Stop{}
	origin, ok := p.StopIndex.Get(fixed.Origin)
	if !ok {
		return empty, empty, fmt.Errorf("origin stop %s not found", fixed.Origin)
	}

	destination, ok := p.StopIndex.Get(fixed.Destination)
	if !ok {
		return empty, empty, fmt.Errorf("destination stop %s not found", fixed.Destination)
	}
	return origin, destination, nil
}

func (p *FixedPlanner) stopTime(stopId string, all []model.StopTime) (model.StopTime, error) {
	for _, stopTime := range all {
		if stopTime.StopId == stopId {
			return stopTime, nil
		}
	}

	return model.StopTime{}, fmt.Errorf("stoptime not found stop:%s", stopId)
}
