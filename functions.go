// Module eventspacer provides packages that can modify Firebase documents with gps points
// evenly spaced out, taking into account the curvature of the Earth
package eventspacer

import (
	"context"
	"log"
	"path"
	"time"

	"google.golang.org/genproto/googleapis/type/latlng"

	"eventspacer.org/spacer/internal/firestorerepo"
	"eventspacer.org/spacer/spacer"
)

// The maximum number of geopoints allowed
const maxNumPoints = 10000

// SpaceEvent responds to the firestore trigger of an event document being added to the
// database and fills the provided areas with gps points spaced out the provided distance.
// It then calls for this updated Event to be pushed back to the database.
func SpaceEvent(ctx context.Context, fireEvent firestorerepo.FirestoreEvent) error {
	log.Println("FirestoreEvent: ", fireEvent)

	// Create a correctly-formatted event from the payload
	var event = firestorerepo.Event{
		ID:            path.Base(fireEvent.Value.Name),
		ClaimedSpots:  map[string]time.Time{},
		Name:          fireEvent.Value.Fields.Name.Value,
		NWBoundary:    fireEvent.Value.Fields.NWBoundary.Value,
		SEBoundary:    fireEvent.Value.Fields.SEBoundary.Value,
		SpacedPoints:  []*latlng.LatLng{},
		SpacingMeters: fireEvent.Value.Fields.SpacingMeters.Value,
	}
	log.Println("firestorerepo.Event created: ", event)

	// Add the geo-points to the event
	spacer.CreateSpacedPoints(&event, maxNumPoints)
	log.Println("firestorerepo.Event filled with spaced points: ", event)

	// Update the doc
	err := firestorerepo.UpdateEvent(ctx, event)
	if err != nil {
		return err
	}

	return nil
}
