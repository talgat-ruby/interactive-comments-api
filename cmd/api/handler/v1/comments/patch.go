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

type PatchRequestParam struct {
	ID *int `param:"id" validate:"required,gt=0"`
}

type PatchRequestQuery struct {
	User *string `query:"user" validate:"required"`
}

type PatchRequestBody struct {
	Content string `xml:"content" json:"content" form:"content" validate:"required"`
}

func (h *Handler) Edit(c echo.Context) error {
	ctx := c.Request().Context()
	h.log.InfoContext(ctx, "start Edit", "path", c.Path())

	reqParam := new(PatchRequestParam)
	if err := (&echo.DefaultBinder{}).BindPathParams(c, reqParam); err != nil {
		h.log.ErrorContext(
			ctx,
			"fail Delete:: body binding error",
			"path", c.Path(),
			"error", err,
		)
		return c.JSON(http.StatusBadRequest, response.ErrorWithMessage{Error: response.WithMessage{Message: err.Error()}})
	}

	if validationError := h.patchRequestParamValidationErrors(ctx, reqParam); validationError != nil {
		h.log.ErrorContext(
			ctx,
			"fail Delete:: validation errors",
			"path", c.Path(),
		)
		return c.JSON(http.StatusBadRequest, response.Error{Error: validationError})
	}

	reqQuery := new(PatchRequestQuery)
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, reqQuery); err != nil {
		h.log.ErrorContext(
			ctx,
			"fail Edit:: body binding error",
			"path", c.Path(),
			"error", err,
		)
		return c.JSON(http.StatusBadRequest, response.ErrorWithMessage{Error: response.WithMessage{Message: err.Error()}})
	}

	if validationError := h.patchRequestQueryValidationErrors(ctx, reqQuery); validationError != nil {
		h.log.ErrorContext(
			ctx,
			"fail Delete:: validation errors",
			"path", c.Path(),
		)
		return c.JSON(http.StatusBadRequest, response.Error{Error: validationError})
	}

	reqBody, err := h.patchRequestBody(ctx, c)
	if err != nil {
		h.log.ErrorContext(
			ctx,
			"fail Edit:: body binding error",
			"path", c.Path(),
			"error", err,
		)
		return c.JSON(http.StatusBadRequest, response.ErrorWithMessage{Error: response.WithMessage{Message: err.Error()}})
	}

	if validationError := h.patchRequestValidationErrors(ctx, reqBody); validationError != nil {
		h.log.ErrorContext(
			ctx,
			"fail Edit:: validation errors",
			"path", c.Path(),
		)
		return c.JSON(http.StatusBadRequest, response.Error{Error: validationError})
	}

	dbInput := patchDBInput(reqBody, reqParam.ID, reqQuery.User)
	if err := h.db.UpdateComment(ctx, dbInput); err != nil {
		h.log.ErrorContext(
			ctx,
			"fail Edit:: db add fail",
			"path", c.Path(),
			"error", err,
		)
		return c.JSON(http.StatusBadRequest, response.ErrorWithMessage{Error: response.WithMessage{Message: err.Error()}})
	}

	h.log.InfoContext(ctx, "success Edit", "path", c.Path())
	return c.NoContent(http.StatusNoContent)
}

type patchParamValidationError struct {
	ID *string `json:"id,omitempty"`
}

func (h *Handler) patchRequestParamValidationErrors(_ context.Context, reqParam *PatchRequestParam) *patchParamValidationError {
	if err := h.validate.Struct(reqParam); err != nil {
		vErr := new(patchParamValidationError)

		for _, err := range err.(validator.ValidationErrors) {
			switch err.StructField() {
			case "ID":
				vErr.ID = utils.ToPtr("id is invalid")
			}
		}

		return vErr
	}

	return nil
}

type patchQueryValidationError struct {
	User *string `json:"user,omitempty"`
}

func (h *Handler) patchRequestQueryValidationErrors(_ context.Context, reqParam *PatchRequestQuery) *patchQueryValidationError {
	if err := h.validate.Struct(reqParam); err != nil {
		vErr := new(patchQueryValidationError)

		for _, err := range err.(validator.ValidationErrors) {
			switch err.StructField() {
			case "User":
				vErr.User = utils.ToPtr("user is invalid")
			}
		}

		return vErr
	}

	return nil
}

type patchValidationError struct {
	Content *string `json:"content,omitempty"`
}

func (h *Handler) patchRequestBody(_ context.Context, c echo.Context) (*PatchRequestBody, error) {
	reqBody := new(PatchRequestBody)
	if err := (&echo.DefaultBinder{}).BindBody(c, reqBody); err != nil {
		return nil, err
	}

	return reqBody, nil
}

func (h *Handler) patchRequestValidationErrors(_ context.Context, reqBody *PatchRequestBody) *patchValidationError {
	if err := h.validate.Struct(reqBody); err != nil {
		vErr := new(patchValidationError)

		for _, err := range err.(validator.ValidationErrors) {
			switch err.StructField() {
			case "Content":
				vErr.Content = utils.ToPtr("content is required")
			}
		}

		return vErr
	}

	return nil
}

func patchDBInput(reqBody *PatchRequestBody, id *int, username *string) *model.UpdateCommentInput {
	inp := new(model.UpdateCommentInput)

	if reqBody == nil {
		return inp
	}

	inp.ID = id
	inp.Author = username
	inp.Content = reqBody.Content

	return inp
}
