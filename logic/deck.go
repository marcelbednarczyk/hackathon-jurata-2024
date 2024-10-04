package logic

import (
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"os"
	"reflect"
)

var Deck = loadDeck()

const DECK_JSON = "deck.json"

var cardTypes = map[pointType]reflect.Type{
	PointsPerVegetableOne:   reflect.TypeOf(pointsPerVegetable{}),
	PointsPerVegetableTwo:   reflect.TypeOf(pointsPerVegetable{}),
	PointsPerVegetableThree: reflect.TypeOf(pointsPerVegetable{}),
	SumTwo:                  reflect.TypeOf(sum{}),
	SumThree:                reflect.TypeOf(sum{}),
	EvenOdd:                 reflect.TypeOf(evenOdd{}),
	Fewest:                  reflect.TypeOf(fewestMost{}),
	Most:                    reflect.TypeOf(fewestMost{}),
	FewestTotal:             reflect.TypeOf(other{}),
	MostTotal:               reflect.TypeOf(other{}),
	CompleteSet:             reflect.TypeOf(other{}),
	AtLeastTwo:              reflect.TypeOf(other{}),
	AtLeastThree:            reflect.TypeOf(other{}),
	MissingVegetable:        reflect.TypeOf(other{}),
}

var errInvalidVegetableType = errors.New("invalid vegetableType")

func loadDeck() map[string]cardInterface {
	deck := map[string]cardInterface{}

	f, err := os.Open(DECK_JSON)
	if err != nil {
		log.Fatalf("os.Open: %v", err)
	}
	defer f.Close()

	deckRaw := []map[string]any{}
	err = json.NewDecoder(f).Decode(&deckRaw)
	if err != nil {
		log.Fatalf("json.NewDecoder: %v", err)
	}

	for _, cardRaw := range deckRaw {
		card := unmarshalCard(cardRaw)
		deck[card.getID()] = card
	}

	return deck
}

func unmarshalCard(cardRaw map[string]any) cardInterface {
	cardTypeStr, ok := cardRaw["pointType"].(string)
	if !ok {
		slog.Error("cardRaw does not have pointType", slog.Any("cardRaw", cardRaw))
		os.Exit(1)
	}

	cardType := pointType(cardTypeStr)
	_, ok = AllPointTypes[cardType]
	if !ok {
		slog.Error("cardTypeStr is not a valid PointType", slog.Any("cardTypeStr", cardTypeStr))
		os.Exit(1)
	}

	card, ok := reflect.New(cardTypes[cardType]).Interface().(cardInterface)
	if !ok {
		slog.Error("reflect.New(cardTypes[cardType]).Interface().(cardInterface) failed")
		os.Exit(1)
	}

	cardData, err := json.Marshal(cardRaw)
	if err != nil {
		slog.Error("json.Marshal: %v", slog.Any("error", err))
		os.Exit(1)
	}

	if err = json.Unmarshal(cardData, card); err != nil {
		slog.Error("json.Unmarshal: %v", slog.Any("error", err))
		os.Exit(1)
	}

	if _, ok = AllVegetablesMap[card.getVegetable()]; !ok {
		slog.Error("card.getVegetable(): not a valid vegetableType", slog.Any("vegetableType", card.getVegetable()), slog.Any("card", card))
		os.Exit(1)
	}

	if err = card.validateDetails(); err != nil {
		slog.Error("card.validateDetails", slog.Any("error", err), slog.Any("card", card))
		os.Exit(1)
	}

	card.createDetailsMap()

	return card
}

type cardInterface interface {
	getID() string
	getVegetable() vegetableType
	getPointType() pointType
	createDetailsMap()
	getDetailsMap() map[string]interface{}
	validateDetails() error
}

type cardBase struct {
	CardID     string        `json:"cardID"`
	Vegetable  vegetableType `json:"vegetable"`
	PointType  pointType     `json:"pointType"`
	detailsMap map[string]interface{}
}

