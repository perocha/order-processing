package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

// CosmosDB container, with all methods to interact with the container (CRUD operations)
type ContainerClientInterface interface {
	CreateItem(ctx context.Context, partitionKey azcosmos.PartitionKey, item interface{}, options *azcosmos.ItemOptions) (azcosmos.ItemResponse, error)
	UpsertItem(ctx context.Context, partitionKey azcosmos.PartitionKey, item interface{}, options *azcosmos.ItemOptions) (azcosmos.ItemResponse, error)
	DeleteItem(ctx context.Context, partitionKey azcosmos.PartitionKey, id string, options *azcosmos.ItemOptions) (azcosmos.ItemResponse, error)
	ReadItem(ctx context.Context, partitionKey azcosmos.PartitionKey, id string, options *azcosmos.ItemOptions) (azcosmos.ItemResponse, error)
}

type CosmosContainer struct {
	container *azcosmos.ContainerClient
}

func (c *CosmosContainer) CreateItem(ctx context.Context, partitionKey azcosmos.PartitionKey, item interface{}, options *azcosmos.ItemOptions) (azcosmos.ItemResponse, error) {
	itemBytes, ok := item.([]byte)
	if !ok {
		return azcosmos.ItemResponse{}, fmt.Errorf("item is not of type []byte")
	}
	return c.container.CreateItem(ctx, partitionKey, itemBytes, options)
}

func (c *CosmosContainer) UpsertItem(ctx context.Context, partitionKey azcosmos.PartitionKey, item interface{}, options *azcosmos.ItemOptions) (azcosmos.ItemResponse, error) {
	itemBytes, ok := item.([]byte)
	if !ok {
		return azcosmos.ItemResponse{}, fmt.Errorf("item is not of type []byte")
	}
	return c.container.UpsertItem(ctx, partitionKey, itemBytes, options)
}

func (c *CosmosContainer) DeleteItem(ctx context.Context, partitionKey azcosmos.PartitionKey, id string, options *azcosmos.ItemOptions) (azcosmos.ItemResponse, error) {
	return c.container.DeleteItem(ctx, partitionKey, id, options)
}

func (c *CosmosContainer) ReadItem(ctx context.Context, partitionKey azcosmos.PartitionKey, id string, options *azcosmos.ItemOptions) (azcosmos.ItemResponse, error) {
	return c.container.ReadItem(ctx, partitionKey, id, options)
}
