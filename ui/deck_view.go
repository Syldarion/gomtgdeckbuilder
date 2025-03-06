package ui

import (
	"gomtgdeckbuilder/deck"
	"gomtgdeckbuilder/scryfall"
	"strconv"
	"strings"

	"github.com/rivo/tview"
)

type DeckDetailsView struct {
	Container    *tview.Flex
	Name         *tview.TextView
	CardList     *FixedHeaderTable
	CardDetails  *CardDetailsView
	SelectedCard scryfall.Card
}

func NewDeckDetailsView() *DeckDetailsView {
	view := &DeckDetailsView{
		Name: tview.NewTextView().SetDynamicColors(true),
		CardList: NewFixedHeaderTable(
			[]string{"Name", "Type", "Mana Cost", "Quantity"},
			[]int{1, 0, 0, 0},
		),
		CardDetails: NewCardDetailsView(),
	}

	view.CardList.Container.SetTitle("Decklist").SetBorder(true)
	view.CardDetails.Container.SetTitle("Card Details").SetBorder(true)

	leftColumn := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(view.CardList.Container, 0, 1, true)

	rightColumn := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(view.CardDetails.Container, 0, 1, false)

	deckMainBody := tview.NewFlex().
		AddItem(leftColumn, 0, 1, true).
		AddItem(rightColumn, 0, 1, false)

	view.Container = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(view.Name, 1, 1, false).
		AddItem(deckMainBody, 0, 6, true).
		AddItem(tview.NewTextView().SetDynamicColors(true).SetText("[yellow]A[white] - Add more cards, [yellow]R[white] - Remove selected card, [yellow]Esc[white] - Return to deck list"), 1, 1, false)

	return view
}

func (dv *DeckDetailsView) SetDeck(deck *deck.Deck) {
	dv.Name.SetText(deck.Name)

	newTableData := make([][]string, len(deck.Cards))
	i := 0

	for _, card := range deck.Cards {
		cardData := card.Card
		newTableData[i] = []string{cardData.Name, cardData.TypeLine, cardData.PrettyMana(), strconv.Itoa(card.Quantity)}
		i++
	}

	dv.CardList.UpdateData(newTableData)
	dv.CardList.SetSelectionChangedFunc(func(row int, col int) {
		nameText := strings.Trim(dv.CardList.DataTable.GetCell(row, 0).Text, " ")
		cardAtRow := deck.Cards[nameText].Card
		dv.SelectedCard = cardAtRow
		dv.CardDetails.Update(cardAtRow)
	})
}
