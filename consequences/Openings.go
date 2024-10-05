package consequences

import (
	"log/slog"
	"sort"

	"github.com/marcelbednarczyk/hackathon-jurata-2024/proto"
)

func OpeningVegeCard(gameState *proto.GameState) []string {
	if len(gameState.YourHand.PointCards) == 0 {
		return []string{}
	}

	positiveBonusesInHand := map[proto.VegetableType]int{}
	for _, card := range gameState.YourHand.PointCards {
		if ppv := card.PointsPerVegetable; ppv != nil {
			for _, vp := range ppv.Points {
				if vp.Points > 0 {
					positiveBonusesInHand[vp.Vegetable] += int(vp.Points)
				}
			}
		}
	}

	sortedVegetableCards := make([]*proto.Card, 0)
	for _, card := range gameState.Market.VegetableCards {
		card.Value = calculateCardValue(card, positiveBonusesInHand)
		if card.Value > 0 {
			sortedVegetableCards = append(sortedVegetableCards, card)
		}
	}

	sort.Slice(sortedVegetableCards, func(i, j int) bool {
		return sortedVegetableCards[i].Value > sortedVegetableCards[j].Value
	})

	selectedCards := []string{}
	if len(sortedVegetableCards) >= 2 {
		selectedCards = append(selectedCards, sortedVegetableCards[0].CardID, sortedVegetableCards[1].CardID)
		slog.Info("Opening", slog.String("bestCardID1", selectedCards[0]), slog.String("bestCardID2", selectedCards[1]))
	} else if len(sortedVegetableCards) == 1 {
		selectedCards = append(selectedCards, sortedVegetableCards[0].CardID)
	}

	return selectedCards
}

func calculateCardValue(card *proto.Card, bonusesInHand map[proto.VegetableType]int) int32 {
	value := 0
	if ppv := card.PointsPerVegetable; ppv != nil {
		for _, vp := range ppv.Points {
			if bonus, exists := bonusesInHand[vp.Vegetable]; !exists || bonus == 0 {
				value += 1
			}
		}
	}
	return int32(value)
}

func OpeningPointCard(gameState *proto.GameState) string {
	pointCardBonusesInHand := map[proto.VegetableType]int{}
	for _, card := range gameState.YourHand.PointCards {
		if ppv := card.PointsPerVegetable; ppv != nil {
			for _, vp := range ppv.Points {
				pointCardBonusesInHand[vp.Vegetable] += int(vp.Points)
			}
		}
	}
	slog.Info("Point card bonuses in hand", slog.Any("pointCardBonusesInHand", pointCardBonusesInHand))

	sortedPointCards := make([]*proto.Card, len(gameState.Market.PointCards))
	copy(sortedPointCards, gameState.Market.PointCards)
	sort.Slice(sortedPointCards, func(i, j int) bool {
		return calculateCardScore(sortedPointCards[i], pointCardBonusesInHand) > calculateCardScore(sortedPointCards[j], pointCardBonusesInHand)
	})

	bestCardID := ""
	maxScore := 0

	for _, card := range sortedPointCards {
		if card.PointType != proto.PointType_BAD_POINT_TYPE {
			modifier := shouldAddCard(card, pointCardBonusesInHand)
			currentScore := calculateCardScore(card, pointCardBonusesInHand) + modifier
			if currentScore > maxScore {
				maxScore = currentScore
				bestCardID = card.CardID
			}
		}
	}
	slog.Info("Opening", slog.String("bestCardID", bestCardID), slog.Int("maxScore", maxScore))
	return bestCardID
}

