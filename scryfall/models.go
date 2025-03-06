package scryfall

import (
	"fmt"
	"regexp"
	"strings"
)

type Card struct {
	Name       string   `json:"name"`
	Mana       string   `json:"mana_cost"`
	CMC        float32  `json:"cmc"`
	TypeLine   string   `json:"type_line"`
	SuperTypes []string `json:"-"`
	SubTypes   []string `json:"-"`
	Text       string   `json:"oracle_text"`
	Colors     []string `json:"colors"`
	Identity   []string `json:"color_identity"`
	Power      string   `json:"power"`
	Toughness  string   `json:"toughness"`
	ImageURIs  struct {
		ArtCrop string `json:"art_crop"`
	} `json:"image_uris"`
	InCurrentDeck bool              `json:"-"`
	Legalities    map[string]string `json:"legalities"`
}

func (c *Card) ParseTypeLine() {
	parts := strings.Split(c.TypeLine, " — ")
	c.SuperTypes = parseTypes(parts[0])
	if len(parts) > 1 {
		c.SubTypes = parseTypes(parts[1])
	} else {
		c.SubTypes = []string{}
	}
}

func (c *Card) LegalFormats() []string {
	formats := []string{}
	for format, legal := range c.Legalities {
		if legal == "legal" {
			formats = append(formats, format)
		}
	}
	return formats
}

func (c *Card) PrettyMana() string {
	manaCost := c.Mana

	re := regexp.MustCompile(`^\{([\d+X])}`)

	manaCost = re.ReplaceAllStringFunc(manaCost, func(match string) string {
		number := re.FindStringSubmatch(match)[1]
		return number
	})

	manaColorizer := strings.NewReplacer(
		"W", "W[white]",
		"U", "[aqua]U[white]",
		"B", "[gray]B[white]",
		"R", "[red]R[white]",
		"G", "[green]G[white]",
	)

	manaReplacer := strings.NewReplacer(
		"{W}", manaColorizer.Replace("W"),
		"{U}", manaColorizer.Replace("U"),
		"{B}", manaColorizer.Replace("B"),
		"{R}", manaColorizer.Replace("R"),
		"{G}", manaColorizer.Replace("G"),
		"{W/U}", manaColorizer.Replace("[W/U]"),
		"{W/B}", manaColorizer.Replace("[W/B]"),
		"{B/R}", manaColorizer.Replace("[B/R]"),
		"{B/G}", manaColorizer.Replace("[B/G]"),
		"{U/B}", manaColorizer.Replace("[U/B]"),
		"{U/R}", manaColorizer.Replace("[U/R]"),
		"{R/G}", manaColorizer.Replace("[R/G]"),
		"{R/W}", manaColorizer.Replace("[R/W]"),
		"{G/W}", manaColorizer.Replace("[G/W]"),
		"{G/U}", manaColorizer.Replace("[G/U]"),
		"{B/G/P}", manaColorizer.Replace("[B/G/P]"),
		"{B/R/P}", manaColorizer.Replace("[B/R/P]"),
		"{G/U/P}", manaColorizer.Replace("[G/U/P]"),
		"{G/W/P}", manaColorizer.Replace("[G/W/P]"),
		"{R/G/P}", manaColorizer.Replace("[R/G/P]"),
		"{R/W/P}", manaColorizer.Replace("[R/W/P]"),
		"{U/B/P}", manaColorizer.Replace("[U/B/P]"),
		"{U/R/P}", manaColorizer.Replace("[U/R/P]"),
		"{W/B/P}", manaColorizer.Replace("[W/B/P]"),
		"{W/U/P}", manaColorizer.Replace("[W/U/P]"),
		"{C/W}", manaColorizer.Replace("[C/W]"),
		"{C/U}", manaColorizer.Replace("[C/U]"),
		"{C/B}", manaColorizer.Replace("[C/B]"),
		"{C/R}", manaColorizer.Replace("[C/R]"),
		"{C/G}", manaColorizer.Replace("[C/G]"),
		"{2/W}", manaColorizer.Replace("[2/W]"),
		"{2/U}", manaColorizer.Replace("[2/U]"),
		"{2/B}", manaColorizer.Replace("[2/B]"),
		"{2/R}", manaColorizer.Replace("[2/R]"),
		"{2/G}", manaColorizer.Replace("[2/G]"),
		"{W/P}", manaColorizer.Replace("[W/P]"),
		"{U/P}", manaColorizer.Replace("[U/P]"),
		"{B/P}", manaColorizer.Replace("[B/P]"),
		"{R/P}", manaColorizer.Replace("[R/P]"),
		"{G/P}", manaColorizer.Replace("[G/P]"),
		"{C}", "C",
		"{S}", "Snow",
	)

	return manaReplacer.Replace(manaCost)
}

