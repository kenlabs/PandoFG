package metastore

import "context"

type ProviderLocations interface {
	Get(ctx context.Context, minerID string)
	GetAll(ctx context.Context)
}
