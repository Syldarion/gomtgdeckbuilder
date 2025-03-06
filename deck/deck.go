package deck

import (
	"gomtgdeckbuilder/scryfall"
)

type Deck struct {
	Name  string                `json:"name"`
	Cards map[string]*DeckEntry `json:"cards"`
}

type DeckEntry struct {
	Card     scryfall.Card `json:"card"`
	Quantity int           `json:"quantity"`
}

func NewDeck(name string) *Deck {
	return &Deck{Name: name, Cards: make(map[string]*DeckEntry)}
}

func (d *Deck) AddCard(card scryfall.Card) {
	if entry, exists := d.Cards[card.Name]; exists {
		entry.Quantity++ // Increase count
	} else {
		d.Cards[card.Name] = &DeckEntry{Card: card, Quantity: 1} // New entry
	}
}

func (d *Deck) RemoveCard(card scryfall.Card) {
	if entry, exists := d.Cards[card.Name]; exists {
		if entry.Quantity > 1 {
			entry.Quantity-- // Decrease count
		} else {
			delete(d.Cards, card.Name) // Remove if last copy
		}
	}
}
