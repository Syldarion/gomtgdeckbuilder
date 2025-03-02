package ui

import (
	"gomtgdeckbuilder/scryfall"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type CardListView struct {
	Container   *tview.Flex
	HeaderTable *tview.Table
	CardTable   *tview.Table
	Headers     []string
}

func NewCardListView() *CardListView {
	view := &CardListView{
		HeaderTable: tview.NewTable().SetSelectable(false, false),
		CardTable:   tview.NewTable().SetSelectable(true, false),
		Headers:     []string{"Name", "Type", "Mana Cost"},
	}

	for col, header := range view.Headers {
		view.HeaderTable.SetCell(0, col, tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetSelectable(false).
			SetAlign(tview.AlignLeft))
	}

	view.Container = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(view.HeaderTable, 1, 1, false).
		AddItem(view.CardTable, 0, 1, false)

	return view
}

func (v *CardListView) UpdateCards(cards []scryfall.Card, onSelectionChanged func(scryfall.Card), onCardSelected func(scryfall.Card)) {
	v.CardTable.Clear()

	for row, card := range cards {
		v.CardTable.SetCell(row, 0, tview.NewTableCell(card.Name).SetAlign(tview.AlignLeft).SetSelectable(true))
		v.CardTable.SetCell(row, 1, tview.NewTableCell(strings.Join(card.SuperTypes, " ")).SetAlign(tview.AlignLeft).SetSelectable(true))
		v.CardTable.SetCell(row, 2, tview.NewTableCell(card.PrettyMana()).SetAlign(tview.AlignLeft).SetSelectable(true))
	}

	for col, header := range v.Headers {
		maxWidth := 0
		for row := 0; row < len(cards); row++ {
			cell := v.CardTable.GetCell(row, col)
			if cell != nil {
				width := len(cell.Text)
				if width > maxWidth {
					maxWidth = width
				}
			}
		}

		paddedHeader := header + strings.Repeat(" ", maxWidth-len(header))

		v.HeaderTable.SetCell(0, col, tview.NewTableCell(paddedHeader).
			SetTextColor(tcell.ColorYellow).
			SetSelectable(false).
			SetAlign(tview.AlignLeft).
			SetMaxWidth(maxWidth))
	}

	v.CardTable.
		SetSelectable(true, false).
		SetSelectionChangedFunc(func(row int, column int) {
			onSelectionChanged(cards[row])
		}).
		SetSelectedFunc(func(row int, column int) {
			onCardSelected(cards[row])
		})
}
