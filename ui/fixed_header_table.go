package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FixedHeaderTable struct {
	Container   *tview.Flex
	HeaderTable *tview.Table
	DataTable   *tview.Table
	Headers     []string
	Expansions  []int
}

func NewFixedHeaderTable(headers []string, expansions []int) *FixedHeaderTable {
	if len(headers) != len(expansions) {
		panic("Headers and Expansions must have the same length!")
	}

	view := &FixedHeaderTable{
		HeaderTable: tview.NewTable().SetSelectable(false, false),
		DataTable:   tview.NewTable().SetSelectable(true, false),
		Headers:     headers,
		Expansions:  expansions,
	}

	for col, header := range headers {
		view.HeaderTable.SetCell(0, col, tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetSelectable(false).
			SetAlign(tview.AlignLeft).
			SetExpansion(expansions[col]))
	}

	view.Container = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(view.HeaderTable, 1, 1, false). // Fixed header
		AddItem(view.DataTable, 0, 1, true)     // Scrollable table

	return view
}

func (v *FixedHeaderTable) UpdateData(data [][]string) {
	v.DataTable.Clear()

	// Ensure headers are as wide as their widest column
	colWidths := v.SyncHeaders(data)

	for row, rowData := range data {
		for col, value := range rowData {
			// Ensure text is at least as wide as the header
			paddedValue, _ := padTaggedString(value, colWidths[col])

			v.DataTable.SetCell(row, col, tview.NewTableCell(paddedValue).
				SetAlign(tview.AlignLeft).
				SetSelectable(true).
				SetExpansion(v.Expansions[col])) // Apply expansion settings
		}
	}
}

func padTaggedString(text string, minWidth int) (string, int) {
	textLen := tview.TaggedStringWidth(text)
	if textLen < minWidth {
		padded := text + strings.Repeat(" ", minWidth-textLen)
		return padded, minWidth
	}
	return text, textLen
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func (v *FixedHeaderTable) SyncHeaders(data [][]string) []int {
	maxColWidths := make([]int, len(v.Headers))

	// Start with header widths
	for col, header := range v.Headers {
		maxColWidths[col] = len(header) // Header defines minimum width
	}

	// Find the widest cell in each column
	for row := 0; row < len(data); row++ {
		for col := 0; col < len(data[row]); col++ {
			cellWidth := tview.TaggedStringWidth(data[row][col])
			maxColWidths[col] = max(maxColWidths[col], cellWidth)
		}
	}

	// Update headers with correct width
	for col, header := range v.Headers {
		paddedHeader, _ := padTaggedString(header, maxColWidths[col])

		v.HeaderTable.SetCell(0, col, tview.NewTableCell(paddedHeader).
			SetTextColor(tcell.ColorYellow).
			SetSelectable(false).
			SetAlign(tview.AlignLeft).
			SetMaxWidth(maxColWidths[col]).
			SetExpansion(v.Expansions[col]))
	}

	return maxColWidths
}

func (v *FixedHeaderTable) SetSelectedFunc(f func(row int, col int)) {
	v.DataTable.SetSelectedFunc(f)
}

func (v *FixedHeaderTable) SetSelectionChangedFunc(f func(row int, col int)) {
	v.DataTable.SetSelectionChangedFunc(f)
}
