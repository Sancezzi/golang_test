package handler

import (
	"context"
	"errors"
	"strconv"
	"test-api/internal/common/command"
	"test-api/internal/common/errs"
	"test-api/internal/http/responses"
	"test-api/internal/service/model"

	"github.com/gofiber/fiber/v3"
)

type ArticleService interface {
	Update(ctx context.Context, cmd command.UpdateArticle, id int64) (model.Article, error)
	List(ctx context.Context) ([]model.Article, error)
}

type articleHandler struct {
	articleService ArticleService
}

func NewArticleHandler(
	articleService ArticleService,
) *articleHandler {
	return &articleHandler{
		articleService: articleService,
	}
}

func (h *articleHandler) Register(app *fiber.App) {
	app.Post("/edit/:Id", h.Edit)
	app.Get("/list", h.List)
}

func (h *articleHandler) Edit(c fiber.Ctx) error {
	ctx := c.UserContext()

	idParam := c.Params("Id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	cmd := command.UpdateArticle{}
	if err := c.Bind().Body(&cmd); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	article, err := h.articleService.Update(ctx, cmd, int64(id))
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(responses.NewArticleResponse(article))
}

func (h *articleHandler) List(c fiber.Ctx) error {
	ctx := c.UserContext()

	list, err := h.articleService.List(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	resp := make([]responses.Article, 0, len(list))
	for _, v := range list {
		resp = append(resp, responses.NewArticleResponse(v))
	}

	return c.JSON(map[string]any{
		"Success": true,
		"News":    resp,
	})
}
