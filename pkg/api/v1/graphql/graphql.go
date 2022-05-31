package graphql

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/ipfs/go-cid"
	"github.com/kenlabs/pandofg/pkg/types/schema"
	"github.com/libp2p/go-libp2p-core/peer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"html/template"
	"io"
	"net/http"
	"time"
)

func (a *API) registerGraphql() {
	tmpl, err := template.New("index.html").Parse(indexHTML())
	if err != nil {
		panic(err)
	}
	err = a.NewSchema()
	if err != nil {
		panic(err)
	}
	a.router.SetHTMLTemplate(tmpl)
	a.router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	a.router.POST("/search", a.search)
}

func (a *API) search(ctx *gin.Context) {
	var result *graphql.Result

	if ctx.Request.Header.Get("Content-Type") == "application/json" {
		var p postData
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(ctx.Request.Body)

		if err := json.NewDecoder(ctx.Request.Body).Decode(&p); err != nil {
			http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
			return
		}

		result = graphql.Do(graphql.Params{
			Context:        context.Background(),
			Schema:         a.schema,
			RequestString:  p.Query,
			VariableValues: p.Variables,
			OperationName:  p.Operation,
		})
	} else {
		err := ctx.Request.ParseForm()
		if err != nil {
			logger.Warnf("failed to read req: %v", err)
			return
		}
		result = graphql.Do(graphql.Params{
			Context:       context.Background(),
			Schema:        a.schema,
			RequestString: ctx.Request.Form.Get("query"),
		})
	}

	if len(result.Errors) > 0 {
		logger.Infof("Query had errors: %s, %v",
			ctx.Request.URL.Query().Get("query"), result.Errors)
	}
	if err := json.NewEncoder(ctx.Writer).Encode(result); err != nil {
		logger.Errorf("Failed to encode response: %s", err)
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *API) NewSchema() error {
	var err error
	a.schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"State":    a.newStateField(),
				"SnapShot": a.newSnapshotField(),
				"GetProvidersLocations": &graphql.Field{
					Type: schema.FilecoinGreenMinerLocationsType(),
					Args: graphql.FieldConfigArgument{
						"minerID": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						var minerID string
						minerIDArg, exists := p.Args["minerID"]
						if exists {
							minerID = minerIDArg.(string)
						}

						if minerID != "" {
							q := []interface{}{
								bson.M{"$match": bson.M{"minerLocations.miner": minerID}},
								bson.M{"$sort": bson.M{"date": -1}},
								bson.M{"$addFields": bson.M{
									"minerLocations": bson.M{
										"$filter": bson.M{
											"input": "$minerLocations",
											"as":    "minerLocation",
											"cond": bson.M{
												"$in": []interface{}{
													"$$minerLocation.miner",
													[]string{minerID},
												},
											},
										},
									},
								}},
							}
							ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)
							opts := options.Aggregate().SetAllowDiskUse(true)
							result, err := a.core.StoreInstance.MetaStore.
								Database("pando-fg").
								Collection("locations").
								Aggregate(ctx, q, opts)
							if err != nil {
								logger.Errorf("decode location results failed: %v", err)
								return nil, err
							}
							var results []schema.MinerLocationsModel
							err = result.All(ctx, &results)
							if err != nil {
								return nil, err
							}
							return results[0], nil
						} else {
							sortOpt := options.FindOne().SetSort(bson.M{"date": -1})
							ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)
							var result schema.MinerLocationsModel
							err = a.core.StoreInstance.MetaStore.
								Database("pando-fg").
								Collection("locations").
								FindOne(ctx, bson.M{}, sortOpt).
								Decode(&result)
							if err != nil {
								return nil, err
							}
							return result, err
						}
					},
				},
			},
		},
		),
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *API) newStateField() *graphql.Field {
	return &graphql.Field{
		Name: "State",
		Type: StateType,
		Args: graphql.FieldConfigArgument{
			"PeerID": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "account id of a provider",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			peerIDStr := p.Args["PeerID"].(string)
			peerID, err := peer.Decode(peerIDStr)
			if err != nil {
				return nil, err
			}
			providerState, err := a.core.StateTree.GetProviderStateByPeerID(peerID)
			if err != nil {
				return nil, err
			}

			return providerState, nil
		},
	}
}

func (a *API) newSnapshotField() *graphql.Field {
	return &graphql.Field{
		Name: "SnapShot",
		Type: SnapShotType,
		Args: graphql.FieldConfigArgument{
			"cid": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "cid of the snapshot",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			cidStr := p.Args["cid"].(string)
			snapshotCid, err := cid.Decode(cidStr)
			if err != nil {
				return nil, err
			}
			ss, err := a.core.StateTree.GetSnapShot(snapshotCid)
			if err != nil {
				return nil, err
			}
			return ss, nil
		},
	}
}
