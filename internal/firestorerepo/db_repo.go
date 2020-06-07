package firestorerepo

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

// client is a Firestore client, reused between function invocations
var client *firestore.Client

func init() {
	// Use context.Background() because the app/client should persist across
	// invocations
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

// UpdateEvent updates the Firebase collection with the provided event
func UpdateEvent(ctx context.Context, event Event) error {
	result, err := client.Doc("events/"+event.ID).Set(ctx, event)
	log.Println("client.Doc: ", result)

	return err
}
