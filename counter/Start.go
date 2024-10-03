package counter

import (
	"strconv"

	"gitlab.com/hackathon-rainbow-2024/go-client/proto"
)

type Counter struct {
	TomatosRemaining  int
	CarrotsRemaining  int
	LettucesRemaining int
	CabbagesRemaining int
	PeppersRemaining  int
	OnionsRemaining   int
	LastGameState     []proto.Card
}

func Start(gameState *proto.GameState) Counter {
	c := Counter{
		TomatosRemaining:  6,
		CarrotsRemaining:  6,
		LettucesRemaining: 6,
		CabbagesRemaining: 6,
		PeppersRemaining:  6,
		OnionsRemaining:   6,
	}
	c.setLastGameState(gameState)
	for _, card := range c.LastGameState {
		switch card.Vegetable {
		case proto.VegetableType_TOMATO:
			c.TomatosRemaining--
		case proto.VegetableType_CARROT:
			c.CarrotsRemaining--
		case proto.VegetableType_LETTUCE:
			c.LettucesRemaining--
		case proto.VegetableType_CABBAGE:
			c.CabbagesRemaining--
		case proto.VegetableType_PEPPER:
			c.PeppersRemaining--
		case proto.VegetableType_ONION:
			c.OnionsRemaining--
		}
	}
	return c
}

func (c *Counter) String() string {
	return "Tomatos: " + strconv.Itoa(c.TomatosRemaining) + ", Carrots: " + strconv.Itoa(c.CarrotsRemaining) + ", Lettuces: " + strconv.Itoa(c.LettucesRemaining) + ", Cabbages: " + strconv.Itoa(c.CabbagesRemaining) + ", Peppers: " + strconv.Itoa(c.PeppersRemaining) + ", Onions: " + strconv.Itoa(c.OnionsRemaining)
}

func (c *Counter) Update(gameState *proto.GameState) {
	changes := c.getChangesInGameState(gameState)
	for _, card := range changes {
		switch card.Vegetable {
		case proto.VegetableType_TOMATO:
			c.TomatosRemaining--
		case proto.VegetableType_CARROT:
			c.CarrotsRemaining--
		case proto.VegetableType_LETTUCE:
			c.LettucesRemaining--
		case proto.VegetableType_CABBAGE:
			c.CabbagesRemaining--
		case proto.VegetableType_PEPPER:
			c.PeppersRemaining--
		case proto.VegetableType_ONION:
			c.OnionsRemaining--
		}
	}
	c.setLastGameState(gameState)
}

func (c *Counter) getChangesInGameState(newGameState *proto.GameState) []proto.Card {
	changes := []proto.Card{}
	for _, card := range newGameState.Market.VegetableCards {
		if !contains(c.LastGameState, card) {
			changes = append(changes, *card)
		}
	}
	for _, card := range newGameState.Market.PointCards {
		if !contains(c.LastGameState, card) {
			changes = append(changes, *card)
		}
	}
	return changes
}

func contains(cards []proto.Card, card *proto.Card) bool {
	for _, c := range cards {
		if c.CardID == card.CardID {
			return true
		}
	}
	return false
}

func (c *Counter) setLastGameState(gameState *proto.GameState) {
	c.LastGameState = []proto.Card{}
	for _, card := range gameState.Market.VegetableCards {
		c.LastGameState = append(c.LastGameState, *card)
	}
	for _, card := range gameState.Market.PointCards {
		c.LastGameState = append(c.LastGameState, *card)
	}
}