func (c *Card) PrettyOracle() string {
	oracleText := c.Text

	genericManaRegex := regexp.MustCompile(`\{([\d+X])}`)

	oracleText = genericManaRegex.ReplaceAllStringFunc(oracleText, func(match string) string {
		number := genericManaRegex.FindStringSubmatch(match)[1]
		return number
	})

	energyRegex := regexp.MustCompile(`(\{E\})+`)

	oracleText = energyRegex.ReplaceAllStringFunc(oracleText, func(match string) string {
		count := strings.Count(match, "{E}")
		if count == 1 {
			return "Energy Counter(s)"
		}
		return fmt.Sprintf("%d Energy Counters", count)
	})

	pawsRegex := regexp.MustCompile(`(\{P\})+`)

	oracleText = pawsRegex.ReplaceAllStringFunc(oracleText, func(match string) string {
		count := strings.Count(match, "{P}")
		if count == 1 {
			return "Paw(s)"
		}
		return fmt.Sprintf("%d Paws", count)
	})

	ticketRegex := regexp.MustCompile(`(\{TK\})+`)

	oracleText = ticketRegex.ReplaceAllStringFunc(oracleText, func(match string) string {
		count := strings.Count(match, "{TK}")
		if count == 1 {
			return "Ticket Counter(s)"
		}
		return fmt.Sprintf("%d Ticket Counters", count)
	})

	manaColorizer := strings.NewReplacer(
		"W", "W[white]",
		"U", "[aqua]U[white]",
		"B", "[gray]B[white]",
		"R", "[red]R[white]",
		"G", "[green]G[white]",
	)

	symbolReplacer := strings.NewReplacer(
		"{W}", manaColorizer.Replace("W"),
		"{U}", manaColorizer.Replace("U"),
		"{B}", manaColorizer.Replace("B"),
		"{R}", manaColorizer.Replace("R"),
		"{G}", manaColorizer.Replace("G"),
		"{W/U}", manaColorizer.Replace("[W/U]"),
		"{W/B}", manaColorizer.Replace("[W/B]"),
		"{B/R}", manaColorizer.Replace("[B/R]"),
		"{B/G}", manaColorizer.Replace("[B/G]"),
		"{U/B}", manaColorizer.Replace("[U/B]"),
		"{U/R}", manaColorizer.Replace("[U/R]"),
		"{R/G}", manaColorizer.Replace("[R/G]"),
		"{R/W}", manaColorizer.Replace("[R/W]"),
		"{G/W}", manaColorizer.Replace("[G/W]"),
		"{G/U}", manaColorizer.Replace("[G/U]"),
		"{B/G/P}", manaColorizer.Replace("[B/G/P]"),
		"{B/R/P}", manaColorizer.Replace("[B/R/P]"),
		"{G/U/P}", manaColorizer.Replace("[G/U/P]"),
		"{G/W/P}", manaColorizer.Replace("[G/W/P]"),
		"{R/G/P}", manaColorizer.Replace("[R/G/P]"),
		"{R/W/P}", manaColorizer.Replace("[R/W/P]"),
		"{U/B/P}", manaColorizer.Replace("[U/B/P]"),
		"{U/R/P}", manaColorizer.Replace("[U/R/P]"),
		"{W/B/P}", manaColorizer.Replace("[W/B/P]"),
		"{W/U/P}", manaColorizer.Replace("[W/U/P]"),
		"{C/W}", manaColorizer.Replace("[C/W]"),
		"{C/U}", manaColorizer.Replace("[C/U]"),
		"{C/B}", manaColorizer.Replace("[C/B]"),
		"{C/R}", manaColorizer.Replace("[C/R]"),
		"{C/G}", manaColorizer.Replace("[C/G]"),
		"{2/W}", manaColorizer.Replace("[2/W]"),
		"{2/U}", manaColorizer.Replace("[2/U]"),
		"{2/B}", manaColorizer.Replace("[2/B]"),
		"{2/R}", manaColorizer.Replace("[2/R]"),
		"{2/G}", manaColorizer.Replace("[2/G]"),
		"{H}", "Phyrexian mana",
		"{W/P}", manaColorizer.Replace("[W/P]"),
		"{U/P}", manaColorizer.Replace("[U/P]"),
		"{B/P}", manaColorizer.Replace("[B/P]"),
		"{R/P}", manaColorizer.Replace("[R/P]"),
		"{G/P}", manaColorizer.Replace("[G/P]"),
		"{C}", "C",
		"{S}", "Snow",
		"{T}", "Tap",
		"{Q}", "Untap",
		"{CHAOS}", "Chaos",
	)

	return symbolReplacer.Replace(oracleText)
}

func parseTypes(typeString string) []string {
	words := strings.Fields(typeString)
	result := []string{}

	for i := 0; i < len(words); i++ {
		if words[i] == "Time" && i+1 < len(words) && words[i+1] == "Lord" {
			result = append(result, "Time Lord")
			i++
		} else {
			result = append(result, words[i])
		}
	}

	return result
}

type ScryfallResponse struct {
	TotalCards int    `json:"total_cards"`
	HasMore    bool   `json:"has_more"`
	NextPage   string `json:"next_page"`
	Data       []Card `json:"data"`
}
