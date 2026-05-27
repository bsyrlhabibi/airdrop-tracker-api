package handler

import (
	"net/http"
	"strconv"

	"github.com/bsyrlhabibi/airdrop/internal/model"
	"github.com/bsyrlhabibi/airdrop/internal/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AccountHandler struct {
	Repo        *repository.AccountRepo
	AirdropRepo *repository.AirdropRepo
	AARepo      *repository.AccountAirdropRepo
	DB          *gorm.DB
}

func NewAccountHandler(repo *repository.AccountRepo, airdropRepo *repository.AirdropRepo, aaRepo *repository.AccountAirdropRepo, db *gorm.DB) *AccountHandler {
	return &AccountHandler{Repo: repo, AirdropRepo: airdropRepo, AARepo: aaRepo, DB: db}
}

type CreateAccountRequest struct {
	Name  string `json:"name" example:"Akun 1"`
	Color string `json:"color" example:"#3B82F6"`
	Notes string `json:"notes" example:"Main sybil account"`
}

type UpdateAccountRequest struct {
	Name  string `json:"name" example:"Akun 1"`
	Color string `json:"color" example:"#EF4444"`
	Notes string `json:"notes" example:"Updated notes"`
}

type AssignAirdropRequest struct {
	AirdropID uint   `json:"airdrop_id" binding:"required" example:"1"`
	Notes     string `json:"notes" example:"Focus on bridging"`
}

type CloneAccountRequest struct {
	Name  string `json:"name" example:"Akun 2"`
	Color string `json:"color" example:"#EF4444"`
}

// List Accounts godoc
// @Summary      List all accounts
// @Description  Get all sybil accounts for authenticated user
// @Tags         Accounts
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   model.Account
// @Failure      401  {object}  map[string]string
// @Router       /api/accounts [get]
func (h *AccountHandler) List(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	accounts, err := h.Repo.FindByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, accounts)
}

// Create Account godoc
// @Summary      Create account
// @Description  Create new sybil account
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body CreateAccountRequest true "Account data"
// @Success      201  {object}  model.Account
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /api/accounts [post]
func (h *AccountHandler) Create(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Color == "" {
		req.Color = "#3B82F6"
	}

	account := &model.Account{
		UserID: userID,
		Name:   req.Name,
		Color:  req.Color,
		Notes:  req.Notes,
	}

	if err := h.Repo.Create(account); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, account)
}

// Get Account godoc
// @Summary      Get account detail
// @Description  Get single account by ID with wallets and account-airdrops
// @Tags         Accounts
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  model.Account
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/accounts/{id} [get]
func (h *AccountHandler) Get(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	account, err := h.Repo.FindByID(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}
	c.JSON(http.StatusOK, account)
}

// Update Account godoc
// @Summary      Update account
// @Description  Update account name, color, notes
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int                  true  "Account ID"
// @Param        body body      UpdateAccountRequest  true  "Updated data"
// @Success      200  {object}  model.Account
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/accounts/{id} [put]
func (h *AccountHandler) Update(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	existing, err := h.Repo.FindByID(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	var req UpdateAccountRequest
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
	if req.Notes != "" {
		existing.Notes = req.Notes
	}

	h.Repo.Update(existing)
	c.JSON(http.StatusOK, existing)
}

// Delete Account godoc
// @Summary      Delete account
// @Description  Delete account (only if empty, or use ?force=true)
// @Tags         Accounts
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int   true  "Account ID"
// @Param        force query     bool  false "Force delete with all data"
// @Success      200   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Router       /api/accounts/{id} [delete]
func (h *AccountHandler) Delete(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	force := c.Query("force") == "true"

	if force {
		if err := h.Repo.DeleteCascade(uint(id), userID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Account and all related data deleted"})
		return
	}

	if err := h.Repo.Delete(uint(id), userID); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Account has wallets or airdrops. Use ?force=true to delete all."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted"})
}

