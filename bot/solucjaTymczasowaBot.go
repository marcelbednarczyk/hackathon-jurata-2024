package bot

import (
	"os"

	"github.com/marcelbednarczyk/hackathon-jurata-2024/consequences"
	"github.com/marcelbednarczyk/hackathon-jurata-2024/counter"
	"github.com/marcelbednarczyk/hackathon-jurata-2024/proto"
)

type solucjaTymczasowaBot struct {
	greedy *greedyBot
}

func NewSolucjaTymczasowaBot() *solucjaTymczasowaBot {
	return &solucjaTymczasowaBot{
		greedy: NewGreedyBot(),
	}
}

func (b *solucjaTymczasowaBot) MakeTakeCardsMove(gameState *proto.GameState, cou counter.Counter, i int) []string {
	maxPoints, maxIds := 0, []string{}
	cons := consequences.Start(gameState, cou)
	for id, points := range cons.PointCards {
		// slog.Info("Point card", slog.String("id", id), slog.Int("points", points))
		if points > maxPoints && id != "" {
			maxPoints = points
			maxIds = []string{id}
		}
	}
	for id, cards := range cons.VegetableCards {
		for id2, points := range cards {
			// slog.Info("Vegetable card", slog.String("id", id), slog.String("id2", id2), slog.Int("points", points))
			if points > maxPoints && id != "" && id2 != "" {
				maxPoints = points
				maxIds = []string{id, id2}
			}
		}
	}
	sum := 0
	if len(maxIds) == 0 || (i < 3 && os.Getenv("PREDICT") != "true") {
		return b.greedy.MakeTakeCardsMove(gameState, cou, i)
	} else {
		for _, card := range gameState.Market.PointCards {
			if card.CardID == maxIds[0] {
				sum++
			}
		}
		for _, card := range gameState.Market.VegetableCards {
			for _, id := range maxIds {
				if card.CardID == id {
					sum++
				}
			}
		}
		if sum == len(maxIds) {
			return maxIds
		} else {
			return b.greedy.MakeTakeCardsMove(gameState, cou, i)
		}
	}
}

func (b *solucjaTymczasowaBot) MakeFlipMove(state *proto.GameState, i int) []string {
	if i < 13 {
		return []string{}
	}
	return []string{}
	// TODO: implement Flip function
	maxPoints, maxIds := 0, []string{}
	cons := consequences.Flip(state)
	for id, points := range cons.PointCards {
		// slog.Info("Point card", slog.String("id", id), slog.Int("points", points))
		if points > maxPoints && id != "" {
			maxPoints = points
			maxIds = []string{id}
		}
	}
	sum := 0
	if len(maxIds) == 0 {
		return []string{}
	} else {
		for _, card := range state.Market.PointCards {
			if card.CardID == maxIds[0] {
				sum++
			}
		}
		if sum == len(maxIds) {
			return maxIds
		} else {
			return []string{}
		}
	}
}
