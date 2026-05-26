package handler

import (
	"net/http"

	"github.com/bsyrlhabibi/airdrop/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardHandler struct {
	DB *gorm.DB
}

func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{DB: db}
}

type DashboardSummary struct {
	TotalAirdrops  int64 `json:"total_airdrops" example:"12"`
	ActiveAirdrops int64 `json:"active_airdrops" example:"8"`
	TotalTasks     int64 `json:"total_tasks" example:"45"`
	CompletedTasks int64 `json:"completed_tasks" example:"30"`
	PendingTasks   int64 `json:"pending_tasks" example:"15"`
	TotalWallets   int64 `json:"total_wallets" example:"3"`
}

// Dashboard Summary godoc
// @Summary      Get dashboard stats
// @Description  Get summary statistics for the authenticated user
// @Tags         Dashboard
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  DashboardSummary
// @Failure      401  {object}  map[string]string
// @Router       /api/dashboard [get]
func (h *DashboardHandler) Summary(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var totalAirdrops int64
	h.DB.Model(&model.Airdrop{}).Where("user_id = ?", userID).Count(&totalAirdrops)

	var activeAirdrops int64
	h.DB.Model(&model.Airdrop{}).Where("user_id = ? AND status = ?", userID, "active").Count(&activeAirdrops)

	var totalTasks int64
	h.DB.Model(&model.Task{}).
		Joins("JOIN airdrops ON airdrops.id = tasks.airdrop_id").
		Where("airdrops.user_id = ?", userID).
		Count(&totalTasks)

	var completedTasks int64
	h.DB.Model(&model.Task{}).
		Joins("JOIN airdrops ON airdrops.id = tasks.airdrop_id").
		Where("airdrops.user_id = ? AND tasks.is_completed = ?", userID, true).
		Count(&completedTasks)

	var totalWallets int64
	h.DB.Model(&model.Wallet{}).Where("user_id = ?", userID).Count(&totalWallets)

	c.JSON(http.StatusOK, DashboardSummary{
		TotalAirdrops:  totalAirdrops,
		ActiveAirdrops: activeAirdrops,
		TotalTasks:     totalTasks,
		CompletedTasks: completedTasks,
		PendingTasks:   totalTasks - completedTasks,
		TotalWallets:   totalWallets,
	})
}
