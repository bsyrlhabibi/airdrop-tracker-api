package handler

import (
	"net/http"
	"strconv"

	"github.com/bsyrlhabibi/airdrop/internal/repository"
	"github.com/gin-gonic/gin"
)

type AccountAirdropHandler struct {
	Repo *repository.AccountAirdropRepo
}

func NewAccountAirdropHandler(repo *repository.AccountAirdropRepo) *AccountAirdropHandler {
	return &AccountAirdropHandler{Repo: repo}
}

type UpdateAccountAirdropRequest struct {
	Status string `json:"status" example:"active"`
	Notes  string `json:"notes" example:"Focus on bridging"`
}

// Get AccountAirdrop godoc
// @Summary      Get account-airdrop detail
// @Description  Get a specific account-airdrop with tasks
// @Tags         AccountAirdrops
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Account Airdrop ID"
// @Success      200  {object}  model.AccountAirdrop
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/account-airdrops/{id} [get]
func (h *AccountAirdropHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	aa, err := h.Repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account airdrop not found"})
		return
	}
	c.JSON(http.StatusOK, aa)
}

// Update AccountAirdrop godoc
// @Summary      Update account-airdrop
// @Description  Update status or notes for an account-airdrop assignment
// @Tags         AccountAirdrops
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int                           true  "Account Airdrop ID"
// @Param        body body      UpdateAccountAirdropRequest   true  "Updated data"
// @Success      200  {object}  model.AccountAirdrop
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/account-airdrops/{id} [put]
func (h *AccountAirdropHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	aa, err := h.Repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account airdrop not found"})
		return
	}

	var req UpdateAccountAirdropRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Status != "" {
		aa.Status = req.Status
	}
	if req.Notes != "" {
		aa.Notes = req.Notes
	}

	if err := h.Repo.Update(aa); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, aa)
}
