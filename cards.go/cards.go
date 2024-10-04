package cards

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Card struct {
	CardID    string   `json:"cardId"`
	Vegetable string   `json:"vegetable"`
	PointType string   `json:"pointType"`
	Details   []Detail `json:"details,omitempty"`
}

type Detail struct {
	Vegetable string `json:"vegetable"`
	Points    int    `json:"points"`
}

func LoadCards(filename string) ([]Card, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var cards []Card
	err = json.Unmarshal(byteValue, &cards)
	if err != nil {
		return nil, err
	}

	return cards, nil
}

func VegetableType(cards []Card, cardId string) (string, error) {
	for _, card := range cards {
		if card.CardID == cardId {
			return card.Vegetable, nil
		}
	}
	return "", fmt.Errorf("cardId %s not found", cardId)
}
