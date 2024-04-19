package cosmosdb

import "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"

// CosmosDB client interface
type ClientInterface interface {
	NewDatabase(id string) (DatabaseClientInterface, error)
	Endpoint() string
}

// ClientInterface implementation
type CosmosClient struct {
	client *azcosmos.Client
}

func (c *CosmosClient) NewDatabase(id string) (DatabaseClientInterface, error) {
	db, err := c.client.NewDatabase(id)
	if err != nil {
		return nil, err
	}
	return &CosmosDatabase{database: db}, nil
}

func (c *CosmosClient) Endpoint() string {
	return c.client.Endpoint()
}
