package usecase

import (
	"context"
	"log"
)

type DeleteOrder struct {
	// Define your dependencies such as repositories
}

func (do *DeleteOrder) Execute(ctx context.Context, orderID string) error {
	log.Printf("DeleteOrder::Execute::order=%v", orderID)
	return nil
}
