package deck

import (
	"gomtgdeckbuilder/scryfall"
	"slices"
)

type Deck struct {
	Name   string                `json:"name"`
	Format string                `json:"format"`
	Cards  map[string]*DeckEntry `json:"cards"`
}

type DeckEntry struct {
	Card     scryfall.Card `json:"card"`
	Quantity int           `json:"quantity"`
}

func NewDeck(name string, format string) *Deck {
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

func (d *Deck) DeckSize() int {
	count := 0
	for _, card := range d.Cards {
		count += card.Quantity
	}
	return count
}

func (d *Deck) ValidateDeck() bool {
	switch d.Format {
	case "standard":
		return d.ValidateStandardDeck()
	case "commander":
		return d.ValidateCommanderDeck()
	}

	return false
}

func (d *Deck) ValidateStandardDeck() bool {
	if d.DeckSize() < 60 {
		return false
	}

	for _, card := range d.Cards {
		if !slices.Contains(card.Card.LegalFormats(), "standard") {
			return false
		}

		if !slices.Contains(card.Card.SuperTypes, "Basic") && card.Quantity > 4 {
			return false
		}
	}

	return true
}

func (d *Deck) ValidateCommanderDeck() bool {
	if d.DeckSize() < 100 {
		return false
	}

	for _, card := range d.Cards {
		if !slices.Contains(card.Card.LegalFormats(), "commander") {
			return false
		}

		if card.Quantity > 1 {
			return false
		}
	}

	return true
}
