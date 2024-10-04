package bot

import (
	"log/slog"

	"github.com/marcelbednarczyk/hackathon-jurata-2024/cards.go"
	"github.com/marcelbednarczyk/hackathon-jurata-2024/consequences"
	"github.com/marcelbednarczyk/hackathon-jurata-2024/proto"
)

type solucjaTymczasowaBot struct {
	greedy *greedyBot
	Cards  []cards.Card
}

func NewSolucjaTymczasowaBot(wczytaneKarty []cards.Card) *solucjaTymczasowaBot {
	return &solucjaTymczasowaBot{
		greedy: NewGreedyBot(),
		Cards:  wczytaneKarty,
	}
}

func (b *solucjaTymczasowaBot) MakeTakeCardsMove(gameState *proto.GameState, i int) []string {
	maxPoints, maxIds := 0, []string{}
	cons := consequences.Start(gameState)
	dynamicPriority := cons.Priorities

	for id, points := range cons.PointCards {
		priority := dynamicPriority.GetPointPriority(proto.VegetableType(id))
		slog.Info("Point card", slog.String("id", id), slog.Int("points", points), slog.Int("priority", int(priority)))
		adjustedPoints := float64(points) + priority
		if adjustedPoints > float64(maxPoints) && id != "" {
			maxPoints = int(adjustedPoints)
			maxIds = []string{id}
		}
	}
	for id, cards := range cons.VegetableCards {
		for id2, points := range cards {
			priority1 := dynamicPriority.GetVegetablePriority(proto.VegetableType(id))
			priority2 := dynamicPriority.GetVegetablePriority(proto.VegetableType(id2))

			slog.Info("Vegetable card", slog.String("id", id), slog.String("id2", id2), slog.Int("points", points))
			if points > maxPoints && id != "" && id2 != "" {
				maxPoints = points
				maxIds = []string{id, id2}
			}
		}
	}
	sum := 0
	if len(maxIds) == 0 || i < 3 {
		return b.greedy.MakeTakeCardsMove(gameState, i)
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
			return b.greedy.MakeTakeCardsMove(gameState, i)
		}
	}
}

func (b *solucjaTymczasowaBot) MakeFlipMove(state *proto.GameState) []string {
	return []string{}
}

//TODO: Dynamiczne priorytezowanie kart prawodobieństwo
