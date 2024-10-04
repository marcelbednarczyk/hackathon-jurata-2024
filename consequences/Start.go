package consequences

import (
	"encoding/json"
	"log/slog"

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
	bytes, _ := json.Marshal(gameState)
	logicState := logic.GetStateLog(gameState)
	for _, card := range gameState.Market.PointCards {
		var game proto.GameState
		json.Unmarshal(bytes, &game)
		newState := SimulateMove(game, card)
		if newState.Market == nil {
			slog.Warn("Cannot simulate move without market")
			return Consequences{}
		}
		newLogicState := logic.GetStateLog(&newState)
		c.PointCards[card.CardID] =
			newLogicState.Hands["me"].CurrentScore -
				-logicState.Hands["me"].CurrentScore
	}

	for i, card := range gameState.Market.VegetableCards {
		for j, card2 := range gameState.Market.VegetableCards {
			if i >= j {
				continue
			}
			var game proto.GameState
			json.Unmarshal(bytes, &game)
			newState := SimulateMove(game, card, card2)
			if newState.Market == nil {
				slog.Warn("Cannot simulate move without market")
				return Consequences{}
			}
			newLogicState := logic.GetStateLog(&newState)
			if _, ok := c.VegetableCards[card.CardID]; !ok {
				c.VegetableCards[card.CardID] = make(map[string]int)
			}
			c.VegetableCards[card.CardID][card2.CardID] =
				newLogicState.Hands["me"].CurrentScore -
					-logicState.Hands["me"].CurrentScore
		}
	}

	return c
}
