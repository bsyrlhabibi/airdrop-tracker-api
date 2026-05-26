package handler

import (
	"net/http"
	"strconv"

	"github.com/bsyrlhabibi/airdrop/internal/model"
	"github.com/bsyrlhabibi/airdrop/internal/repository"
	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	Repo *repository.WalletRepo
}

func NewWalletHandler(repo *repository.WalletRepo) *WalletHandler {
	return &WalletHandler{Repo: repo}
}

type CreateWalletRequest struct {
	AccountID uint   `json:"account_id" example:"1"`
	Label     string `json:"label" example:"Main Wallet"`
	Address   string `json:"address" example:"0x7245...139b"`
	Chain     string `json:"chain" example:"Ethereum"`
}

// List Wallets godoc
// @Summary      List all wallets
// @Description  Get all wallets, optionally filtered by account_id
// @Tags         Wallets
// @Produce      json
// @Security     BearerAuth
// @Param        account_id query int false "Filter by Account ID"
// @Success      200  {array}   model.Wallet
// @Failure      401  {object}  map[string]string
// @Router       /api/wallets [get]
func (h *WalletHandler) List(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	accountIDStr := c.Query("account_id")

	if accountIDStr != "" {
		accountID, _ := strconv.ParseUint(accountIDStr, 10, 64)
		wallets, err := h.Repo.FindByAccount(uint(accountID), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, wallets)
		return
	}

	wallets, _ := h.Repo.FindByUser(userID)
	c.JSON(http.StatusOK, wallets)
}

// Create Wallet godoc
// @Summary      Add wallet
// @Description  Register a new wallet address to an account
// @Tags         Wallets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body CreateWalletRequest true "Wallet data"
// @Success      201  {object}  model.Wallet
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /api/wallets [post]
func (h *WalletHandler) Create(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var req CreateWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet := &model.Wallet{
		UserID:    userID,
		AccountID: req.AccountID,
		Label:     req.Label,
		Address:   req.Address,
		Chain:     req.Chain,
	}

	if err := h.Repo.Create(wallet); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, wallet)
}

// Delete Wallet godoc
// @Summary      Remove wallet
// @Description  Delete a wallet address
// @Tags         Wallets
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Wallet ID"
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /api/wallets/{id} [delete]
func (h *WalletHandler) Delete(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.Repo.Delete(uint(id), userID)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
