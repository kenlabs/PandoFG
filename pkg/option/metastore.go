package option

import "go.mongodb.org/mongo-driver/mongo"

const (
	defaultMetaStoreType          = "mongodb"
	defaultMetaStoreConnectionURI = "mongodb://52.14.211.248:27018"
)

type MetaStore struct {
	Type          string `yaml:"Type"`
	ConnectionURI string `yaml:"ConnectionURI"`
	Client        *mongo.Client
}
