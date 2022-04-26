package legs

import (
	"context"
	"github.com/ipld/go-ipld-prime"
	"github.com/kenlabs/pandofg/pkg/types/schema"
	"go.mongodb.org/mongo-driver/mongo"
)

func CommitPayloadToMetastore(data ipld.Node, client *mongo.Client) error {
	//var payload map[string]interface{}
	//if err := json.Unmarshal(data, &payload); err != nil {
	//	return err
	//}
	locations, err := schema.UnwrapLocation(data)
	locationCollection := client.Database("pando-fg").Collection("locations")
	result, err := locationCollection.InsertOne(context.TODO(), locations)
	if err != nil {
		return err
	}
	log.Debugf("insert a doc into mongo, ID: %s", result.InsertedID)
	return nil
}