// AssignAirdrop godoc
// @Summary      Assign airdrop to account
// @Description  Assign a global airdrop to an account and sync global tasks
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int                   true  "Account ID"
// @Param        body body      AssignAirdropRequest  true  "Assignment data"
// @Success      201  {object}  model.AccountAirdrop
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/accounts/{id}/airdrops [post]
func (h *AccountHandler) AssignAirdrop(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	accountID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	_, err := h.Repo.FindByID(uint(accountID), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	var req AssignAirdropRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = h.AirdropRepo.FindByID(req.AirdropID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Airdrop not found"})
		return
	}

	// Check if already assigned — prevent duplicates
	var existingAA model.AccountAirdrop
	if err := h.DB.Where("account_id = ? AND airdrop_id = ?", accountID, req.AirdropID).First(&existingAA).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Airdrop already assigned to this account"})
		return
	}

	// Auto-sync tasks from global airdrop tasks
	var globalTasks []model.AirdropTask
	h.DB.Where("airdrop_id = ?", req.AirdropID).Order("sort_order ASC, created_at ASC").Find(&globalTasks)

	var tasks []model.Task
	for _, gt := range globalTasks {
		tasks = append(tasks, model.Task{
			Name:       gt.Name,
			CategoryID: gt.CategoryID,
			Status:     "pending",
			Frequency:  "daily",
			Date:       gt.StartDate,
		})
	}

	aa, err := h.AARepo.AssignAirdrop(uint(accountID), req.AirdropID, "active", req.Notes, tasks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, aa)
}

// GetAccountAirdrops godoc
// @Summary      List airdrops for account
// @Description  Get all airdrop assignments for a specific account
// @Tags         Accounts
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Account ID"
// @Success      200  {array}   model.AccountAirdrop
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/accounts/{id}/airdrops [get]
func (h *AccountHandler) GetAccountAirdrops(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	accountID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	_, err := h.Repo.FindByID(uint(accountID), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	aas, err := h.AARepo.FindByAccount(uint(accountID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, aas)
}

// RemoveAirdrop godoc
// @Summary      Remove airdrop from account
// @Description  Unlink an airdrop from an account and delete associated tasks
// @Tags         Accounts
// @Produce      json
// @Security     BearerAuth
// @Param        id         path      int  true  "Account ID"
// @Param        airdrop_id path      int  true  "Airdrop ID"
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/accounts/{id}/airdrops/{airdrop_id} [delete]
func (h *AccountHandler) RemoveAirdrop(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	accountID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	airdropID, _ := strconv.ParseUint(c.Param("airdrop_id"), 10, 64)

	_, err := h.Repo.FindByID(uint(accountID), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	if err := h.AARepo.RemoveAirdrop(uint(accountID), uint(airdropID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Airdrop removed from account"})
}

// CloneAccount godoc
// @Summary      Clone account
// @Description  Clone an account with its airdrop assignments and tasks (not wallets)
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int                 true  "Source Account ID"
// @Param        body body      CloneAccountRequest true  "New account data"
// @Success      201  {object}  model.Account
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/accounts/{id}/clone [post]
func (h *AccountHandler) CloneAccount(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	sourceID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req CloneAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	if req.Color == "" {
		req.Color = "#3B82F6"
	}

	newAccount, err := h.Repo.CloneAccount(uint(sourceID), userID, req.Name, req.Color)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newAccount)
}

// GetComparison godoc
// @Summary      Get comparison table
// @Description  Get per-account progress stats for comparison
// @Tags         Dashboard
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   handler.ComparisonRow
// @Failure      401  {object}  map[string]string
// @Router       /api/dashboard/comparison [get]
func (h *AccountHandler) GetComparison(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	accounts, err := h.Repo.FindByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var rows []ComparisonRow
	for _, acc := range accounts {
		totalAirdrops := int64(len(acc.AccountAirdrops))
		completedAirdrops := int64(0)
		totalTasks := int64(0)
		completedTasks := int64(0)

		for _, aa := range acc.AccountAirdrops {
			if aa.Status == "completed" {
				completedAirdrops++
			}
			for _, t := range aa.Tasks {
				totalTasks++
				if t.Status == "finish" {
					completedTasks++
				}
			}
		}

		rows = append(rows, ComparisonRow{
			AccountID:         acc.ID,
			AccountName:       acc.Name,
			AccountColor:      acc.Color,
			TotalAirdrops:     totalAirdrops,
			CompletedAirdrops: completedAirdrops,
			TotalTasks:        totalTasks,
			CompletedTasks:    completedTasks,
			PendingTasks:      totalTasks - completedTasks,
			WalletCount:       int64(len(acc.Wallets)),
		})
	}

	c.JSON(http.StatusOK, rows)
}

type ComparisonRow struct {
	AccountID         uint   `json:"account_id"`
	AccountName       string `json:"account_name"`
	AccountColor      string `json:"account_color"`
	TotalAirdrops     int64  `json:"total_airdrops"`
	CompletedAirdrops int64  `json:"completed_airdrops"`
	TotalTasks        int64  `json:"total_tasks"`
	CompletedTasks    int64  `json:"completed_tasks"`
	PendingTasks      int64  `json:"pending_tasks"`
	WalletCount       int64  `json:"wallet_count"`
}
