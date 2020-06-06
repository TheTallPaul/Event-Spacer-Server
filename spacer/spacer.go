// Package spacer provides functions for
package spacer

import (
	//"log"
	"math"

	"eventspacer.org/spacer/internal/firestorerepo"
	"google.golang.org/genproto/googleapis/type/latlng"
)

// Curvature constants
const EarthRadius = 6378137.0 // meters
const DegreesToRadians = math.Pi / 180
const RadiansToDegrees = 180 / math.Pi

// Bearing directions
const N = 0.0
const E = 90.0
const SE = 135.0
const S = 180.0
const W = 270.0

// The maximum number of points allowed
const MaxNumPoints = 50000

// CreateSpacedPoints adds points spaced out the provided distance, starting in a
// clockwise spiral pattern from the center
func CreateSpacedPoints(event *firestorerepo.Event) {
	// Find and add centerpoint
	center := asTheCrowFlies(
		event.NWBoundary,
		haversineDistance(event.NWBoundary, event.SEBoundary)/2,
		SE,
	)
	event.SpacedPoints = append(event.SpacedPoints, center)

	// Prepare variables for spiral iteration
	current := center
	direction := 'N'
	spiralSegments := 1
	spiralTurn := false
	numPoints := 0
	var maxPoint = latlng.LatLng{
		Longitude: center.Longitude,
		Latitude:  center.Latitude,
	}
	var minPoint = latlng.LatLng{
		Longitude: center.Longitude,
		Latitude:  center.Latitude,
	}

	// Plop down points in a 90 degree spiral pattern, increasing the length of a
	// spiral side every two straight lines
	for (maxPoint.Latitude <= event.NWBoundary.Latitude ||
		minPoint.Latitude >= event.SEBoundary.Latitude ||
		maxPoint.Longitude <= event.SEBoundary.Longitude ||
		minPoint.Longitude >= event.NWBoundary.Longitude) &&
		numPoints <= MaxNumPoints {

		// Move the number of segments in the chosen direction
		for i := 0; i < spiralSegments; i++ {
			switch direction {
			// North
			case 'N':
				current = asTheCrowFlies(current, event.SpacingMeters, N)
				maxPoint.Latitude = current.Latitude
			// East
			case 'E':
				current = asTheCrowFlies(current, event.SpacingMeters, E)
				maxPoint.Longitude = current.Longitude
			// South
			case 'S':
				current = asTheCrowFlies(current, event.SpacingMeters, S)
				minPoint.Latitude = current.Latitude
			// West
			case 'W':
				current = asTheCrowFlies(current, event.SpacingMeters, W)
				minPoint.Longitude = current.Longitude
			}

			// Add point if in boundaries
			if inBoundaries(current, event.NWBoundary, event.SEBoundary) {
				event.SpacedPoints = append(event.SpacedPoints, current)
				numPoints = numPoints + 1
			}
		}

		// Turn
		switch direction {
		case 'N':
			direction = 'E'
		case 'E':
			direction = 'S'
		case 'S':
			direction = 'W'
		case 'W':
			direction = 'N'
		}

		// Check if there's been enough turns to increase a spiral length
		if spiralTurn {
			spiralSegments = spiralSegments + 1
			spiralTurn = false
		} else {
			spiralTurn = true
		}
	}
}

// inBoundaries finds if a point is within the boundaries of two other points, inclusive.
// For example, (1,2) would be in the boundary created by (0,0) and (2,2).
func inBoundaries(point, boundaryA, boundaryB *latlng.LatLng) bool {
	maxLat := math.Max(boundaryA.Latitude, boundaryB.Latitude)
	minLat := math.Min(boundaryA.Latitude, boundaryB.Latitude)
	maxLng := math.Max(boundaryA.Longitude, boundaryB.Longitude)
	minLng := math.Min(boundaryA.Longitude, boundaryB.Longitude)

	if point.Latitude <= maxLat && point.Latitude >= minLat &&
		point.Longitude <= maxLng && point.Longitude >= minLng {
		return true
	}

	return false
}

// asTheCrowFlies finds a new coordinate the provided distance (in meters) away in the
// bearing provided. North is bearing 0, E is bearing 90, etc. Assumes crows take into
// account the curvature of the Earth.
func asTheCrowFlies(point *latlng.LatLng, distance, bearing float64) *latlng.LatLng {
	lat := point.Latitude * DegreesToRadians
	lng := point.Longitude * DegreesToRadians
	angularDistance := distance / EarthRadius
	trueCourse := bearing * DegreesToRadians

	newLat := math.Asin(
		math.Sin(lat)*math.Cos(angularDistance) +
			math.Cos(lat)*math.Sin(angularDistance)*math.Cos(trueCourse))

	lngDist := math.Atan2(
		math.Sin(trueCourse)*math.Sin(angularDistance)*math.Cos(lat),
		math.Cos(angularDistance)-math.Sin(lat)*math.Sin(newLat))

	newLng := math.Mod(lng+lngDist+math.Pi, 2*math.Pi) - math.Pi

	newPoint := latlng.LatLng{
		Latitude:  newLat * RadiansToDegrees,
		Longitude: newLng * RadiansToDegrees,
	}

	return &newPoint
}

// haversineDistance finds the Haversine distance in meters between two points
func haversineDistance(pointA, pointB *latlng.LatLng) float64 {
	latDist := (pointB.Latitude - pointA.Latitude) * DegreesToRadians
	lngDist := (pointB.Longitude - pointA.Longitude) * DegreesToRadians

	latA := pointA.Latitude * DegreesToRadians
	latB := pointB.Latitude * DegreesToRadians

	return 2 * EarthRadius * math.Asin(math.Sqrt(
		math.Pow(math.Sin(latDist/2), 2)+
			math.Pow(math.Sin(lngDist/2), 2)*math.Cos(latA)*math.Cos(latB)))
}
