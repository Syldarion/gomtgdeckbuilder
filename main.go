package main

import (
	"gomtgdeckbuilder/deck"
	"gomtgdeckbuilder/scryfall"
	"gomtgdeckbuilder/ui"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

/*
func CreateImageView(imageURL string) *tview.Image {
	image := tview.NewImage()

	// Download and convert the image
	base64Str, err := imageutils.DownloadImageAsBase64(imageURL)
	if err != nil {
		fmt.Println("Failed to load image:", err)
		return image
	}

	// Decode the Base64 image
	photo, err := imageutils.DecodeBase64Image(base64Str)
	if err != nil {
		fmt.Println("Failed to decode image:", err)
		return image
	}

	image.SetImage(photo)
	return image
}
*/

func NewQueryInput(onSearch func(query string)) *tview.InputField {
	var input tview.InputField
	input = *tview.NewInputField().
		SetLabel("[yellow]Search: ").
		SetFieldWidth(40).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				query := input.GetText()
				onSearch(query)
			}
		})
	return &input
}

func main() {
	app := tview.NewApplication()
	statusBar := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetText("[green]Type a search and press [yellow]Enter[green]. Press [yellow]Q / Esc[green] to quit.")
	detailsView := ui.NewCardDetailsView()
	detailsView.Container.SetBorder(true).SetTitle("Details")
	resultsCardList := ui.NewCardListView()
	resultsCardList.Container.SetBorder(true).SetTitle("Results")

	queryInput := NewQueryInput(func(query string) {
		statusBar.SetText("[blue]Searching...")
		cards, err := scryfall.FetchCards(query)
		if err != nil {
			statusBar.SetText("[red]Error fetching cards.")
			log.Println("Error:", err)
		}

		if len(cards) == 0 {
			statusBar.SetText("[yellow]No cards found.")
		} else {
			resultsCardList.UpdateCards(
				cards,
				func(card scryfall.Card) {
					// onSelectionChanged
					detailsView.Update(card)
				},
				func(card scryfall.Card) {
					// onCardSelected
					app.SetFocus(detailsView.Container)
				})

			// Update details view with first card in the list
			detailsView.Update(resultsCardList.SelectedCard)

			statusBar.SetText("[green]Type a search and press [yellow]Enter[green]. Press [yellow]Esc[green] to quit.")
			app.SetFocus(resultsCardList.Table.DataTable)
		}
	})

	mainLayout := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(resultsCardList.Container, 0, 2, true).
		AddItem(detailsView.Container, 0, 1, false)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(queryInput, 3, 1, true).
		AddItem(mainLayout, 0, 6, false).
		AddItem(statusBar, 1, 1, false)

	// Don't ever set this, we're doing keyboard ONLY
	// app.EnableMouse(true)

	appPages := tview.NewPages()
	page := 0

	appPages.AddPage("Search", layout, true, false)

	// Deck Pages
	deckCollection, err := deck.LoadFromFile("decks.json")
	if err != nil {
		statusBar.SetText("[red]FAILED TO LOAD USER DECKS")
	}

	for _, deck := range deckCollection.Decks {
		for c, card := range deck.Cards {
			newCard := card
			newCard.Card.ParseTypeLine()
			deck.Cards[c] = newCard
		}
	}

	deckDetailsTable := ui.NewDeckDetailsView()
	var activeDeck *deck.Deck

	deckCollectionView := ui.NewDeckCollectionView()
	deckCollectionView.Table.SetSelectedFunc(func(row int, col int) {
		activeDeck = deckCollection.Decks[row]
		deckDetailsTable.SetDeck(activeDeck)
		appPages.SwitchToPage("DeckList")
		app.SetFocus(deckDetailsTable.CardList.DataTable)
	})
	deckCollectionView.SetDeckCollection(deckCollection)

	deckCollectionView.Container.SetTitle("Decks").SetBorder(true)

	deckLayout := tview.NewFlex().
		AddItem(deckCollectionView.Container, 0, 1, true).
		AddItem(tview.NewBox().SetTitle("Details").SetBorder(true), 0, 1, false)

	appPages.AddPage("Decks", deckLayout, true, page == 0)

	deckViewLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(deckDetailsTable.Container, 0, 6, true)

	appPages.AddPage("DeckList", deckViewLayout, true, false)

	// App Setup
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		activePage, _ := appPages.GetFrontPage()

		switch event.Key() {
		case tcell.KeyEscape:
			if activePage == "Search" {
				if app.GetFocus() == detailsView.Container {
					app.SetFocus(resultsCardList.Table.DataTable)
				} else if app.GetFocus() == resultsCardList.Table.DataTable {
					app.SetFocus(queryInput)
				} else {
					deckDetailsTable.SetDeck(activeDeck)
					appPages.SwitchToPage("DeckList")
					app.SetFocus(deckDetailsTable.CardList.DataTable)
				}
			} else if activePage == "DeckList" {
				deckCollectionView.RefreshView()
				appPages.SwitchToPage("Decks")
			} else if activePage == "Decks" {
				app.Stop()
			}
		}

		switch event.Rune() {
		case 's':
			deckCollection.SaveToFile("decks.json")
		case 'a':
			if activePage == "DeckList" {
				appPages.SwitchToPage("Search")
			} else if activePage == "Search" {
				activeDeck.AddCard(resultsCardList.SelectedCard)
			}
		case 'c':
			if activePage == "Decks" {
				activeDeck = deck.NewDeck("New Deck")
				deckCollection.AddDeck(activeDeck)
				appPages.SwitchToPage("DeckList")
			}
		}

		return event
	})

	if err := app.SetRoot(appPages, true).SetFocus(appPages).Run(); err != nil {
		panic(err)
	}
}
