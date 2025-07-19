package wallets

import "context"

type Wallet interface {
	ActivateCheque(ctx context.Context, chequeID string) error
}
