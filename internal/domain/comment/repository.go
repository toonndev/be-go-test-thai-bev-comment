package comment

import "context"

type Repository interface {
	Create(ctx context.Context, comment *Comment) (*Comment, error)
	FindAll(ctx context.Context) ([]*Comment, error)
}
