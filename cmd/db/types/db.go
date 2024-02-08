package types

import (
	"context"

	"github.com/talgat-ruby/interactive-comments-api/cmd/db/model"
)

type DB interface {
	AddForm(ctx context.Context, input model.Form) error
	GetComments(ctx context.Context, username string) ([]*model.Comment, error)
}
