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
	TotalAirdrops    int64          `json:"total_airdrops"`
	ActiveAirdrops   int64          `json:"active_airdrops"`
	UpcomingAirdrops int64          `json:"upcoming_airdrops"`
	EndedAirdrops    int64          `json:"ended_airdrops"`
	MissedAirdrops   int64          `json:"missed_airdrops"`
	TotalTasks       int64          `json:"total_tasks"`
	CompletedTasks   int64          `json:"completed_tasks"`
	PendingTasks     int64          `json:"pending_tasks"`
	OngoingTasks     int64          `json:"ongoing_tasks"`
	MissedTasks      int64          `json:"missed_tasks"`
	TotalWallets     int64          `json:"total_wallets"`
	TotalAccounts    int64          `json:"total_accounts"`
	Accounts         []AccountStats `json:"accounts,omitempty"`
}

type AccountStats struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Color          string `json:"color"`
	TotalAirdrops  int64  `json:"total_airdrops"`
	ActiveAirdrops int64  `json:"active_airdrops"`
	TotalTasks     int64  `json:"total_tasks"`
	CompletedTasks int64  `json:"completed_tasks"`
	PendingTasks   int64  `json:"pending_tasks"`
	OngoingTasks   int64  `json:"ongoing_tasks"`
	MissedTasks    int64  `json:"missed_tasks"`
	TotalWallets   int64  `json:"total_wallets"`
}

// Dashboard Summary godoc
// @Summary      Get dashboard stats
// @Description  Get summary statistics with per-account breakdown
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

	var upcomingAirdrops int64
	h.DB.Model(&model.Airdrop{}).Where("user_id = ? AND status = ?", userID, "upcoming").Count(&upcomingAirdrops)

	var endedAirdrops int64
	h.DB.Model(&model.Airdrop{}).Where("user_id = ? AND status = ?", userID, "end").Count(&endedAirdrops)

	var missedAirdrops int64
	h.DB.Model(&model.Airdrop{}).Where("user_id = ? AND status = ?", userID, "missed").Count(&missedAirdrops)

	var totalTasks int64
	h.DB.Model(&model.Task{}).
		Joins("JOIN account_airdrops ON account_airdrops.id = tasks.account_airdrop_id").
		Joins("JOIN accounts ON accounts.id = account_airdrops.account_id").
		Where("accounts.user_id = ?", userID).
		Count(&totalTasks)

	var completedTasks int64
	h.DB.Model(&model.Task{}).
		Joins("JOIN account_airdrops ON account_airdrops.id = tasks.account_airdrop_id").
		Joins("JOIN accounts ON accounts.id = account_airdrops.account_id").
		Where("accounts.user_id = ? AND tasks.status = ?", userID, "finish").
		Count(&completedTasks)

	var ongoingTasks int64
	h.DB.Model(&model.Task{}).
		Joins("JOIN account_airdrops ON account_airdrops.id = tasks.account_airdrop_id").
		Joins("JOIN accounts ON accounts.id = account_airdrops.account_id").
		Where("accounts.user_id = ? AND tasks.status = ?", userID, "ongoing").
		Count(&ongoingTasks)

	var missedTasks int64
	h.DB.Model(&model.Task{}).
		Joins("JOIN account_airdrops ON account_airdrops.id = tasks.account_airdrop_id").
		Joins("JOIN accounts ON accounts.id = account_airdrops.account_id").
		Where("accounts.user_id = ? AND tasks.status = ?", userID, "missed").
		Count(&missedTasks)

	var totalWallets int64
	h.DB.Model(&model.Wallet{}).Where("user_id = ?", userID).Count(&totalWallets)

	var totalAccounts int64
	h.DB.Model(&model.Account{}).Where("user_id = ?", userID).Count(&totalAccounts)

	// Per-account stats
	var accounts []model.Account
	h.DB.Where("user_id = ?", userID).Order("created_at ASC").Find(&accounts)

	var accountStats []AccountStats
	for _, acc := range accounts {
		var accAirdrops int64
		h.DB.Model(&model.AccountAirdrop{}).Where("account_id = ?", acc.ID).Count(&accAirdrops)

		var accActiveAirdrops int64
		h.DB.Model(&model.AccountAirdrop{}).Where("account_id = ? AND status = ?", acc.ID, "active").Count(&accActiveAirdrops)

		var accTotalTasks int64
		h.DB.Model(&model.Task{}).
			Joins("JOIN account_airdrops ON account_airdrops.id = tasks.account_airdrop_id").
			Where("account_airdrops.account_id = ?", acc.ID).
			Count(&accTotalTasks)

		var accCompletedTasks int64
		h.DB.Model(&model.Task{}).
			Joins("JOIN account_airdrops ON account_airdrops.id = tasks.account_airdrop_id").
			Where("account_airdrops.account_id = ? AND tasks.status = ?", acc.ID, "finish").
			Count(&accCompletedTasks)

		var accOngoingTasks int64
		h.DB.Model(&model.Task{}).
			Joins("JOIN account_airdrops ON account_airdrops.id = tasks.account_airdrop_id").
			Where("account_airdrops.account_id = ? AND tasks.status = ?", acc.ID, "ongoing").
			Count(&accOngoingTasks)

		var accMissedTasks int64
		h.DB.Model(&model.Task{}).
			Joins("JOIN account_airdrops ON account_airdrops.id = tasks.account_airdrop_id").
			Where("account_airdrops.account_id = ? AND tasks.status = ?", acc.ID, "missed").
			Count(&accMissedTasks)

		var accWallets int64
		h.DB.Model(&model.Wallet{}).Where("account_id = ?", acc.ID).Count(&accWallets)

		accountStats = append(accountStats, AccountStats{
			ID:             acc.ID,
			Name:           acc.Name,
			Color:          acc.Color,
			TotalAirdrops:  accAirdrops,
			ActiveAirdrops: accActiveAirdrops,
			TotalTasks:     accTotalTasks,
			CompletedTasks: accCompletedTasks,
			PendingTasks:   accTotalTasks - accCompletedTasks - accOngoingTasks - accMissedTasks,
			OngoingTasks:   accOngoingTasks,
			MissedTasks:    accMissedTasks,
			TotalWallets:   accWallets,
		})
	}

	c.JSON(http.StatusOK, DashboardSummary{
		TotalAirdrops:    totalAirdrops,
		ActiveAirdrops:   activeAirdrops,
		UpcomingAirdrops: upcomingAirdrops,
		EndedAirdrops:    endedAirdrops,
		MissedAirdrops:   missedAirdrops,
		TotalTasks:       totalTasks,
		CompletedTasks:   completedTasks,
		PendingTasks:     totalTasks - completedTasks - ongoingTasks - missedTasks,
		OngoingTasks:     ongoingTasks,
		MissedTasks:      missedTasks,
		TotalWallets:     totalWallets,
		TotalAccounts:    totalAccounts,
		Accounts:         accountStats,
	})
}
