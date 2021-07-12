package client

import (
	"context"
	"fmt"
	"github.com/apirator/apirator/internal/inventory"
)

func (k *Kubernetes) Apply(ctx context.Context, inv inventory.Object) error {
	for _, obj := range inv.Create {
		if err := k.Create(ctx, obj); err != nil {
			return fmt.Errorf("failed to create %T %q: %w", obj, obj.GetName(), err)
		}
	}

	for _, obj := range inv.Update {
		if err := k.Update(ctx, obj); err != nil {
			return fmt.Errorf("failed to update %T %q: %w", obj, obj.GetName(), err)
		}
	}

	for _, obj := range inv.Delete {
		if err := k.Delete(ctx, obj); err != nil {
			return fmt.Errorf("failed to delete %T %q: %w", obj, obj.GetName(), err)
		}
	}

	return nil
}