func (c *cardBase) getID() string {
	return c.CardID
}

func (c *cardBase) getVegetable() vegetableType {
	return c.Vegetable
}

func (c *cardBase) getPointType() pointType {
	return c.PointType
}

func (c *cardBase) createDetailsMap() {}

func (c *cardBase) getDetailsMap() map[string]interface{} {
	return c.detailsMap
}

func (c *cardBase) validateDetails() error {
	return nil
}

type pointsPerVegetable struct {
	cardBase
	Details []struct {
		Vegetable vegetableType `json:"vegetable"`
		Points    int           `json:"points"`
	} `json:"details"`
}

func (c *pointsPerVegetable) createDetailsMap() {
	detailsMap := map[string]interface{}{}
	for _, detail := range c.Details {
		detailsMap[string(detail.Vegetable)] = detail.Points
	}
	c.detailsMap = detailsMap
}

func (c *pointsPerVegetable) validateDetails() error {
	for _, detail := range c.Details {
		if _, ok := AllVegetablesMap[detail.Vegetable]; !ok {
			return errInvalidVegetableType
		}
	}
	return nil
}

type sum struct {
	cardBase
	Details struct {
		Vegetables []vegetableType `json:"vegetables"`
		Points     int             `json:"points"`
	} `json:"details"`
}

func (c *sum) createDetailsMap() {
	anySlice := []interface{}{}
	for _, vegetable := range c.Details.Vegetables {
		anySlice = append(anySlice, vegetable)
	}
	c.detailsMap = map[string]interface{}{
		"vegetables": anySlice,
		"points":     c.Details.Points,
	}
}

func (c *sum) validateDetails() error {
	for _, vegetable := range c.Details.Vegetables {
		if _, ok := AllVegetablesMap[vegetable]; !ok {
			return errInvalidVegetableType
		}
	}
	return nil
}

type evenOdd struct {
	cardBase
	Details struct {
		Vegetable vegetableType `json:"vegetable"`
		Even      int           `json:"even"`
		Odd       int           `json:"odd"`
	} `json:"details"`
}

func (c *evenOdd) createDetailsMap() {
	c.detailsMap = map[string]interface{}{
		"even": c.Details.Even,
		"odd":  c.Details.Odd,
	}
}

func (c *evenOdd) validateDetails() error {
	if _, ok := AllVegetablesMap[c.Details.Vegetable]; !ok {
		return errInvalidVegetableType
	}
	return nil
}

type fewestMost struct {
	cardBase
	Details struct {
		Vegetable vegetableType `json:"vegetable"`
		Points    int           `json:"points"`
	} `json:"details"`
}

func (c *fewestMost) createDetailsMap() {
	c.detailsMap = map[string]interface{}{
		"vegetable": c.Details.Vegetable,
		"points":    c.Details.Points,
	}
}

func (c *fewestMost) validateDetails() error {
	if _, ok := AllVegetablesMap[c.Details.Vegetable]; !ok {
		return errInvalidVegetableType
	}
	return nil
}

type other struct {
	cardBase
	Details struct {
		Points int `json:"points"`
	} `json:"details"`
}

func (c *other) createDetailsMap() {
	c.detailsMap = map[string]interface{}{
		"points": c.Details.Points,
	}
}

func (c *other) validateDetails() error {
	return nil
}

var cardsForPlayerCount = map[int]int{
	2: 6,
	3: 9,
	4: 12,
	5: 15,
	6: 18,
}

func (l *gameLogic) createPlayingDeck(playerCount int) []string {
	playingDeck := []string{}
	countForType := cardsForPlayerCount[playerCount]
	counts := map[vegetableType]int{}

	for _, card := range Deck {
		if counts[card.getVegetable()] < countForType {
			playingDeck = append(playingDeck, card.getID())
			counts[card.getVegetable()]++
		}

		if len(playingDeck) == playerCount*countForType*len(AllVegetables) {
			break
		}
	}

	return playingDeck
}
