// Package spacer provides functions for
package spacer

import (
	//"fmt"
	"math"

	"google.golang.org/genproto/googleapis/type/latlng"
	"eventspacer.org/spacer/internal/firestorerepo"
)

func FindSpacedLocation(event firestorerepo.Event) *latlng.LatLng {

	return event.NWBoundary

}

// asTheCrowFlies finds a new coordinate the provided distance (in meters) away
/// in the bearing provided. North is bearing 0, E is bearing 90, etc.
func asTheCrowFlies(point *latlng.LatLng, distance, bearing float64) *latlng.LatLng {
	earthRadius := 6378137.0
	degreesToRadians := math.Pi / 180
	radiansToDegrees := 180 / math.Pi

	latA := point.Latitude * degreesToRadians
	lngA := point.Longitude * degreesToRadians
	angularDistance := distance / earthRadius
	trueCourse := bearing * degreesToRadians

	latB := math.Asin(
		math.Sin(latA) * math.Cos(angularDistance) +
		math.Cos(latA) * math.Sin(angularDistance) * math.Cos(trueCourse))

	lngDist := math.Atan2(
		math.Sin(trueCourse) * math.Sin(angularDistance) * math.Cos(latA),
		math.Cos(angularDistance) - math.Sin(latA) * math.Sin(latB))

	lngB := math.Mod(lngA + lngDist + math.Pi, 2 * math.Pi) - math.Pi

	newPoint := latlng.LatLng{
		Latitude: latB * radiansToDegrees,
		Longitude: lngB * radiansToDegrees,
	}

	return &newPoint
}
