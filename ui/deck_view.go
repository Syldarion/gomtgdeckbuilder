package ui

import "github.com/rivo/tview"

type DeckDetailsView struct {
	Container *tview.Flex
	DeckName  *tview.TextView
}

func NewDeckDetailsView() *DeckDetailsView {
	view := &DeckDetailsView{
		DeckName: tview.NewTextView().SetDynamicColors(true),
	}

	view.Container = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(view.DeckName, 1, 1, false)

	return view
}
