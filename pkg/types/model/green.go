package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type MinerLocationsModel struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Epoch          int64                `bson:"epoch,omitempty" json:"epoch,omitempty"`
	Date           string               `bson:"date,omitempty" json:"date,omitempty"`
	MinerLocations []MinerLocationModel `bson:"minerLocations,omitempty" json:"minerLocations,omitempty"`
}

type MinerLocationModel struct {
	Miner        string  `bson:"miner,omitempty" json:"miner,omitempty"`
	Region       string  `bson:"region,omitempty" json:"region,omitempty"`
	Long         float64 `bson:"long,omitempty" json:"long,omitempty"`
	Lat          float64 `bson:"lat,omitempty" json:"lat,omitempty"`
	NumLocations int     `bson:"numLocations,omitempty" json:"numLocations,omitempty"`
	Country      string  `bson:"country,omitempty" json:"country,omitempty"`
	City         string  `bson:"city,omitempty" json:"city,omitempty"`
	Subdiv1      string  `bson:"subdiv1,omitempty" json:"subdiv1,omitempty"`
}
