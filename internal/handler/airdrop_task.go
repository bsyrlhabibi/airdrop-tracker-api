package handler

import (
	"net/http"
	"strconv"

	"github.com/bsyrlhabibi/airdrop/internal/model"
	"github.com/bsyrlhabibi/airdrop/internal/repository"
	"github.com/gin-gonic/gin"
)

type AirdropTaskHandler struct {
	Repo        *repository.AirdropTaskRepo
	AirdropRepo *repository.AirdropRepo
}

func NewAirdropTaskHandler(repo *repository.AirdropTaskRepo, airdropRepo *repository.AirdropRepo) *AirdropTaskHandler {
	return &AirdropTaskHandler{Repo: repo, AirdropRepo: airdropRepo}
}

type CreateAirdropTaskRequest struct {
	Name      string  `json:"name" binding:"required"`
	CategoryID *uint  `json:"category_id"`
	Status    string  `json:"status"`
	StartDate string  `json:"start_date"`
	EndDate   string  `json:"end_date"`
	GasSpent  float64 `json:"gas_spent"`
	TxHash    string  `json:"tx_hash"`
}

type UpdateAirdropTaskRequest struct {
	Name      string  `json:"name"`
	CategoryID *uint  `json:"category_id"`
	Status    string  `json:"status"`
	StartDate string  `json:"start_date"`
	EndDate   string  `json:"end_date"`
	GasSpent  float64 `json:"gas_spent"`
	TxHash    string  `json:"tx_hash"`
}

// List godoc
// @Summary      List tasks for an airdrop
// @Tags         AirdropTasks
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Airdrop ID"
// @Success      200  {array}  model.AirdropTask
// @Router       /api/airdrops/{airdrop_id}/tasks [get]
func (h *AirdropTaskHandler) List(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	airdropID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	_, err := h.AirdropRepo.FindByID(uint(airdropID), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Airdrop not found"})
		return
	}

	tasks, err := h.Repo.FindByAirdropID(uint(airdropID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// Create godoc
// @Summary      Add task to airdrop
// @Tags         AirdropTasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Airdrop ID"
// @Param        body body CreateAirdropTaskRequest true "Task data"
// @Success      201  {object}  model.AirdropTask
// @Router       /api/airdrops/{airdrop_id}/tasks [post]
func (h *AirdropTaskHandler) Create(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	airdropID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	_, err := h.AirdropRepo.FindByID(uint(airdropID), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Airdrop not found"})
		return
	}

	var req CreateAirdropTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := req.Status
	if status == "" {
		status = "pending"
	}

	task := &model.AirdropTask{
		AirdropID:  uint(airdropID),
		Name:       req.Name,
		CategoryID: req.CategoryID,
		Status:     status,
		StartDate:  parseDate(req.StartDate),
		EndDate:    parseDate(req.EndDate),
		GasSpent:   req.GasSpent,
		TxHash:     req.TxHash,
	}

	if err := h.Repo.Create(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, task)
}

// Update godoc
// @Summary      Update task
// @Tags         AirdropTasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int                     true  "Task ID"
// @Param        body body      CreateAirdropTaskRequest true  "Updated data"
// @Success      200  {object}  model.AirdropTask
// @Router       /api/airdrop-tasks/{id} [put]
func (h *AirdropTaskHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	task, err := h.Repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var req UpdateAirdropTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		task.Name = req.Name
	}
	if req.CategoryID != nil {
		task.CategoryID = req.CategoryID
	}
	if req.Status != "" {
		task.Status = req.Status
	}
	if req.StartDate != "" {
		task.StartDate = parseDate(req.StartDate)
	}
	if req.EndDate != "" {
		task.EndDate = parseDate(req.EndDate)
	}
	if req.GasSpent > 0 {
		task.GasSpent = req.GasSpent
	}
	if req.TxHash != "" {
		task.TxHash = req.TxHash
	}

	if err := h.Repo.Update(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

// Delete godoc
// @Summary      Delete task
// @Tags         AirdropTasks
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Task ID"
// @Success      200  {object}  map[string]string
// @Router       /api/airdrop-tasks/{id} [delete]
func (h *AirdropTaskHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.Repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
