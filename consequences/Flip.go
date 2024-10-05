package consequences

import (
	"encoding/json"

	"github.com/marcelbednarczyk/hackathon-jurata-2024/logic"
	"github.com/marcelbednarczyk/hackathon-jurata-2024/proto"
)

func Flip(gameState *proto.GameState) Consequences {
	c := Consequences{
		PointCards: make(map[string]int),
	}

	logicState := logic.GetStateLog(gameState)
	for i, card := range gameState.YourHand.PointCards {
		bytes, _ := json.Marshal(gameState)
		var game proto.GameState
		json.Unmarshal(bytes, &game)

		veggieHand(game.YourHand, game.YourHand.PointCards[i])
		game.YourHand.PointCards = removeCard(game.YourHand.PointCards, i)
		newLogicState := logic.GetStateLog(&game)

		c.PointCards[card.CardID] =
			newLogicState.Hands["me"].CurrentScore -
				-logicState.Hands["me"].CurrentScore
	}

	return c
}

func removeCard(cards []*proto.Card, i int) []*proto.Card {
	return append(cards[:i], cards[i+1:]...)
}
