package consequences

import (
	"github.com/marcelbednarczyk/hackathon-jurata-2024/logic"
	"github.com/marcelbednarczyk/hackathon-jurata-2024/proto"
)

type Consequences struct {
	PointCards     map[string]int
	VegetableCards map[string]map[string]int
}

func Start(gameState *proto.GameState) Consequences {
	c := Consequences{
		PointCards:     make(map[string]int),
		VegetableCards: make(map[string]map[string]int),
	}
	for _, card := range gameState.Market.PointCards {
		logicState := logic.GetStateLog(SimulateMove(gameState, card))
		c.PointCards[card.CardID] = logicState.Hands["me"].CurrentScore + logicState.Hands["opponent"].CurrentScore
	}

	// for i, card := range gameState.Market.VegetableCards {
	// 	for j := i + 1; j < len(gameState.Market.VegetableCards); j++ {
	// 		logicState := logic.GetStateLog(SimulateMove(gameState, card, gameState.Market.VegetableCards[j]))
	// 		if _, ok := c.VegetableCards[card.CardID]; !ok {
	// 			c.VegetableCards[card.CardID] = make(map[string]int)
	// 		}
	// 		c.VegetableCards[card.CardID][gameState.Market.VegetableCards[j].CardID] = logicState.Hands["me"].CurrentScore + logicState.Hands["opponent"].CurrentScore
	// 	}
	// }

	return c
}
