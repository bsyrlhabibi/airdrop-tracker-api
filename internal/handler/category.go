package handler

import (
	"net/http"
	"strconv"

	"github.com/bsyrlhabibi/airdrop/internal/model"
	"github.com/bsyrlhabibi/airdrop/internal/repository"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	Repo *repository.CategoryRepo
}

func NewCategoryHandler(repo *repository.CategoryRepo) *CategoryHandler {
	return &CategoryHandler{Repo: repo}
}

type CreateCategoryRequest struct {
	Name  string `json:"name" binding:"required" example:"Bridge"`
	Color string `json:"color" example:"#3B82F6"`
}

type UpdateCategoryRequest struct {
	Name  string `json:"name" example:"Bridge"`
	Color string `json:"color" example:"#3B82F6"`
}

// List godoc
// @Summary      List categories
// @Description  Get all categories for authenticated user
// @Tags         Categories
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  model.Category
// @Router       /api/categories [get]
func (h *CategoryHandler) List(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	categories, err := h.Repo.FindByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// Create godoc
// @Summary      Create category
// @Description  Create a new task category
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body CreateCategoryRequest true "Category data"
// @Success      201  {object}  model.Category
// @Router       /api/categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Color == "" {
		req.Color = "#6B7280"
	}

	cat := &model.Category{
		UserID: userID,
		Name:   req.Name,
		Color:  req.Color,
	}

	if err := h.Repo.Create(cat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cat)
}

// Update godoc
// @Summary      Update category
// @Description  Update category name or color
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int                  true  "Category ID"
// @Param        body body      UpdateCategoryRequest true  "Updated data"
// @Success      200  {object}  model.Category
// @Router       /api/categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	existing, err := h.Repo.FindByID(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Color != "" {
		existing.Color = req.Color
	}

	h.Repo.Update(existing)
	c.JSON(http.StatusOK, existing)
}

// Delete godoc
// @Summary      Delete category
// @Description  Remove a category
// @Tags         Categories
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  map[string]string
// @Router       /api/categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.Repo.Delete(uint(id), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
