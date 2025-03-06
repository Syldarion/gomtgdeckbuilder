package ui

import (
	"gomtgdeckbuilder/scryfall"
	"strings"

	"github.com/rivo/tview"
)

type CardListView struct {
	Container    *tview.Flex
	Table        *FixedHeaderTable
	SelectedCard scryfall.Card
}

func NewCardListView() *CardListView {
	view := &CardListView{
		Table: NewFixedHeaderTable(
			[]string{"Name", "Type", "Mana Cost"},
			[]int{1, 0, 0},
		),
	}

	view.Container = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(view.Table.Container, 0, 1, false)

	return view
}

func (v *CardListView) UpdateCards(cards []scryfall.Card, onSelectionChanged func(scryfall.Card), onCardSelected func(scryfall.Card)) {
	newTableData := make([][]string, len(cards))
	for i, card := range cards {
		newTableData[i] = []string{card.Name, strings.Join(card.SuperTypes, " "), card.PrettyMana()}
	}

	v.Table.UpdateData(newTableData)
	v.Table.SetSelectionChangedFunc(func(row int, column int) {
		v.SelectedCard = cards[row]
		onSelectionChanged(cards[row])
	})
	v.Table.SetSelectedFunc(func(row int, column int) {
		v.SelectedCard = cards[row]
		onCardSelected(cards[row])
	})

	v.Table.DataTable.Select(0, 0)
	v.Table.DataTable.ScrollToBeginning()
}
