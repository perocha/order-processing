package storage

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/perocha/order-processing/pkg/domain"
)

type CosmosDBOrderRepository struct {
	client *azcosmos.Client
}

// Initialize CosmosDB repository using the provided connection string
func NewCosmosDBOrderRepository(connectionString string) (*CosmosDBOrderRepository, error) {
	client, err := azcosmos.NewClientFromConnectionString(connectionString, nil)
	log.Printf("CosmosDBOrderRepository::NewCosmosDBOrderRepository::Client=%v::Error=%v", client, err)

	return nil, nil
	/*
		if err != nil {
			return nil, err
		}

		return &CosmosDBOrderRepository{
			client: client,
		}, nil
	*/
}

// Creates a new order in CosmosDB
func (r *CosmosDBOrderRepository) CreateOrder(ctx context.Context, order domain.Order) error {
	log.Printf("CosmosDBOrderRepository::CreateOrder::OrderID=%v", order.OrderID)
	return nil
}

// Updates an existing order in CosmosDB
func (r *CosmosDBOrderRepository) UpdateOrder(ctx context.Context, order domain.Order) error {
	log.Printf("CosmosDBOrderRepository::UpdateOrder::OrderID=%v", order.OrderID)
	return nil
}

// Deletes an order from CosmosDB
func (r *CosmosDBOrderRepository) DeleteOrder(ctx context.Context, orderID string) error {
	log.Printf("CosmosDBOrderRepository::DeleteOrder::OrderID=%v", orderID)
	return nil
}
