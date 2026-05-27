package handler

import (
	"net/http"
	"strconv"
	"time"

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
	Name      string `json:"name" binding:"required" example:"Bridge 0.1 ETH"`
	CategoryID *uint `json:"category_id"`
	Status    string `json:"status" example:"pending"`
	Frequency string `json:"frequency" example:"once"`
	Date      string `json:"date" example:"2025-01-15"`
}

type UpdateTaskRequest struct {
	Name      string  `json:"name"`
	CategoryID *uint  `json:"category_id"`
	Status    string  `json:"status"`
	Frequency string  `json:"frequency"`
	Date      string  `json:"date"`
	GasSpent  float64 `json:"gas_spent"`
	TxHash    string  `json:"tx_hash"`
}

// List Tasks godoc
// @Summary      List tasks for account-airdrop
// @Description  Get all tasks for a specific account-airdrop
// @Tags         Tasks
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Account Airdrop ID"
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
	freq := req.Frequency
	if freq == "" {
		freq = "once"
	}

	task := &model.Task{
		AccountAirdropID: uint(accountAirdropID),
		Name:             req.Name,
		CategoryID:       req.CategoryID,
		Status:           status,
		Frequency:        freq,
	}

	if req.Date != "" {
		task.Date = parseDate(req.Date)
	} else {
		// Default to today
		now := time.Now()
		task.Date = &now
	}

	if err := h.Repo.Create(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// Update Task godoc
// @Summary      Update task
// @Description  Update task name, category, status, date, frequency
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

	var req UpdateTaskRequest
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
	if req.Frequency != "" {
		task.Frequency = req.Frequency
	}
	if req.Date != "" {
		task.Date = parseDate(req.Date)
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

// Delete Task godoc
// @Summary      Delete task
// @Description  Remove a task permanently
// @Tags         Tasks
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Task ID"
// @Success      200  {object}  map[string]string
// @Router       /api/tasks/{id} [delete]
func (h *TaskHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	h.Repo.Delete(uint(id))
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// TodayTasks godoc
// @Summary      Get today's tasks for account
// @Description  Get all tasks scheduled for today across all airdrops in an account
// @Tags         Tasks
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Account ID"
// @Success      200  {array}   model.Task
// @Router       /api/accounts/{id}/tasks/today [get]
func (h *TaskHandler) TodayTasks(c *gin.Context) {
	accountID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	today := time.Now().Format("2006-01-02")
	tasks, err := h.Repo.FindTodayByAccount(uint(accountID), today)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// DateTasks godoc
// @Summary      Get tasks for a specific date
// @Description  Get all tasks for a specific date in an account
// @Tags         Tasks
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int    true  "Account ID"
// @Param        date  query     string true  "Date (YYYY-MM-DD)"
// @Success      200   {array}   model.Task
// @Router       /api/accounts/{id}/tasks/by-date [get]
func (h *TaskHandler) DateTasks(c *gin.Context) {
	accountID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	date := c.Query("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	tasks, err := h.Repo.FindTodayByAccount(uint(accountID), date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}
