package schema

import (
	_ "embed"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/schema"
)

var (
	MinerLocationPrototype schema.TypedPrototype

	//go:embed green.ipldsch
	greenSchemaBytes []byte
)

func init() {
	typeSystem, err := ipld.LoadSchemaBytes(greenSchemaBytes)
	if err != nil {
		panic(fmt.Errorf("failed to load schema: %w", err))
	}
	MinerLocationPrototype = bindnode.Prototype((*MinerLocationsModel)(nil), typeSystem.TypeByName("MinerLocationsModel"))
}

func FilecoinGreenMinerLocationsType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Locations",
		Fields: graphql.Fields{
			"epoch": &graphql.Field{
				Type: graphql.Int,
			},
			"date": &graphql.Field{
				Type: graphql.String,
			},
			"minerLocations": &graphql.Field{
				Type: graphql.NewList(FilecoinGreenMinerLocationType()),
			},
		},
	})
}

func FilecoinGreenMinerLocationType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Location",
		Fields: graphql.Fields{
			"miner": &graphql.Field{
				Type: graphql.String,
			},
			"region": &graphql.Field{
				Type: graphql.String,
			},
			"lat": &graphql.Field{
				Type: graphql.Float,
			},
			"long": &graphql.Field{
				Type: graphql.Float,
			},
			"numLocations": &graphql.Field{
				Type: graphql.Int,
			},
			"subdiv1": &graphql.Field{
				Type: graphql.String,
			},
			"country": &graphql.Field{
				Type: graphql.String,
			},
			"city": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
}

type MinerLocationsModel struct {
	//ID             primitive.ObjectID    `bson:"_id,omitempty" json:"_id,omitempty"`
	Epoch          uint64                `bson:"epoch,omitempty" json:"epoch,omitempty"`
	Date           string                `bson:"date,omitempty" json:"date,omitempty"`
	MinerLocations []*MinerLocationModel `bson:"minerLocations,omitempty" json:"minerLocations,omitempty"`
}

type MinerLocationModel struct {
	Miner        string  `bson:"miner,omitempty" json:"miner,omitempty"`
	Region       string  `bson:"region,omitempty" json:"region,omitempty"`
	Long         float32 `bson:"long,omitempty" json:"long,omitempty"`
	Lat          float32 `bson:"lat,omitempty" json:"lat,omitempty"`
	NumLocations int     `bson:"numLocations,omitempty" json:"numLocations,omitempty"`
	Country      string  `bson:"country,omitempty" json:"country,omitempty"`
	City         string  `bson:"city,omitempty" json:"city,omitempty"`
	SubDiv1      string  `bson:"subdiv1,omitempty" json:"subdiv1,omitempty"`
}

func UnwrapLocation(node ipld.Node) (*MinerLocationsModel, error) {
	if node.Prototype() != MinerLocationPrototype {
		adBuilder := MinerLocationPrototype.NewBuilder()
		err := adBuilder.AssignNode(node)
		if err != nil {
			return nil, fmt.Errorf("faild to convert node prototype: %v", err)
		}
		node = adBuilder.Build()
	}

	ad, ok := bindnode.Unwrap(node).(*MinerLocationsModel)
	if !ok || ad == nil {
		return nil, fmt.Errorf("unwrapped node does not match MinerLocationsModel")
	}
	return ad, nil
}