func calculateCardScore(card *proto.Card, bonusesInHand map[proto.VegetableType]int) int {
	score := baseCardScore(card, bonusesInHand)
	if ppv := card.PointsPerVegetable; ppv != nil {
		uniqueVegetables := make(map[proto.VegetableType]bool)
		for _, vp := range ppv.Points {
			uniqueVegetables[vp.Vegetable] = true
			score += int(vp.Points)
			if bonus, ok := bonusesInHand[vp.Vegetable]; ok && bonus > 0 && vp.Points < 0 {
				score += int(vp.Points)
			} else if bonus, ok := bonusesInHand[vp.Vegetable]; ok && bonus < 0 && vp.Points > 0 {
				score += bonus
			}
		}
		if len(uniqueVegetables) == 3 {
			score += 3
		}
	}

	if card.Sum != nil {
		for _, veg := range card.Sum.Vegetables {
			if bonus, ok := bonusesInHand[veg]; ok && bonus > 0 {
				score += bonus
			}
		}
	}

	return score
}
func baseCardScore(card *proto.Card, bonusesInHand map[proto.VegetableType]int) int {
	switch card.PointType {
	case proto.PointType_POINTS_PER_VEGETABLE_THREE:
		return 6 + bonusModifier(card, bonusesInHand)
	case proto.PointType_POINTS_PER_VEGETABLE_TWO:
		return 3 + bonusModifier(card, bonusesInHand)
	case proto.PointType_POINTS_PER_VEGETABLE_ONE:
		return 1 + bonusModifier(card, bonusesInHand)

	case proto.PointType_SUM_THREE:
		return evaluateSumCard(card.Sum, 8, bonusesInHand)
	case proto.PointType_SUM_TWO:
		return evaluateSumCard(card.Sum, 3, bonusesInHand)

	case proto.PointType_MOST_TOTAL:
		return 6 + bonusModifier(card, bonusesInHand)
	case proto.PointType_FEWEST_TOTAL:
		return 6 + bonusModifier(card, bonusesInHand)

	// case proto.PointType_COMPLETE_SET:
	// 	return -1 + bonusModifier(card, bonusesInHand)
	// case proto.PointType_MISSING_VEGETABLE:
	// 	return missingVegetableModifier(card, bonusesInHand)

	// case proto.PointType_AT_LEAST_TWO:
	// 	return conditionalBonus(card, 2)
	// case proto.PointType_AT_LEAST_THREE:
	// 	return conditionalBonus(card, 3)

	default:
		return 0
	}
}

func bonusModifier(card *proto.Card, bonusesInHand map[proto.VegetableType]int) int {
	score := 0
	if ppv := card.PointsPerVegetable; ppv != nil {
		for _, vp := range ppv.Points {
			if bonus, ok := bonusesInHand[vp.Vegetable]; ok && bonus > 0 {
				score += 2
			}
		}
	}
	return score
}

func evaluateSumCard(sum *proto.Sum, baseValue int, bonusesInHand map[proto.VegetableType]int) int {
	uniqueVegetables := make(map[proto.VegetableType]struct{})
	score := baseValue

	for _, vegetable := range sum.Vegetables {
		uniqueVegetables[vegetable] = struct{}{}
		if _, exists := bonusesInHand[vegetable]; !exists {
			score += 1
		}
	}

	if len(uniqueVegetables) >= 3 {
		score += 3
	}
	return score
}

func missingVegetableModifier(card *proto.Card, bonusesInHand map[proto.VegetableType]int) int {
	score := 0
	if ppv := card.PointsPerVegetable; ppv != nil {
		for _, vp := range ppv.Points {
			if _, exists := bonusesInHand[vp.Vegetable]; !exists {
				score += 3
			}
		}
	}
	return score
}

func conditionalBonus(card *proto.Card, requiredCount int) int {
	if card.Sum != nil && len(card.Sum.Vegetables) >= requiredCount {
		return requiredCount * 2
	}
	return 0
}

// Sprawdza czy już nie ma debufów na dany rodaj warzywa
func shouldAddCard(card *proto.Card, bonusesInHand map[proto.VegetableType]int) int {
	if ppv := card.PointsPerVegetable; ppv != nil {
		for _, vp := range ppv.Points {
			if bonus, ok := bonusesInHand[vp.Vegetable]; ok && bonus > 0 {
				return -bonus
			}
		}
	}
	return 0
}
