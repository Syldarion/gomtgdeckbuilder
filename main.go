package main

import (
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
	detailsView := ui.NewCardDetailsView(func() {
		log.Println("ADD")
	})
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

			statusBar.SetText("[green]Type a search and press [yellow]Enter[green]. Press [yellow]Q / Esc[green] to quit.")
			app.SetFocus(resultsCardList.CardTable)
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

	appPages.AddPage("Search", layout, true, page == 0)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			if app.GetFocus() == detailsView.Container {
				app.SetFocus(resultsCardList.CardTable)
			} else if app.GetFocus() == resultsCardList.CardTable {
				app.SetFocus(queryInput)
			} else {
				app.Stop()
			}
		}

		switch event.Rune() {
		case 'q':
			app.Stop()
		}

		return event
	})

	if err := app.SetRoot(appPages, true).SetFocus(queryInput).Run(); err != nil {
		panic(err)
	}
}
