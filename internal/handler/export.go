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
	headerBg    = "1E3A5F" // Dark blue
	headerFg    = "FFFFFF" // White text
	altRowBg    = "F0F4FA" // Light blue-gray
	borderColor = "D1D5DB" // Gray border
	greenBg     = "DCFCE7" // Light green for completed
	redBg       = "FEE2E2" // Light red for missed
	yellowBg    = "FEF3C7" // Light yellow for in progress
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

	// Create styles
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: headerFg, Size: 11, Family: "Calibri"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{headerBg}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "left", Color: headerBg, Style: 1},
			{Type: "right", Color: headerBg, Style: 1},
			{Type: "top", Color: headerBg, Style: 1},
			{Type: "bottom", Color: headerBg, Style: 1},
		},
	})

	altRowStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{altRowBg}},
		Border: []excelize.Border{
			{Type: "left", Color: borderColor, Style: 1},
			{Type: "right", Color: borderColor, Style: 1},
			{Type: "top", Color: borderColor, Style: 1},
			{Type: "bottom", Color: borderColor, Style: 1},
		},
	})

	cellStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: borderColor, Style: 1},
			{Type: "right", Color: borderColor, Style: 1},
			{Type: "top", Color: borderColor, Style: 1},
			{Type: "bottom", Color: borderColor, Style: 1},
		},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})

	greenStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{greenBg}},
		Font: &excelize.Font{Color: "166534"},
		Border: []excelize.Border{
			{Type: "left", Color: borderColor, Style: 1},
			{Type: "right", Color: borderColor, Style: 1},
			{Type: "top", Color: borderColor, Style: 1},
			{Type: "bottom", Color: borderColor, Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	redStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{redBg}},
		Font: &excelize.Font{Color: "991B1B"},
		Border: []excelize.Border{
			{Type: "left", Color: borderColor, Style: 1},
			{Type: "right", Color: borderColor, Style: 1},
			{Type: "top", Color: borderColor, Style: 1},
			{Type: "bottom", Color: borderColor, Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	yellowStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{yellowBg}},
		Font: &excelize.Font{Color: "92400E"},
		Border: []excelize.Border{
			{Type: "left", Color: borderColor, Style: 1},
			{Type: "right", Color: borderColor, Style: 1},
			{Type: "top", Color: borderColor, Style: 1},
			{Type: "bottom", Color: borderColor, Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	centerStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: borderColor, Style: 1},
			{Type: "right", Color: borderColor, Style: 1},
			{Type: "top", Color: borderColor, Style: 1},
			{Type: "bottom", Color: borderColor, Style: 1},
		},
	})

	// Load all data
	var accounts []model.Account
	h.DB.Where("user_id = ?", userID).Order("created_at ASC").Preload("Wallets").Preload("AccountAirdrops").Preload("AccountAirdrops.Airdrop").Preload("AccountAirdrops.Tasks").Find(&accounts)

	var airdrops []model.Airdrop
	h.DB.Where("user_id = ?", userID).Order("name ASC").Find(&airdrops)

	// ========== Sheet 1: Overview ==========
	{
		sheet := "Overview"
		f.SetSheetName("Sheet1", sheet)
		f.SetSheetProps(sheet, &excelize.SheetPropsOptions{TabColorRGB: strPtr("1E3A5F")})

		headers := []string{"Account Name", "Total Airdrops", "Completed Airdrops", "Active Airdrops", "Total Tasks", "Completed Tasks", "Pending Tasks", "Completion %", "Total Wallets", "Chains Active", "Status"}
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
					if t.IsCompleted {
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
				status = "🟢 Completed"
			} else if completedTasks > 0 {
				status = "🟡 In Progress"
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
			// Alternate row colors
			if rowIdx%2 == 1 {
				startCell, _ := excelize.CoordinatesToCellName(1, row)
				endCell, _ := excelize.CoordinatesToCellName(len(headers), row)
				f.SetCellStyle(sheet, startCell, endCell, altRowStyle)
			}
			f.SetRowHeight(sheet, row, 24)
		}
	}

	// ========== Sheet 2: Airdrop Detail ==========
	{
		sheet := "Airdrop Detail"
		f.NewSheet(sheet)
		f.SetSheetProps(sheet, &excelize.SheetPropsOptions{TabColorRGB: strPtr("7C3AED")})

		headers := []string{"Account Name", "Airdrop Name", "Chain", "Category", "Status", "Total Tasks", "Completed", "Completion %", "Assigned Date", "Notes"}
		widths := []float64{22, 20, 14, 14, 16, 14, 14, 15, 16, 30}

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
				rowIdx++
				row := rowIdx + 1

				totalTasks := len(aa.Tasks)
				completed := 0
				for _, t := range aa.Tasks {
					if t.IsCompleted {
						completed++
					}
				}

				pct := 0.0
				if totalTasks > 0 {
					pct = float64(completed) / float64(totalTasks) * 100
				}

				statusStr := aa.Status
				if aa.Airdrop != nil {
					statusStr = aa.Airdrop.Status
				}

				airdropName := ""
				chain := ""
				category := ""
				notes := ""
				if aa.Airdrop != nil {
					airdropName = aa.Airdrop.Name
					chain = aa.Airdrop.Chain
					category = aa.Airdrop.Category
					notes = aa.Notes
				}

				vals := []interface{}{
					acc.Name,
					airdropName,
					chain,
					category,
					statusStr,
					totalTasks,
					completed,
					fmt.Sprintf("%.0f%%", pct),
					aa.CreatedAt.Format("2006-01-02"),
					notes,
				}

				for colIdx, v := range vals {
					cell, _ := excelize.CoordinatesToCellName(colIdx+1, row)
					f.SetCellValue(sheet, cell, v)
					f.SetCellStyle(sheet, cell, cell, cellStyle)
				}

				// Color code status
				statusCell, _ := excelize.CoordinatesToCellName(5, row)
				switch statusStr {
				case "completed":
					f.SetCellStyle(sheet, statusCell, statusCell, greenStyle)
				case "missed", "dropped":
					f.SetCellStyle(sheet, statusCell, statusCell, redStyle)
				default:
					f.SetCellStyle(sheet, statusCell, statusCell, yellowStyle)
				}

				// Center numeric cols
				for _, c := range []int{3, 4, 6, 7, 8, 9} {
					cell, _ := excelize.CoordinatesToCellName(c, row)
					f.SetCellStyle(sheet, cell, cell, centerStyle)
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

	// ========== Sheet 3: Tasks ==========
	{
		sheet := "Tasks"
		f.NewSheet(sheet)
		f.SetSheetProps(sheet, &excelize.SheetPropsOptions{TabColorRGB: strPtr("059669")})

		headers := []string{"Account Name", "Airdrop Name", "Task Description", "Frequency", "Completed", "Completed At", "Gas Spent", "Tx Hash", "Notes"}
		widths := []float64{22, 20, 30, 14, 14, 16, 14, 20, 30}

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

					completedStr := "No"
					completedAt := ""
					if task.IsCompleted {
						completedStr = "Yes"
						if task.CompletedAt != nil {
							completedAt = task.CompletedAt.Format("2006-01-02 15:04")
						}
					}

					gasSpent := ""
					if task.GasSpent > 0 {
						gasSpent = fmt.Sprintf("$%.2f", task.GasSpent)
					}

					vals := []interface{}{
						acc.Name,
						airdropName,
						task.Description,
						task.Frequency,
						completedStr,
						completedAt,
						gasSpent,
						task.TxHash,
						"",
					}

					for colIdx, v := range vals {
						cell, _ := excelize.CoordinatesToCellName(colIdx+1, row)
						f.SetCellValue(sheet, cell, v)
						f.SetCellStyle(sheet, cell, cell, cellStyle)
					}

					// Color code completed
					compCell, _ := excelize.CoordinatesToCellName(5, row)
					if task.IsCompleted {
						f.SetCellStyle(sheet, compCell, compCell, greenStyle)
					} else {
						f.SetCellStyle(sheet, compCell, compCell, redStyle)
					}

					// Center cols
					for _, c := range []int{4, 5, 7} {
						cell, _ := excelize.CoordinatesToCellName(c, row)
						f.SetCellStyle(sheet, cell, cell, centerStyle)
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

	// ========== Sheet 4: Wallets ==========
	{
		sheet := "Wallets"
		f.NewSheet(sheet)
		f.SetSheetProps(sheet, &excelize.SheetPropsOptions{TabColorRGB: strPtr("D97706")})

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

				vals := []interface{}{
					acc.Name,
					w.Label,
					w.Address,
					w.Chain,
					w.CreatedAt.Format("2006-01-02"),
				}

				for colIdx, v := range vals {
					cell, _ := excelize.CoordinatesToCellName(colIdx+1, row)
					f.SetCellValue(sheet, cell, v)
					f.SetCellStyle(sheet, cell, cell, cellStyle)
				}

				// Center chain & date
				for _, c := range []int{4, 5} {
					cell, _ := excelize.CoordinatesToCellName(c, row)
					f.SetCellStyle(sheet, cell, cell, centerStyle)
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

	// ========== Sheet 5: Quick Reference ==========
	{
		sheet := "Quick Reference"
		f.NewSheet(sheet)
		f.SetSheetProps(sheet, &excelize.SheetPropsOptions{TabColorRGB: strPtr("DC2626")})

		// Section 1: All Airdrops
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
			// Count how many accounts have this airdrop
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

		// Section 2: Chains Summary
		chainRow := len(airdrops) + 5
		f.SetCellValue(sheet, fmt.Sprintf("A%d", chainRow), "⛓️ All Chains")
		f.SetCellStyle(sheet, fmt.Sprintf("A%d", chainRow), fmt.Sprintf("A%d", chainRow), sectionStyle)

		chainHeaders := []string{"Chain", "Total Wallets", "Accounts Using"}
		for i, h := range chainHeaders {
			cell, _ := excelize.CoordinatesToCellName(i+1, chainRow+1)
			f.SetCellValue(sheet, cell, h)
			f.SetCellStyle(sheet, cell, cell, headerStyle)
		}

		chainStats := make(map[string]struct {
			Wallets  int
			Accounts map[uint]bool
		})
		for _, acc := range accounts {
			for _, w := range acc.Wallets {
				cs := chainStats[w.Chain]
				cs.Wallets++
				if cs.Accounts == nil {
					cs.Accounts = make(map[uint]bool)
				}
				cs.Accounts[acc.ID] = true
				chainStats[w.Chain] = cs
			}
		}

		chainIdx := 0
		for chain, cs := range chainStats {
			r := chainRow + 2 + chainIdx
			f.SetCellValue(sheet, fmt.Sprintf("A%d", r), chain)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", r), cs.Wallets)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", r), len(cs.Accounts))
			for _, c := range []string{"A", "B", "C"} {
				cell := fmt.Sprintf("%s%d", c, r)
				f.SetCellStyle(sheet, cell, cell, cellStyle)
			}
			chainIdx++
		}

		// Section 3: Account-Wallet Matrix
		matrixRow := chainRow + 2 + chainIdx + 3
		f.SetCellValue(sheet, fmt.Sprintf("A%d", matrixRow), "🔗 Account → Wallet Matrix")
		f.SetCellStyle(sheet, fmt.Sprintf("A%d", matrixRow), fmt.Sprintf("A%d", matrixRow), sectionStyle)

		matrixHeaders := []string{"Account", "Wallet Address", "Chain", "Label"}
		for i, h := range matrixHeaders {
			cell, _ := excelize.CoordinatesToCellName(i+1, matrixRow+1)
			f.SetCellValue(sheet, cell, h)
			f.SetCellStyle(sheet, cell, cell, headerStyle)
		}

		mIdx := 0
		for _, acc := range accounts {
			for _, w := range acc.Wallets {
				r := matrixRow + 2 + mIdx
				f.SetCellValue(sheet, fmt.Sprintf("A%d", r), acc.Name)
				f.SetCellValue(sheet, fmt.Sprintf("B%d", r), w.Address)
				f.SetCellValue(sheet, fmt.Sprintf("C%d", r), w.Chain)
				f.SetCellValue(sheet, fmt.Sprintf("D%d", r), w.Label)
				for _, c := range []string{"A", "B", "C", "D"} {
					cell := fmt.Sprintf("%s%d", c, r)
					f.SetCellStyle(sheet, cell, cell, cellStyle)
				}
				mIdx++
			}
		}
	}

	// Delete default Sheet1 if renamed
	f.DeleteSheet("Sheet1")

	// Set first sheet as active
	f.SetActiveSheet(0)

	// Freeze panes (header row) for each sheet
	for _, sheet := range []string{"Overview", "Airdrop Detail", "Tasks", "Wallets"} {
		f.SetPanes(sheet, &excelize.Panes{
			Freeze:      true,
			Split:       false,
			XSplit:      0,
			YSplit:      1,
			TopLeftCell: "A2",
			ActivePane:  "bottomLeft",
		})
	}

	// Generate file
	filename := fmt.Sprintf("airdrop-tracker-export-%s.xlsx", time.Now().Format("2006-01-02"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Transfer-Encoding", "binary")

	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate Excel file"})
		return
	}
}
