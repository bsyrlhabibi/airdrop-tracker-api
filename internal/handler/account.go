package handler

import (
	"net/http"
	"strconv"

	"github.com/bsyrlhabibi/airdrop/internal/model"
	"github.com/bsyrlhabibi/airdrop/internal/repository"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	Repo *repository.AccountRepo
}

func NewAccountHandler(repo *repository.AccountRepo) *AccountHandler {
	return &AccountHandler{Repo: repo}
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
// @Description  Get single account by ID with wallets and airdrops
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
