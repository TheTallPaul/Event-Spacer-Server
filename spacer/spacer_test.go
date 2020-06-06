package spacer

import (
	"reflect"
	"testing"
	"time"

	"eventspacer.org/spacer/internal/firestorerepo"
	"google.golang.org/genproto/googleapis/type/latlng"
)

var zeroPoint = latlng.LatLng{Latitude: 0.0, Longitude: 0.0}
var halfPoint = latlng.LatLng{
	Latitude:  0.4999809639485908,
	Longitude: -0.5000000014497166,
}
var onePoint = latlng.LatLng{Latitude: 1.0, Longitude: -1.0}
var caseOnePoint = latlng.LatLng{Latitude: 0.000017966305682390428, Longitude: 0.0}
var broadwayYamhill = latlng.LatLng{Latitude: 45.518673, Longitude: -122.679996}
var pioneerCourthouse = latlng.LatLng{
	Latitude:  45.51867299282918,
	Longitude: -122.67871393198214,
}
var eventPoint1 = latlng.LatLng{
	Latitude:  0.7515092435020568,
	Longitude: -0.5000000014497166,
}
var eventPoint2 = latlng.LatLng{
	Latitude:  0.7515020015288297,
	Longitude: -0.24845008446180564,
}
var eventPoint3 = latlng.LatLng{
	Latitude:  0.49997372197536377,
	Longitude: -0.24845008446180564,
}
var eventPoint4 = latlng.LatLng{
	Latitude:  0.2484454424218978,
	Longitude: -0.24845008446180564,
}
var eventPoint5 = latlng.LatLng{
	Latitude:  0.2484430483816972,
	Longitude: -0.4999807286902337,
}
var eventPoint6 = latlng.LatLng{
	Latitude:  0.24844065436456605,
	Longitude: -0.7515113728730652,
}
var eventPoint7 = latlng.LatLng{
	Latitude:  0.4999689339180321,
	Longitude: -0.7515113728730652,
}
var eventPoint8 = latlng.LatLng{
	Latitude:  0.751497213471498,
	Longitude: -0.7515113728730652,
}

var eventTest1 = firestorerepo.Event{
	ClaimedSpots:  map[string]time.Time{},
	Expiration:    time.Time{},
	Name:          "eventTest1",
	NWBoundary:    &onePoint,
	SEBoundary:    &zeroPoint,
	SpacedPoints:  []*latlng.LatLng{},
	SpacingMeters: 28000,
}

var expectedEventTest1 = firestorerepo.Event{
	Expiration: time.Time{},
	Name:       "eventTest1",
	NWBoundary: &onePoint,
	SEBoundary: &zeroPoint,
	SpacedPoints: []*latlng.LatLng{
		&halfPoint,
		&eventPoint1,
		&eventPoint2,
		&eventPoint3,
		&eventPoint4,
		&eventPoint5,
		&eventPoint6,
		&eventPoint7,
		&eventPoint8,
	},
	SpacingMeters: 28000,
}

var CreateSpacedPointsTestCases = []struct {
	event         *firestorerepo.Event
	expectedEvent *firestorerepo.Event
}{
	{
		&eventTest1,
		&expectedEventTest1,
	},
}

func TestCreateSpacedPoints(t *testing.T) {
	for _, input := range CreateSpacedPointsTestCases {
		CreateSpacedPoints(input.event)

		if !reflect.DeepEqual(
			input.event.SpacedPoints,
			input.expectedEvent.SpacedPoints) {
			t.Errorf(
				"FAIL: Want spaced event points to be %v but we got %v",
				input.expectedEvent.SpacedPoints,
				input.event.SpacedPoints,
			)
		}
	}
}

var inBoundariesTestCases = []struct {
	point         *latlng.LatLng
	boundaryA     *latlng.LatLng
	boundaryB     *latlng.LatLng
	expectedTruth bool
}{
	{
		&halfPoint,
		&zeroPoint,
		&onePoint,
		true,
	},
	{
		&pioneerCourthouse,
		&zeroPoint,
		&onePoint,
		false,
	},
}

func TestInBoundaries(t *testing.T) {
	for _, input := range inBoundariesTestCases {
		truthiness := inBoundaries(input.point, input.boundaryA, input.boundaryB)

		if truthiness != input.expectedTruth {
			t.Errorf(
				"FAIL: Want evaluation of (%v, %v) being in the "+
					"boundary of (%v, %v) and (%v, %v) to be: %v, "+
					"but we got %v",
				input.point.Latitude,
				input.point.Longitude,
				input.boundaryA.Latitude,
				input.boundaryA.Longitude,
				input.boundaryB.Latitude,
				input.boundaryB.Longitude,
				input.expectedTruth,
				truthiness,
			)
		}
	}
}

var asTheCrowFliesTestCases = []struct {
	point         *latlng.LatLng
	distance      float64
	bearing       float64
	expectedPoint *latlng.LatLng
}{
	{
		&zeroPoint,
		2,
		0,
		&caseOnePoint,
	},
	{
		&broadwayYamhill,
		100,
		90,
		&pioneerCourthouse,
	},
}

func TestAsTheCrowFlies(t *testing.T) {
	for _, input := range asTheCrowFliesTestCases {
		coord := asTheCrowFlies(
			input.point,
			input.distance,
			input.bearing,
		)

		if coord.Latitude != input.expectedPoint.Latitude ||
			coord.Longitude != input.expectedPoint.Longitude {
			t.Errorf(
				"FAIL: Want correct point from (%v, %v) %v meters away "+
					"in bearing %v° to be: %v but we got (%v, %v)",
				input.point.Latitude,
				input.point.Longitude,
				input.distance,
				input.bearing,
				input.expectedPoint.Latitude,
				input.expectedPoint.Longitude,
				coord,
			)
		}
	}
}

var haversineDistanceTestCases = []struct {
	pointA           *latlng.LatLng
	pointB           *latlng.LatLng
	expectedDistance float64
}{
	{
		&zeroPoint,
		&onePoint,
		157425.537108412,
	},
	{
		&broadwayYamhill,
		&pioneerCourthouse,
		99.99999999829213,
	},
}

func TestHaversineDistance(t *testing.T) {
	for _, input := range haversineDistanceTestCases {
		distance := haversineDistance(input.pointA, input.pointB)

		if distance != input.expectedDistance {
			t.Errorf(
				"FAIL: Want distance from (%v, %v) to (%v, %v) to be: "+
					"%v but we got %v",
				input.pointA.Latitude,
				input.pointA.Longitude,
				input.pointB.Latitude,
				input.pointB.Longitude,
				input.expectedDistance,
				distance,
			)
		}
	}
}
