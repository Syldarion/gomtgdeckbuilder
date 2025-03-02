package ui

import (
	"fmt"

	"gomtgdeckbuilder/scryfall"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CardDetailsView holds various fields for displaying card details.
type CardDetailsView struct {
	Container         *tview.Flex
	Name              *tview.TextView
	Cost              *tview.TextView
	TypeLine          *tview.TextView
	PowerAndToughness *tview.TextView
	OracleText        *tview.TextView
	AddToDeckButton   *tview.Button
}

// NewCardDetailsView creates a structured card details view.
func NewCardDetailsView(onAdd func()) *CardDetailsView {
	view := &CardDetailsView{
		Name:              tview.NewTextView().SetDynamicColors(true),
		Cost:              tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignRight),
		TypeLine:          tview.NewTextView(),
		PowerAndToughness: tview.NewTextView(),
		OracleText:        tview.NewTextView().SetDynamicColors(true).SetWrap(true),
		AddToDeckButton:   tview.NewButton("(a) Add to Deck").SetSelectedFunc(onAdd),
	}

	headerFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(view.Name, 0, 1, false).
		AddItem(view.Cost, 0, 1, false)

	// Layout with labels
	view.Container = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(headerFlex, 1, 1, false).
		AddItem(view.TypeLine, 1, 1, false).
		AddItem(view.PowerAndToughness, 1, 1, false).
		AddItem(view.OracleText, 0, 5, false).
		AddItem(view.AddToDeckButton, 1, 1, false)

	view.Container.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// a -> add to current deck
		switch event.Rune() {
		case 'a':
			onAdd()
		}

		return event
	})

	return view
}

// Update updates the details section when a card is selected.
func (v *CardDetailsView) Update(card scryfall.Card) {
	v.Name.SetText(fmt.Sprintf("[yellow]%s", card.Name))
	v.Cost.SetText(card.PrettyMana())
	v.TypeLine.SetText(card.TypeLine)

	// Handle Power/Toughness (not all cards have these)
	if card.Power != "" && card.Toughness != "" {
		v.PowerAndToughness.SetText(fmt.Sprintf("%s / %s", card.Power, card.Toughness))
	} else {
		v.PowerAndToughness.SetText("")
	}

	v.OracleText.SetText(card.PrettyOracle())
}
