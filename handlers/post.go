package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lunebakami/futtodos-api/models"
	"github.com/lunebakami/futtodos-api/storage"
)

type PostHandler struct {
	storage *storage.PostStorage
}

func NewPostHandler(storage *storage.PostStorage) *PostHandler {
	return &PostHandler{storage: storage}
}

func (h *PostHandler) CreatePost(c echo.Context) error {
  post := new(models.BlogPost)

  if err := c.Bind(post); err != nil {
    return err
  }

  err := h.storage.Create(post)
  if err != nil {
    return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
  }

  return c.JSON(http.StatusCreated, post)
}

func (h *PostHandler) GetAllPost(c echo.Context) error {
  posts, err := h.storage.GetAll()
  if err != nil {
    return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
  }

  return c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) GetPost(c echo.Context) error {
  id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

  post, err := h.storage.GetByID(id)

  if err != nil {
    return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
  }

  return c.JSON(http.StatusOK, post)
}

func (h *PostHandler) UpdatePost(c echo.Context) error {
  id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
  updatedPost := new(models.BlogPost)
  if err := c.Bind(updatedPost); err != nil {
    return err
  }

  err := h.storage.Update(id, updatedPost)

  if err != nil {
    return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
  }

  return c.JSON(http.StatusOK, updatedPost)
}

func (h *PostHandler) DeletePost(c echo.Context) error {
  id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

  err := h.storage.Delete(id)
  if err != nil {
    return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
  }

  return c.NoContent(http.StatusNoContent)
}
