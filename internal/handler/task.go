package handler

import (
	"net/http"
	"strconv"

	"github.com/bsyrlhabibi/airdrop/internal/model"
	"github.com/bsyrlhabibi/airdrop/internal/repository"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	Repo *repository.TaskRepo
}

func NewTaskHandler(repo *repository.TaskRepo) *TaskHandler {
	return &TaskHandler{Repo: repo}
}

type CreateTaskRequest struct {
	Description string  `json:"description" example:"Swap ETH to USDC"`
	Frequency   string  `json:"frequency" example:"weekly"`
	WalletID    *uint   `json:"wallet_id,omitempty"`
}

// List Tasks godoc
// @Summary      List tasks for airdrop
// @Description  Get all tasks for a specific airdrop
// @Tags         Tasks
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Airdrop ID"
// @Success      200  {array}   model.Task
// @Failure      401  {object}  map[string]string
// @Router       /api/airdrops/{id}/tasks [get]
func (h *TaskHandler) List(c *gin.Context) {
	airdropID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	tasks, err := h.Repo.FindByAirdrop(uint(airdropID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// Create Task godoc
// @Summary      Create task
// @Description  Add new task to an airdrop
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int               true  "Airdrop ID"
// @Param        body body      CreateTaskRequest  true  "Task data"
// @Success      201  {object}  model.Task
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /api/airdrops/{id}/tasks [post]
func (h *TaskHandler) Create(c *gin.Context) {
	airdropID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var t model.Task
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t.AirdropID = uint(airdropID)
	h.Repo.Create(&t)
	c.JSON(http.StatusCreated, t)
}

// Complete Task godoc
// @Summary      Mark task as completed
// @Description  Mark a task as done
// @Tags         Tasks
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Task ID"
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/tasks/{id}/complete [put]
func (h *TaskHandler) Complete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.Repo.Complete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Completed"})
}

// Reset Task godoc
// @Summary      Reset task
// @Description  Reset a completed task back to pending
// @Tags         Tasks
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Task ID"
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/tasks/{id}/reset [put]
func (h *TaskHandler) Reset(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.Repo.Reset(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Reset"})
}

// Delete Task godoc
// @Summary      Delete task
// @Description  Remove a task permanently
// @Tags         Tasks
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Task ID"
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /api/tasks/{id} [delete]
func (h *TaskHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.Repo.Delete(uint(id))
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
