package port

import "context"

// TxManager abstracts database transaction management.
type TxManager interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}
