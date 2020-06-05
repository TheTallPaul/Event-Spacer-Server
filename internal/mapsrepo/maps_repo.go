package mapsrepo

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/genproto/googleapis/type/latlng"
	"googlemaps.github.io/maps"
)

// client is a Google Maps client, reused between function invocations.
var client *maps.Client

func init() {
	var err error
	client, err = maps.NewClient(maps.WithAPIKey(mapsAPIKey))
	if err != nil {
		log.Fatalf("maps.NewClient: %s", err)
	}
}

// mapsToFireLatLng converts Google Map's LatLng to the Google API LatLng. Why
// the same company has two different datatypes for the same information I'll
// never know...
func mapsToFireLatLng(original maps.LatLng) *latlng.LatLng {
	return *latlng.LatLng{
		Latitude: original.Lat,
		Longitude: original.Lng
	}
}

// fireToMapsLatLng converts Google API LatLng to Google Map's LatLng
func fireToMapsLatLng(original *latlng.LatLng) maps.LatLng {
	return maps.LatLng{
		Lat: original.Latitude,
		Lng: original.Longitude
	}
}

// NearestRoad makes Google Maps API to return a slice of the nearest
// coordinates of a road for the provided coordinates
func NearestRoad(points []*latlng.LatLng) []*latlng.LatLng {
	// Convert points into agreeable format
	var convertedPoints := []maps.LatLng{}
	for _, coord := range points {
		convertedPoints = append(
			convertedPoints, fireToMapsLatLng(coord))
	}

	// Make request and handle errors
	request := &maps.NearestRoadRequest{Points: convertedPoints}
	snappedPoints, _, err := client.NearestRoad(
		context.Background(), request)
	if err != nil {
		log.Fatalf("Nearest Road fatal error: %s", err)
	}

	// Convert response to Firestore format
	var roadPoints := []*latlng.LatLng{}
	for _, snapPoint := range snappedPoints {
		roadPoints = append(
			roadPoints, mapsToFireLatLng(snapPoint.Location))
	}

	return roadPoints
}
