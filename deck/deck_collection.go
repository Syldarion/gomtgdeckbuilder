package deck

import (
	"encoding/json"
	"os"
)

type DeckCollection struct {
	Decks []*Deck `json:"decks"`
}

func NewDeckCollection() *DeckCollection {
	return &DeckCollection{Decks: []*Deck{}}
}

func (dc *DeckCollection) AddDeck(Deck *Deck) {
	dc.Decks = append(dc.Decks, Deck)
}

func (dc *DeckCollection) GetDeck(name string) *Deck {
	for i, deck := range dc.Decks {
		if deck.Name == name {
			return dc.Decks[i]
		}
	}

	return nil
}

func (dc *DeckCollection) SaveToFile(filename string) error {
	data, err := json.MarshalIndent(dc, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func LoadFromFile(filename string) (*DeckCollection, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return NewDeckCollection(), nil
		}

		return nil, err
	}

	var collection DeckCollection
	err = json.Unmarshal(data, &collection)

	if err != nil {
		return nil, err
	}

	return &collection, nil
}
