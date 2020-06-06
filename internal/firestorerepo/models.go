package firestorerepo

import (
	"time"

	"google.golang.org/genproto/googleapis/type/latlng"
)

// An Event is a struct conversion of a Firestore document for events
type Event struct {
	ID            string               `firestore:"-"`
	ClaimedSpots  map[string]time.Time `firestore:"claimed_spots"`
	Expiration    time.Time            `firestore:"expiration"`
	Name          string               `firestore:"name"`
	NWBoundary    *latlng.LatLng       `firestore:"nw_boundary"`
	SEBoundary    *latlng.LatLng       `firestore:"se_boundary"`
	SpacedPoints  []*latlng.LatLng     `firestore:"spaced_points"`
	SpacingMeters float64              `firestore:"spacing_meters"`
}
