package firestorerepo

import (
	"time"

	"google.golang.org/genproto/googleapis/type/latlng"
)

// An Event is a struct conversion of a Firestore document for events
type Event struct {
	ID            string           `firestore:"-"`
	ClaimedSpots  map[string]bool  `firestore:"claimed_spots"`
	Name          string           `firestore:"name"`
	NWBoundary    *latlng.LatLng   `firestore:"nw_boundary"`
	SEBoundary    *latlng.LatLng   `firestore:"se_boundary"`
	SpacedPoints  []*latlng.LatLng `firestore:"spaced_points"`
	SpacingMeters float64          `firestore:"spacing_meters"`
}

// FirestoreEvent is the payload of a Firestore event.
type FirestoreEvent struct {
	OldValue   FirestoreValue `json:"oldValue"`
	Value      FirestoreValue `json:"value"`
	UpdateMask struct {
		FieldPaths []string `json:"fieldPaths"`
	} `json:"updateMask"`
}

// FirestoreValue holds Firestore fields.
type FirestoreValue struct {
	CreateTime time.Time `json:"createTime"`
	Fields     RawEvent  `json:"fields"`
	Name       string    `json:"name"`
	UpdateTime time.Time `json:"updateTime"`
}

// RawEvent is a struct that holds the Firestore fields for eventual conversion to
// Document Snapshot type Event. Yes, this is overcomplicated, watch
// https://github.com/googleapis/google-cloud-go/issues/1438 for an eventual solution.
type RawEvent struct {
	Name struct {
		Value string `json:"stringValue"`
	} `json:"name"`
	NWBoundary struct {
		Value *latlng.LatLng `json:"geoPointValue"`
	} `json:"nw_boundary"`
	SEBoundary struct {
		Value *latlng.LatLng `json:"geoPointValue"`
	} `json:"se_boundary"`
	SpacingMeters struct {
		Value float64 `json:"doubleValue"`
	} `json:"spacing_meters"`
}
