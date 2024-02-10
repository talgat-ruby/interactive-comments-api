package comments

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/talgat-ruby/interactive-comments-api/cmd/db/model"
	"github.com/talgat-ruby/interactive-comments-api/internal/response"
	"github.com/talgat-ruby/interactive-comments-api/pkg/utils"
)

type PostRequestQuery struct {
	User *string `query:"user"`
}

type PostRequestBody struct {
	ParentID  *int    `xml:"parentId" json:"parentId,omitempty" form:"parentId" validate:"omitempty,gt=0"`
	Addressee *string `xml:"addressee" json:"addressee,omitempty" form:"addressee" validate:"required_with=ParentID,omitempty,gt=0"`
	Content   string  `xml:"content" json:"content" form:"content" validate:"required"`
}

func (h *Handler) Add(c echo.Context) error {
	ctx := c.Request().Context()
	h.log.InfoContext(ctx, "start Add", "path", c.Path())

	reqQuery := new(PostRequestQuery)
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, reqQuery); err != nil {
		h.log.ErrorContext(
			ctx,
			"fail Add:: body binding error",
			"path", c.Path(),
			"error", err,
		)
		return c.JSON(http.StatusBadRequest, response.ErrorWithMessage{Error: response.WithMessage{Message: err.Error()}})
	}

	reqBody, err := h.postRequestBody(ctx, c)
	if err != nil {
		h.log.ErrorContext(
			ctx,
			"fail Add:: body binding error",
			"path", c.Path(),
			"error", err,
		)
		return c.JSON(http.StatusBadRequest, response.ErrorWithMessage{Error: response.WithMessage{Message: err.Error()}})
	}

	if validationError := h.postRequestValidationErrors(ctx, reqBody); validationError != nil {
		h.log.ErrorContext(
			ctx,
			"fail Add:: validation errors",
			"path", c.Path(),
		)
		return c.JSON(http.StatusBadRequest, response.Error{Error: validationError})
	}

	dbInput := postDBInput(reqBody, reqQuery.User)
	if err := h.db.CreateComment(ctx, dbInput); err != nil {
		h.log.ErrorContext(
			ctx,
			"fail Add:: db add fail",
			"path", c.Path(),
			"error", err,
		)
		return c.JSON(http.StatusBadRequest, response.ErrorWithMessage{Error: response.WithMessage{Message: err.Error()}})
	}

	h.log.InfoContext(ctx, "success Add", "path", c.Path())
	return c.NoContent(http.StatusNoContent)
}

type postValidationError struct {
	Content   *string `json:"content,omitempty"`
	ParentID  *string `json:"parentId,omitempty"`
	Addressee *string `json:"addressee,omitempty"`
}

func (h *Handler) postRequestBody(_ context.Context, c echo.Context) (*PostRequestBody, error) {
	reqBody := new(PostRequestBody)
	if err := (&echo.DefaultBinder{}).BindBody(c, reqBody); err != nil {
		return nil, err
	}

	return reqBody, nil
}

func (h *Handler) postRequestValidationErrors(_ context.Context, reqBody *PostRequestBody) *postValidationError {
	if err := h.validate.Struct(reqBody); err != nil {
		vErr := new(postValidationError)

		for _, err := range err.(validator.ValidationErrors) {
			switch err.StructField() {
			case "ParentID":
				vErr.ParentID = utils.ToPtr("parentId is invalid")
			case "ReplyID":
				vErr.Addressee = utils.ToPtr("addressee is invalid, is required if parentId presented")
			case "Content":
				vErr.Content = utils.ToPtr("content is required")
			}
		}

		return vErr
	}

	return nil
}

func postDBInput(reqBody *PostRequestBody, username *string) *model.CreateCommentInput {
	inp := new(model.CreateCommentInput)

	if reqBody == nil {
		return inp
	}

	inp.Author = username
	inp.Content = reqBody.Content
	inp.ParentID = reqBody.ParentID
	inp.Addressee = reqBody.Addressee

	return inp
}
