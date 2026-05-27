package handler

import (
	"net/http"
	"strconv"

	"github.com/bsyrlhabibi/airdrop/internal/model"
	"github.com/bsyrlhabibi/airdrop/internal/repository"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	Repo   *repository.TaskRepo
	AARepo *repository.AccountAirdropRepo
}

func NewTaskHandler(repo *repository.TaskRepo, aaRepo *repository.AccountAirdropRepo) *TaskHandler {
	return &TaskHandler{Repo: repo, AARepo: aaRepo}
}

type CreateTaskRequest struct {
	Name       string `json:"name" binding:"required" example:"Bridge 0.1 ETH"`
	CategoryID *uint  `json:"category_id"`
	Status     string `json:"status" example:"pending"`
	Date       string `json:"date" example:"2025-01-15"`
}

// List Tasks godoc
// @Summary      List tasks for account-airdrop
// @Description  Get all tasks for a specific account-airdrop
// @Tags         Tasks
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Account Airdrop ID"
// @Success      200  {array}   model.Task
// @Failure      401  {object}  map[string]string
// @Router       /api/account-airdrops/{id}/tasks [get]
func (h *TaskHandler) List(c *gin.Context) {
	accountAirdropID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	tasks, err := h.Repo.FindByAccountAirdrop(uint(accountAirdropID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// Create Task godoc
// @Summary      Create task
// @Description  Add new task to an account-airdrop
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int               true  "Account Airdrop ID"
// @Param        body body      CreateTaskRequest  true  "Task data"
// @Success      201  {object}  model.Task
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /api/account-airdrops/{id}/tasks [post]
func (h *TaskHandler) Create(c *gin.Context) {
	accountAirdropID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := req.Status
	if status == "" {
		status = "pending"
	}

	task := &model.Task{
		AccountAirdropID: uint(accountAirdropID),
		Name:             req.Name,
		CategoryID:       req.CategoryID,
		Status:           status,
	}

	if req.Date != "" {
		task.Date = parseDate(req.Date)
	}

	if err := h.Repo.Create(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// Update Task godoc
// @Summary      Update task
// @Description  Update task name, category, status, date
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int               true  "Task ID"
// @Param        body body      CreateTaskRequest  true  "Updated data"
// @Success      200  {object}  model.Task
// @Router       /api/tasks/{id} [put]
func (h *TaskHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	task, err := h.Repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var req CreateTaskRequest
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
	if req.Date != "" {
		task.Date = parseDate(req.Date)
	}

	if err := h.Repo.Update(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
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
