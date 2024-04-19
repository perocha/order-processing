package cosmosdb

import (
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

// Database client interface, used to create a new container
type DatabaseClientInterface interface {
	NewContainer(id string) (ContainerClientInterface, error)
}

type CosmosDatabase struct {
	database *azcosmos.DatabaseClient
}

func (d *CosmosDatabase) NewContainer(id string) (ContainerClientInterface, error) {
	container, err := d.database.NewContainer(id)
	if err != nil {
		return nil, err
	}
	return &CosmosContainer{container: container}, nil
}
