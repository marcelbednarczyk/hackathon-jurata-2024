package logic

import (
	"slices"
)

func (l *gameLogic) countPoints(playerID string) int {
	playerHand := l.playersHands[playerID]
	points := 0
	for _, cardID := range playerHand.PointCards {
		card := Deck[cardID]
		if card == nil {
			continue
		}
		points += scoreCard(playerID, l.playersHands, card)
	}
	return points
}

func scoreCard(playerID string, playersHands map[string]*playerHand, card cardInterface) int {
	playerVegetables := playersHands[playerID].Vegetables

	switch c := card.(type) {
	case *pointsPerVegetable:
		return scorePointsPerVegetable(playerVegetables, c)
	case *sum:
		return scoreSum(playerVegetables, c)
	case *evenOdd:
		return scoreEvenOdd(playerVegetables, c)
	case *fewestMost:
		counts := getCountsForVegetable(playersHands, c)
		switch c.PointType {
		case Fewest:
			return scoreFewest(playerVegetables, counts, c)
		case Most:
			return scoreMost(playerVegetables, counts, c)
		}
	case *other:
		switch c.PointType {
		case CompleteSet:
			return scoreCompleteSet(playerVegetables, c)
		case AtLeastTwo:
			return scoreAtLeast(playerVegetables, c, 2)
		case AtLeastThree:
			return scoreAtLeast(playerVegetables, c, 3)
		case MissingVegetable:
			return scoreMissing(playerVegetables, c)
		case FewestTotal:
			return scoreFewestTotal(playerID, playersHands, c)
		case MostTotal:
			return scoreMostTotal(playerID, playersHands, c)
		}
	}
	//unreachable
	return 0
}

func scorePointsPerVegetable(vegetables map[vegetableType]int, c *pointsPerVegetable) int {
	points := 0
	for _, detail := range c.Details {
		points += detail.Points * int(vegetables[detail.Vegetable])
	}
	return points
}

func scoreSum(vegetables map[vegetableType]int, c *sum) int {
	counts := []int{}
	for _, vegetable := range c.Details.Vegetables {
		counts = append(counts, int(vegetables[vegetable]))
	}

	count := slices.Min(counts)
	if isOneVegetableType(c.Details.Vegetables) {
		count /= len(c.Details.Vegetables)
	}

	return count * c.Details.Points
}

func isOneVegetableType(types []vegetableType) bool {
	if len(types) == 0 {
		return false
	}

	for _, t := range types {
		if t != types[0] {
			return false
		}
	}
	return true
}

func scoreEvenOdd(vegetables map[vegetableType]int, c *evenOdd) int {
	if vegetables[c.Details.Vegetable] == 0 {
		return c.Details.Odd
	}

	if vegetables[c.Details.Vegetable]%2 == 0 {
		return c.Details.Even
	}

	return c.Details.Odd
}

func getCountsForVegetable(playerHands map[string]*playerHand, c *fewestMost) map[string]int {
	counts := map[string]int{}
	for playerID, playerHand := range playerHands {
		counts[playerID] = playerHand.Vegetables[c.Details.Vegetable]
	}
	return counts
}

func scoreFewest(vegetables map[vegetableType]int, counts map[string]int, c *fewestMost) int {
	playersVeggie := int(vegetables[c.Details.Vegetable])
	for _, count := range counts {
		if playersVeggie > int(count) {
			return 0
		}
	}
	return c.Details.Points
}

func scoreMost(vegetables map[vegetableType]int, counts map[string]int, c *fewestMost) int {
	playersVeggie := int(vegetables[c.Details.Vegetable])
	for _, count := range counts {
		if playersVeggie < int(count) {
			return 0
		}
	}
	return c.Details.Points
}

func scoreCompleteSet(vegetables map[vegetableType]int, c *other) int {
	counts := []int{}
	for vegetable := range AllVegetablesMap {
		counts = append(counts, int(vegetables[vegetable]))
	}

	return slices.Min(counts) * c.Details.Points
}

func scoreAtLeast(vegetables map[vegetableType]int, c *other, atLeast int) int {
	points := 0
	for vegetable := range AllVegetablesMap {
		if vegetables[vegetable] >= atLeast {
			points += c.Details.Points
		}
	}

	return points
}

func scoreMissing(vegetables map[vegetableType]int, c *other) int {
	points := 0
	for vegetable := range AllVegetablesMap {
		if vegetables[vegetable] == 0 {
			points += c.Details.Points
		}
	}

	return points
}

func scoreFewestTotal(playerID string, playerHands map[string]*playerHand, c *other) int {
	counts := map[string]int{}
	countValues := []int{}
	for pID, playerHand := range playerHands {
		count := 0
		for _, vegetable := range playerHand.Vegetables {
			count += int(vegetable)
		}
		counts[pID] = count
		countValues = append(countValues, count)
	}

	min := slices.Min(countValues)
	if counts[playerID] == min {
		return c.Details.Points
	}
	return 0
}

func scoreMostTotal(playerID string, playerHands map[string]*playerHand, c *other) int {
	counts := map[string]int{}
	countValues := []int{}
	for pID, playerHand := range playerHands {
		count := 0
		for _, vegetable := range playerHand.Vegetables {
			count += int(vegetable)
		}
		counts[pID] = count
		countValues = append(countValues, count)
	}

	max := slices.Max(countValues)
	if counts[playerID] == max {
		return c.Details.Points
	}
	return 0
}
