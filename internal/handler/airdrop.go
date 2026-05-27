package handler

import (
	"net/http"
	"strconv"

	"github.com/bsyrlhabibi/airdrop/internal/model"
	"github.com/bsyrlhabibi/airdrop/internal/repository"
	"github.com/gin-gonic/gin"
)

type AirdropHandler struct {
	Repo *repository.AirdropRepo
}

func NewAirdropHandler(repo *repository.AirdropRepo) *AirdropHandler {
	return &AirdropHandler{Repo: repo}
}

type CreateAirdropRequest struct {
	Name      string `json:"name" example:"zkSync"`
	Chain     string `json:"chain" example:"Ethereum"`
	Category  string `json:"category" example:"rumored"`
	Priority  string `json:"priority" example:"high"`
	URL       string `json:"url" example:"https://zksync.io"`
	DateStart string `json:"date_start" example:"2025-01-01"`
	DateEnd   string `json:"date_end" example:"2025-12-31"`
	Notes     string `json:"notes" example:"Bridge weekly"`
}

type UpdateAirdropRequest struct {
	Name      string `json:"name"`
	Chain     string `json:"chain"`
	Category  string `json:"category"`
	Priority  string `json:"priority"`
	Status    string `json:"status"`
	URL       string `json:"url"`
	DateStart string `json:"date_start"`
	DateEnd   string `json:"date_end"`
	Notes     string `json:"notes"`
}

// List Airdrops godoc
// @Summary      List all airdrops
// @Description  Get all global airdrops for the authenticated user
// @Tags         Airdrops
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   model.Airdrop
// @Failure      401  {object}  map[string]string
// @Router       /api/airdrops [get]
func (h *AirdropHandler) List(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	airdrops, err := h.Repo.FindByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, airdrops)
}

// Create Airdrop godoc
// @Summary      Create airdrop
// @Description  Add new global airdrop to the catalog
// @Tags         Airdrops
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body CreateAirdropRequest true "Airdrop data"
// @Success      201  {object}  model.Airdrop
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /api/airdrops [post]
func (h *AirdropHandler) Create(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var req CreateAirdropRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	airdrop := &model.Airdrop{
		UserID:    userID,
		Name:      req.Name,
		Chain:     req.Chain,
		Category:  req.Category,
		Priority:  req.Priority,
		URL:       req.URL,
		DateStart: parseDate(req.DateStart),
		DateEnd:   parseDate(req.DateEnd),
		Notes:     req.Notes,
	}

	if err := h.Repo.Create(airdrop); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, airdrop)
}

// Get Airdrop godoc
// @Summary      Get airdrop detail
// @Description  Get single airdrop by ID
// @Tags         Airdrops
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Airdrop ID"
// @Success      200  {object}  model.Airdrop
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/airdrops/{id} [get]
func (h *AirdropHandler) Get(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	a, err := h.Repo.FindByID(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, a)
}

// Update Airdrop godoc
// @Summary      Update airdrop
// @Description  Update existing airdrop
// @Tags         Airdrops
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int                  true  "Airdrop ID"
// @Param        body body      UpdateAirdropRequest true  "Updated data"
// @Success      200  {object}  model.Airdrop
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/airdrops/{id} [put]
func (h *AirdropHandler) Update(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	existing, err := h.Repo.FindByID(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	var req UpdateAirdropRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Apply non-empty fields
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Chain != "" {
		existing.Chain = req.Chain
	}
	if req.Category != "" {
		existing.Category = req.Category
	}
	if req.Priority != "" {
		existing.Priority = req.Priority
	}
	if req.Status != "" {
		existing.Status = req.Status
	}
	if req.URL != "" {
		existing.URL = req.URL
	}
	if req.Notes != "" {
		existing.Notes = req.Notes
	}
	// Date fields: use parseDate helper for YYYY-MM-DD support
	if req.DateStart != "" {
		existing.DateStart = parseDate(req.DateStart)
	}
	if req.DateEnd != "" {
		existing.DateEnd = parseDate(req.DateEnd)
	}

	h.Repo.Update(existing)
	c.JSON(http.StatusOK, existing)
}

// Delete Airdrop godoc
// @Summary      Delete airdrop
// @Description  Remove airdrop from the global catalog
// @Tags         Airdrops
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Airdrop ID"
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/airdrops/{id} [delete]
func (h *AirdropHandler) Delete(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.Repo.Delete(uint(id), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
