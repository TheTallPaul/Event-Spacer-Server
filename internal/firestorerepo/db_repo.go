package firestorerepo

import (
	"context"
	"errors"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

// client is a Firestore client, reused between function invocations.
var client *firestore.Client

func init() {
	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()
	config := &firebase.Config{ProjectID: projectID}

	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}
}

// FetchEvent returns a Firestore event document matching the provided ID
func FetchEvent(ctx context.Context, eventID string) (Event, error) {
	var event Event
	doc, err := client.Doc("event/" + eventID).Get(ctx)
	if err != nil {
		return event, errors.New("Fetching event " + eventID + " failed")
	}
	if err := doc.DataTo(&event); err != nil {
		return event, errors.New("Failed to convert event json data to struct")
	}

	return event, nil
}
