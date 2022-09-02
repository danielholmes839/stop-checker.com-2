// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package sdl

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"stop-checker.com/db/model"
	"stop-checker.com/features/travel"
)

type PageInfo struct {
	Cursor    int `json:"cursor"`
	Remaining int `json:"remaining"`
}

type PageInput struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
}

type StopSearchPayload struct {
	Page    *PageInfo     `json:"page"`
	Results []*model.Stop `json:"results"`
}

type TravelLegInput struct {
	Origin      string  `json:"origin"`
	Destination string  `json:"destination"`
	Route       *string `json:"route"`
}

type TravelOptions struct {
	Datetime *time.Time   `json:"datetime"`
	Mode     ScheduleMode `json:"mode"`
}

type TravelRoutePayload struct {
	Route []*travel.FixedLeg `json:"route"`
	Error *string            `json:"error"`
}

type TravelSchedulePayload struct {
	Schedule travel.Schedule `json:"schedule"`
	Error    *string         `json:"error"`
}

type RouteType string

const (
	RouteTypeBus   RouteType = "BUS"
	RouteTypeTrain RouteType = "TRAIN"
)

var AllRouteType = []RouteType{
	RouteTypeBus,
	RouteTypeTrain,
}

func (e RouteType) IsValid() bool {
	switch e {
	case RouteTypeBus, RouteTypeTrain:
		return true
	}
	return false
}

func (e RouteType) String() string {
	return string(e)
}

func (e *RouteType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RouteType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RouteType", str)
	}
	return nil
}

func (e RouteType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ScheduleMode string

const (
	ScheduleModeArriveBy ScheduleMode = "ARRIVE_BY"
	ScheduleModeDepartAt ScheduleMode = "DEPART_AT"
)

var AllScheduleMode = []ScheduleMode{
	ScheduleModeArriveBy,
	ScheduleModeDepartAt,
}

func (e ScheduleMode) IsValid() bool {
	switch e {
	case ScheduleModeArriveBy, ScheduleModeDepartAt:
		return true
	}
	return false
}

func (e ScheduleMode) String() string {
	return string(e)
}

func (e *ScheduleMode) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ScheduleMode(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ScheduleMode", str)
	}
	return nil
}

func (e ScheduleMode) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
