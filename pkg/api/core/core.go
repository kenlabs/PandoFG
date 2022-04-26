package core

import (
	badger "github.com/dgraph-io/badger/v3"
	"github.com/ipfs/go-datastore/sync"
	"github.com/ipfs/go-ds-leveldb"
	"github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipld/go-ipld-prime"
	"github.com/kenlabs/pandofg/pkg/legs"
	"github.com/kenlabs/pandofg/pkg/lotus"
	"github.com/kenlabs/pandofg/pkg/metadata"
	"github.com/kenlabs/pandofg/pkg/registry"
	"github.com/kenlabs/pandofg/pkg/statetree"
	"go.mongodb.org/mongo-driver/mongo"
)

type Core struct {
	MetaManager   *metadata.MetaManager
	StateTree     *statetree.StateTree
	LotusDiscover *lotus.Discoverer
	Registry      *registry.Registry
	LegsCore      *legs.Core
	StoreInstance *StoreInstance
	LinkSystem    *ipld.LinkSystem
}

type StoreInstance struct {
	DataStore      *leveldb.Datastore
	MutexDataStore *sync.MutexDatastore
	BlockStore     blockstore.Blockstore
	CacheStore     *badger.DB
	MetaStore      *mongo.Client
}
