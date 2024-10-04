package consequences

import (
	"log/slog"

	"github.com/marcelbednarczyk/hackathon-jurata-2024/proto"
)

func SimulateMove(gameState *proto.GameState, cards ...*proto.Card) *proto.GameState {
	if gameState.MoveToMake == proto.MoveType_FLIP_CARD {
		slog.Warn("Cannot simulate flip move")
		return nil
	}
	if len(cards) == 1 {
		gameState.YourHand.PointCards = append(gameState.YourHand.PointCards, cards[0])
		gameState.OpponentsHands[0].PointCards = append(gameState.OpponentsHands[0].PointCards, cards[0])
	} else {
		slog.Warn("TODO: implement taking two cards")
	}
	return gameState
}
