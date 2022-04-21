package legs

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
)

func CommitPayloadToMetastore(data []byte, client *mongo.Client) error {
	var payload map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	locationCollection := client.Database("pando-fg").Collection("locations")
	result, err := locationCollection.InsertOne(context.TODO(), payload)
	if err != nil {
		return err
	}
	log.Debugf("insert a doc into mongo, ID: %s", result.InsertedID)
	return nil
}
