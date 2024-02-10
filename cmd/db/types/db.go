package types

import (
	"context"

	"github.com/talgat-ruby/interactive-comments-api/cmd/db/model"
)

type DB interface {
	ReadComments(ctx context.Context, username string) ([]*model.Comment, error)
	InsertComment(ctx context.Context, input model.InsertCommentInput) error
}
