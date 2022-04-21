package schema

import "github.com/graphql-go/graphql"

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
