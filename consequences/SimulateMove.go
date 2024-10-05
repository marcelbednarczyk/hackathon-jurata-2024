package consequences

import (
	"log/slog"

	"github.com/marcelbednarczyk/hackathon-jurata-2024/proto"
)

func SimulateMove(gameState proto.GameState, cards ...*proto.Card) proto.GameState {
	if gameState.MoveToMake == proto.MoveType_FLIP_CARD {
		slog.Warn("Cannot simulate flip move")
		return proto.GameState{}
	}
	for _, card := range cards {
		if card.CardID == "" {
			slog.Warn("Cannot simulate move with empty card")
			return proto.GameState{}
		}
	}
	if len(cards) == 1 {
		gameState.YourHand.PointCards = append(gameState.YourHand.PointCards, cards[0])
	} else {
		for i := 0; i < len(cards); i++ {
			veggieHand(gameState.YourHand, cards[i])
		}
	}
	return gameState
}

func veggieHand(hand *proto.Hand, card *proto.Card) {
	found := false
	for _, held := range hand.Vegetables {
		if held.VegetableType == card.Vegetable {
			held.Count++
			found = true
			break
		}
	}
	if !found {
		hand.Vegetables = append(hand.Vegetables, &proto.VegtableHeld{
			VegetableType: card.Vegetable,
			Count:         1,
		})
	}
}
