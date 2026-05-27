package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bsyrlhabibi/airdrop/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type ExportHandler struct {
	DB *gorm.DB
}

func NewExportHandler(db *gorm.DB) *ExportHandler {
	return &ExportHandler{DB: db}
}

// Style constants
const (
	headerBg    = "1E3A5F"
	headerFg    = "FFFFFF"
	altRowBg    = "F0F4FA"
	borderColor = "D1D5DB"
	greenBg     = "DCFCE7"
	redBg       = "FEE2E2"
	yellowBg    = "FEF3C7"
)

func strPtr(s string) *string { return &s }

// Export godoc
// @Summary      Export data to Excel
// @Description  Export all data as multi-sheet styled Excel file
// @Tags         Export
// @Produce      octet-stream
// @Security     BearerAuth
// @Success      200  {file}  file
// @Failure      401  {object}  map[string]string
// @Router       /api/export/excel [get]
func (h *ExportHandler) ExportExcel(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	f := excelize.NewFile()
	defer f.Close()

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: headerFg, Size: 11, Family: "Calibri"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{headerBg}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "left", Color: headerBg, Style: 1}, {Type: "right", Color: headerBg, Style: 1},
			{Type: "top", Color: headerBg, Style: 1}, {Type: "bottom", Color: headerBg, Style: 1},
		},
	})

	altRowStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{altRowBg}},
		Border: []excelize.Border{
			{Type: "left", Color: borderColor, Style: 1}, {Type: "right", Color: borderColor, Style: 1},
			{Type: "top", Color: borderColor, Style: 1}, {Type: "bottom", Color: borderColor, Style: 1},
		},
	})

	cellStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: borderColor, Style: 1}, {Type: "right", Color: borderColor, Style: 1},
			{Type: "top", Color: borderColor, Style: 1}, {Type: "bottom", Color: borderColor, Style: 1},
		},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})

	greenStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{greenBg}},
		Font: &excelize.Font{Color: "166534"},
		Border: []excelize.Border{
			{Type: "left", Color: borderColor, Style: 1}, {Type: "right", Color: borderColor, Style: 1},
			{Type: "top", Color: borderColor, Style: 1}, {Type: "bottom", Color: borderColor, Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	redStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{redBg}},
		Font: &excelize.Font{Color: "991B1B"},
		Border: []excelize.Border{
			{Type: "left", Color: borderColor, Style: 1}, {Type: "right", Color: borderColor, Style: 1},
			{Type: "top", Color: borderColor, Style: 1}, {Type: "bottom", Color: borderColor, Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	yellowStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{yellowBg}},
		Font: &excelize.Font{Color: "92400E"},
		Border: []excelize.Border{
			{Type: "left", Color: borderColor, Style: 1}, {Type: "right", Color: borderColor, Style: 1},
			{Type: "top", Color: borderColor, Style: 1}, {Type: "bottom", Color: borderColor, Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	centerStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: borderColor, Style: 1}, {Type: "right", Color: borderColor, Style: 1},
			{Type: "top", Color: borderColor, Style: 1}, {Type: "bottom", Color: borderColor, Style: 1},
		},
	})

	// Load data
	var accounts []model.Account
	h.DB.Where("user_id = ?", userID).Order("created_at ASC").Preload("Wallets").Preload("AccountAirdrops").Preload("AccountAirdrops.Airdrop").Preload("AccountAirdrops.Tasks").Preload("AccountAirdrops.Tasks.Category").Find(&accounts)

	var airdrops []model.Airdrop
	h.DB.Where("user_id = ?", userID).Order("name ASC").Find(&airdrops)

	// ========== Sheet 1: Overview ==========
	{
		sheet := "Overview"
		f.SetSheetName("Sheet1", sheet)

		headers := []string{"Account Name", "Total Airdrops", "Completed Airdrops", "Active Airdrops", "Total Tasks", "Finished Tasks", "Pending Tasks", "Completion %", "Total Wallets", "Chains Active", "Status"}
		widths := []float64{22, 15, 18, 15, 14, 17, 16, 15, 15, 16, 14}

		for i, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheet, cell, h)
			f.SetCellStyle(sheet, cell, cell, headerStyle)
			col, _ := excelize.ColumnNumberToName(i + 1)
			f.SetColWidth(sheet, col, col, widths[i])
		}
		f.SetRowHeight(sheet, 1, 30)

		for rowIdx, acc := range accounts {
			row := rowIdx + 2
			totalAirdrops := len(acc.AccountAirdrops)
			completedAirdrops := 0
			activeAirdrops := 0
			totalTasks := 0
			completedTasks := 0
			chains := make(map[string]bool)

			for _, aa := range acc.AccountAirdrops {
				if aa.Status == "completed" {
					completedAirdrops++
				}
				if aa.Status == "active" {
					activeAirdrops++
				}
				totalTasks += len(aa.Tasks)
				for _, t := range aa.Tasks {
					if t.Status == "finish" {
						completedTasks++
					}
				}
			}

			for _, w := range acc.Wallets {
				chains[w.Chain] = true
			}
			chainList := ""
			for ch := range chains {
				if chainList != "" {
					chainList += ", "
				}
				chainList += ch
			}

			pendingTasks := totalTasks - completedTasks
			completionPct := 0.0
			if totalTasks > 0 {
				completionPct = float64(completedTasks) / float64(totalTasks) * 100
			}

			status := "🔴 Not Started"
			if totalTasks > 0 && completionPct == 100 {
				status = "✅ Completed"
			} else if completedTasks > 0 {
				status = "🔄 In Progress"
			}

			vals := []interface{}{acc.Name, totalAirdrops, completedAirdrops, activeAirdrops, totalTasks, completedTasks, pendingTasks, fmt.Sprintf("%.0f%%", completionPct), len(acc.Wallets), chainList, status}

			for colIdx, v := range vals {
				cell, _ := excelize.CoordinatesToCellName(colIdx+1, row)
				f.SetCellValue(sheet, cell, v)
				if colIdx > 0 {
					f.SetCellStyle(sheet, cell, cell, centerStyle)
				} else {
					f.SetCellStyle(sheet, cell, cell, cellStyle)
				}
			}
			if rowIdx%2 == 1 {
				startCell, _ := excelize.CoordinatesToCellName(1, row)
				endCell, _ := excelize.CoordinatesToCellName(len(headers), row)
				f.SetCellStyle(sheet, startCell, endCell, altRowStyle)
			}
			f.SetRowHeight(sheet, row, 24)
		}
	}

	// ========== Sheet 2: Tasks ==========
	{
		sheet := "Tasks"
		f.NewSheet(sheet)

		headers := []string{"Account Name", "Airdrop Name", "Task Name", "Category", "Status", "Date", "Gas Spent", "Tx Hash"}
		widths := []float64{22, 20, 30, 16, 14, 14, 14, 20}

		for i, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheet, cell, h)
			f.SetCellStyle(sheet, cell, cell, headerStyle)
			col, _ := excelize.ColumnNumberToName(i + 1)
			f.SetColWidth(sheet, col, col, widths[i])
		}
		f.SetRowHeight(sheet, 1, 30)

		rowIdx := 0
		for _, acc := range accounts {
			for _, aa := range acc.AccountAirdrops {
				airdropName := ""
				if aa.Airdrop != nil {
					airdropName = aa.Airdrop.Name
				}
				for _, task := range aa.Tasks {
					rowIdx++
					row := rowIdx + 1

					categoryName := ""
					if task.Category != nil {
						categoryName = task.Category.Name
					}

					dateStr := ""
					if task.Date != nil {
						dateStr = task.Date.Format("2006-01-02")
					}

					gasSpent := ""
					if task.GasSpent > 0 {
						gasSpent = fmt.Sprintf("$%.2f", task.GasSpent)
					}

					statusDisplay := task.Status
					switch task.Status {
					case "finish":
						statusDisplay = "Finished"
					case "ongoing":
						statusDisplay = "Ongoing"
					case "missed":
						statusDisplay = "Missed"
					case "pending":
						statusDisplay = "Pending"
					case "cancel":
						statusDisplay = "Cancelled"
					}
					vals := []interface{}{
						acc.Name, airdropName, task.Name, categoryName,
						statusDisplay, dateStr, gasSpent, task.TxHash,
					}

					for colIdx, v := range vals {
						cell, _ := excelize.CoordinatesToCellName(colIdx+1, row)
						f.SetCellValue(sheet, cell, v)
						f.SetCellStyle(sheet, cell, cell, cellStyle)
					}

					// Color code status
					statusCell, _ := excelize.CoordinatesToCellName(5, row)
					switch task.Status {
					case "finish":
						f.SetCellStyle(sheet, statusCell, statusCell, greenStyle)
					case "cancel", "missed":
						f.SetCellStyle(sheet, statusCell, statusCell, redStyle)
					case "ongoing":
						f.SetCellStyle(sheet, statusCell, statusCell, yellowStyle)
					}

					if rowIdx%2 == 0 {
						startCell, _ := excelize.CoordinatesToCellName(1, row)
						endCell, _ := excelize.CoordinatesToCellName(len(headers), row)
						f.SetCellStyle(sheet, startCell, endCell, altRowStyle)
					}
					f.SetRowHeight(sheet, row, 24)
				}
			}
		}
	}

	// ========== Sheet 3: Wallets ==========
	{
		sheet := "Wallets"
		f.NewSheet(sheet)

		headers := []string{"Account Name", "Wallet Label", "Address", "Chain", "Created Date"}
		widths := []float64{22, 20, 48, 16, 16}

		for i, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheet, cell, h)
			f.SetCellStyle(sheet, cell, cell, headerStyle)
			col, _ := excelize.ColumnNumberToName(i + 1)
			f.SetColWidth(sheet, col, col, widths[i])
		}
		f.SetRowHeight(sheet, 1, 30)

		rowIdx := 0
		for _, acc := range accounts {
			for _, w := range acc.Wallets {
				rowIdx++
				row := rowIdx + 1

				vals := []interface{}{acc.Name, w.Label, w.Address, w.Chain, w.CreatedAt.Format("2006-01-02")}
				for colIdx, v := range vals {
					cell, _ := excelize.CoordinatesToCellName(colIdx+1, row)
					f.SetCellValue(sheet, cell, v)
					f.SetCellStyle(sheet, cell, cell, cellStyle)
				}

				if rowIdx%2 == 0 {
					startCell, _ := excelize.CoordinatesToCellName(1, row)
					endCell, _ := excelize.CoordinatesToCellName(len(headers), row)
					f.SetCellStyle(sheet, startCell, endCell, altRowStyle)
				}
				f.SetRowHeight(sheet, row, 24)
			}
		}
	}

	// ========== Sheet 4: Quick Reference ==========
	{
		sheet := "Quick Reference"
		f.NewSheet(sheet)

		sectionStyle, _ := f.NewStyle(&excelize.Style{
			Font:      &excelize.Font{Bold: true, Size: 13, Color: "1E3A5F"},
			Alignment: &excelize.Alignment{Vertical: "center"},
		})

		f.SetCellValue(sheet, "A1", "🪂 All Airdrops")
		f.SetCellStyle(sheet, "A1", "A1", sectionStyle)

		airdropHeaders := []string{"Airdrop Name", "Chain", "Category", "Priority", "Status", "Accounts Assigned"}
		airdropWidths := []float64{25, 16, 14, 12, 12, 20}

		for i, h := range airdropHeaders {
			cell, _ := excelize.CoordinatesToCellName(i+1, 2)
			f.SetCellValue(sheet, cell, h)
			f.SetCellStyle(sheet, cell, cell, headerStyle)
			col, _ := excelize.ColumnNumberToName(i + 1)
			f.SetColWidth(sheet, col, col, airdropWidths[i])
		}

		for i, a := range airdrops {
			row := i + 3
			var accCount int64
			h.DB.Model(&model.AccountAirdrop{}).Where("airdrop_id = ?", a.ID).Count(&accCount)

			vals := []interface{}{a.Name, a.Chain, a.Category, a.Priority, a.Status, accCount}
			for colIdx, v := range vals {
				cell, _ := excelize.CoordinatesToCellName(colIdx+1, row)
				f.SetCellValue(sheet, cell, v)
				f.SetCellStyle(sheet, cell, cell, cellStyle)
			}
			if i%2 == 1 {
				startCell, _ := excelize.CoordinatesToCellName(1, row)
				endCell, _ := excelize.CoordinatesToCellName(len(airdropHeaders), row)
				f.SetCellStyle(sheet, startCell, endCell, altRowStyle)
			}
		}
	}


	// ========== Sheet 5: Airdrop Tasks (Templates) ==========
	{
		sheet := "Airdrop Tasks"
		f.NewSheet(sheet)

		// Load all airdrop tasks with relations
		var allAirdropTasks []model.AirdropTask
		h.DB.Where("airdrop_id IN (?)", h.DB.Model(&model.Airdrop{}).Select("id").Where("user_id = ?", userID)).
			Preload("Airdrop").Preload("Category").Order("airdrop_id ASC, sort_order ASC").Find(&allAirdropTasks)

		headers := []string{"Airdrop Name", "Task Name", "Category", "Status", "Start Date", "End Date"}
		widths := []float64{22, 30, 16, 12, 14, 14}

		for i, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheet, cell, h)
			f.SetCellStyle(sheet, cell, cell, headerStyle)
			col, _ := excelize.ColumnNumberToName(i + 1)
			f.SetColWidth(sheet, col, col, widths[i])
		}
		f.SetRowHeight(sheet, 1, 30)

		for i, t := range allAirdropTasks {
			row := i + 2

			airdropName := ""
			if t.Airdrop != nil {
				airdropName = t.Airdrop.Name
			}
			categoryName := ""
			if t.Category != nil {
				categoryName = t.Category.Name
			}

			startStr := ""
			if t.StartDate != nil {
				startStr = t.StartDate.Format("2006-01-02")
			}
			endStr := ""
			if t.EndDate != nil {
				endStr = t.EndDate.Format("2006-01-02")
			}

			statusDisplay := t.Status
			switch t.Status {
			case "end":
				statusDisplay = "Ended"
			case "ongoing":
				statusDisplay = "Ongoing"
			case "pending":
				statusDisplay = "Pending"
			}

			vals := []interface{}{airdropName, t.Name, categoryName, statusDisplay, startStr, endStr}
			for colIdx, v := range vals {
				cell, _ := excelize.CoordinatesToCellName(colIdx+1, row)
				f.SetCellValue(sheet, cell, v)
				f.SetCellStyle(sheet, cell, cell, cellStyle)
			}

			// Status color
			statusCell, _ := excelize.CoordinatesToCellName(4, row)
			switch t.Status {
			case "end":
				f.SetCellStyle(sheet, statusCell, statusCell, greenStyle)
			case "ongoing":
				f.SetCellStyle(sheet, statusCell, statusCell, yellowStyle)
			}

			if i%2 == 1 {
				startCell, _ := excelize.CoordinatesToCellName(1, row)
				endCell, _ := excelize.CoordinatesToCellName(len(headers), row)
				f.SetCellStyle(sheet, startCell, endCell, altRowStyle)
			}
			f.SetRowHeight(sheet, row, 24)
		}
	}

	f.DeleteSheet("Sheet1")
	f.SetActiveSheet(0)

	// Auto-filter + freeze panes on data sheets
	for _, sheet := range []string{"Overview", "Tasks", "Wallets", "Airdrop Tasks"} {
		rows, _ := f.GetRows(sheet)
		if len(rows) > 1 {
			lastCol, _ := excelize.ColumnNumberToName(len(rows[0]))
			f.AutoFilter(sheet, fmt.Sprintf("A1:%s%d", lastCol, len(rows)), []excelize.AutoFilterOptions{})
		}
		f.SetPanes(sheet, &excelize.Panes{
			Freeze: true, Split: false, XSplit: 0, YSplit: 1,
			TopLeftCell: "A2", ActivePane: "bottomLeft",
		})
	}

	filename := fmt.Sprintf("airdrop-tracker-export-%s.xlsx", time.Now().Format("2006-01-02"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Transfer-Encoding", "binary")

	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate Excel file"})
	}
}
