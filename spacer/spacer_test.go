package spacer

import (
	//"reflect"
	"testing"

	"google.golang.org/genproto/googleapis/type/latlng"
	//"googlemaps.github.io/maps"

	//"eventspacer.org/gridunlockridematch/internal/firebaserepo"
)

var zeroPoint = latlng.LatLng{Latitude: 0.0, Longitude: 0.0}
var caseOnePoint = latlng.LatLng{Latitude: 0.000017966305682390428, Longitude: 0.0}
var broadwayYamhill = latlng.LatLng{Latitude: 45.518673, Longitude: -122.679996}
var pioneerCourthouse = latlng.LatLng{Latitude: 45.51867299282918,
	Longitude: -122.67871393198214}

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
				"FAIL: Want correct point from (%v, %v) %v meters away " +
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
