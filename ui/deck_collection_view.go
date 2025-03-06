package ui

import (
	"gomtgdeckbuilder/deck"
	"strconv"

	"github.com/rivo/tview"
)

type DeckCollectionView struct {
	Container  *tview.Flex
	Table      *FixedHeaderTable
	Collection *deck.DeckCollection
}

func NewDeckCollectionView() *DeckCollectionView {
	view := &DeckCollectionView{
		Table: NewFixedHeaderTable(
			[]string{"Name", "Cards"},
			[]int{1, 0},
		),
	}

	view.Container = tview.NewFlex().
		AddItem(view.Table.Container, 0, 1, true)

	return view
}

func (dc *DeckCollectionView) SetDeckCollection(collection *deck.DeckCollection) {
	dc.Collection = collection
	dc.RefreshView()
}

func (dc *DeckCollectionView) RefreshView() {
	tableData := make([][]string, len(dc.Collection.Decks))
	for i, deck := range dc.Collection.Decks {
		tableData[i] = []string{deck.Name, strconv.Itoa(deck.DeckSize())}
	}
	dc.Table.UpdateData(tableData)
}
