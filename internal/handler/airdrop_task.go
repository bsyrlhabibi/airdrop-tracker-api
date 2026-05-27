package handler

import (
	"net/http"
	"strconv"
	"time"

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
	Description string `json:"description" binding:"required"`
	Frequency   string `json:"frequency"`
}

// List godoc
// @Summary      List tasks for an airdrop
// @Description  Get all tasks for a specific airdrop
// @Tags         AirdropTasks
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Airdrop ID"
// @Success      200  {array}  model.AirdropTask
// @Router       /api/airdrops/{airdrop_id}/tasks [get]
func (h *AirdropTaskHandler) List(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	airdropID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	// Verify airdrop belongs to user
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
// @Description  Create a new task for an airdrop
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

	// Verify airdrop belongs to user
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

	freq := req.Frequency
	if freq == "" {
		freq = "once"
	}

	task := &model.AirdropTask{
		AirdropID:   uint(airdropID),
		Description: req.Description,
		Frequency:   freq,
	}

	if err := h.Repo.Create(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, task)
}

// Complete godoc
// @Summary      Toggle task completion
// @Description  Mark task as completed or uncompleted
// @Tags         AirdropTasks
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Task ID"
// @Success      200  {object}  model.AirdropTask
// @Router       /api/airdrop-tasks/{id}/complete [put]
func (h *AirdropTaskHandler) Complete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	task, err := h.Repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	newStatus := !task.IsCompleted
	if err := h.Repo.ToggleComplete(uint(id), newStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Refresh
	task, _ = h.Repo.FindByID(uint(id))
	c.JSON(http.StatusOK, task)
}

// Delete godoc
// @Summary      Delete task
// @Description  Remove a task from an airdrop
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

// Reorder godoc
// @Summary      Reorder tasks
// @Description  Update sort order of tasks
// @Tags         AirdropTasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Airdrop ID"
// @Param        body body []uint true "Task IDs in order"
// @Success      200  {object}  map[string]string
// @Router       /api/airdrops/{airdrop_id}/tasks/reorder [put]
func (h *AirdropTaskHandler) Reorder(c *gin.Context) {
	airdropID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var taskIDs []uint
	if err := c.ShouldBindJSON(&taskIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, taskID := range taskIDs {
		h.Repo.DB.Model(&model.AirdropTask{}).
			Where("id = ? AND airdrop_id = ?", taskID, uint(airdropID)).
			Update("sort_order", i)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reordered"})
}

// BulkCreate godoc
// @Summary      Bulk add tasks
// @Description  Create multiple tasks at once
// @Tags         AirdropTasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Airdrop ID"
// @Param        body body []CreateAirdropTaskRequest true "Tasks"
// @Success      201  {array}  model.AirdropTask
// @Router       /api/airdrops/{airdrop_id}/tasks/bulk [post]
func (h *AirdropTaskHandler) BulkCreate(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	airdropID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	_, err := h.AirdropRepo.FindByID(uint(airdropID), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Airdrop not found"})
		return
	}

	var reqs []CreateAirdropTaskRequest
	if err := c.ShouldBindJSON(&reqs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var created []model.AirdropTask
	for _, req := range reqs {
		freq := req.Frequency
		if freq == "" {
			freq = "once"
		}
		task := &model.AirdropTask{
			AirdropID:   uint(airdropID),
			Description: req.Description,
			Frequency:   freq,
		}
		if err := h.Repo.Create(task); err == nil {
			created = append(created, *task)
		}
	}

	c.JSON(http.StatusCreated, created)
}

// ResetAll godoc
// @Summary      Reset all tasks
// @Description  Mark all tasks as uncompleted
// @Tags         AirdropTasks
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Airdrop ID"
// @Success      200  {object}  map[string]string
// @Router       /api/airdrops/{airdrop_id}/tasks/reset [put]
func (h *AirdropTaskHandler) ResetAll(c *gin.Context) {
	airdropID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	h.Repo.DB.Model(&model.AirdropTask{}).
		Where("airdrop_id = ?", uint(airdropID)).
		Updates(map[string]interface{}{
			"is_completed": false,
			"completed_at": nil,
		})

	c.JSON(http.StatusOK, gin.H{"message": "All tasks reset"})
}

// unused import guard
var _ = time.Now
