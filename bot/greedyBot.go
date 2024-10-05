package bot

import (
	"math/rand"

	"github.com/marcelbednarczyk/hackathon-jurata-2024/counter"
	"github.com/marcelbednarczyk/hackathon-jurata-2024/proto"
)

type greedyBot struct{}

func NewGreedyBot() *greedyBot {
	return &greedyBot{}
}

func (b *greedyBot) MakeTakeCardsMove(state *proto.GameState, _ counter.Counter, i int) []string {
	if ok, card := niceToTakePointCard(state, i); ok {
		return []string{card}
	}
	if (rand.Intn(2) == 0 && notEmpty(state.Market.PointCards)) || !notEmpty(state.Market.VegetableCards) {
		return []string{takeRandomNotNullCard(state.Market.PointCards)}
	}

	return takeNRandomNotNullCards(state.Market.VegetableCards, 2)
}

func niceToTakePointCard(state *proto.GameState, i int) (bool, string) {
	for _, card := range state.Market.PointCards {
		if card.PointsPerVegetable != nil {
			sum := 0
			for _, detail := range card.PointsPerVegetable.Points {
				sum += int(detail.Points)
			}
			if sum > 2 {
				// TODO: logic with couting points for us and opponents
				return true, card.CardID
			}
		}
	}
	return false, ""
}

func takeNicestCard(vegetables []*proto.VegtableHeld, cards []*proto.Card) string {
	for _, card := range cards {
		if card.PointsPerVegetable != nil {
			sum := 0
			for _, detail := range card.PointsPerVegetable.Points {
				sum += int(detail.Points)
			}
			if sum > 1 {
				return card.CardID
			}
		}
	}
	return takeRandomNotNullCard(cards)
}
func (b *greedyBot) MakeFlipMove(state *proto.GameState, _ int) []string {
	return []string{}
}
