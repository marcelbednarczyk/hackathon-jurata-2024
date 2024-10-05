package bot

import (
	"math/rand"

	"github.com/marcelbednarczyk/hackathon-jurata-2024/counter"
	"github.com/marcelbednarczyk/hackathon-jurata-2024/proto"
)

type randomBot struct{}

func NewRandomBot() *randomBot {
	return &randomBot{}
}

func (b *randomBot) MakeTakeCardsMove(state *proto.GameState, _ counter.Counter, _ int) []string {
	if (rand.Intn(2) == 0 && notEmpty(state.Market.PointCards)) || !notEmpty(state.Market.VegetableCards) {
		return []string{takeRandomNotNullCard(state.Market.PointCards)}
	}

	return TakeNRandomNotNullCards(state.Market.VegetableCards, 2)
}

func (b *randomBot) MakeFlipMove(state *proto.GameState, _ counter.Counter, _ int) []string {
	return []string{}
}

func TakeNRandomNotNullCards(cards []*proto.Card, n int) []string {
	result := map[string]struct{}{}
	count := min(n, avaialbleCards(cards))
	for {
		if len(result) == count {
			break
		}

		card := takeRandomNotNullCard(cards)
		result[card] = struct{}{}
	}

	res := make([]string, 0, n)
	for card := range result {
		res = append(res, card)
	}
	return res
}

func takeRandomNotNullCard(cards []*proto.Card) string {
	for {
		card := cards[rand.Intn(len(cards))]
		if card.CardID != "" {
			return card.CardID
		}
	}
}

func notEmpty(cards []*proto.Card) bool {
	for _, card := range cards {
		if card.CardID != "" {
			return true
		}
	}

	return false
}

func avaialbleCards(cards []*proto.Card) int {
	count := 0
	for _, card := range cards {
		if card.CardID != "" {
			count++
		}
	}

	return count
}
