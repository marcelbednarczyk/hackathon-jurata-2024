package logic

const (
	MARKET_WIDTH  = 3
	MARKET_HEIGHT = 2
)

type market struct {
	pointCards     [MARKET_WIDTH]string
	vegetableCards [MARKET_WIDTH][MARKET_HEIGHT]string
}

func createEmptyMarket() *market {
	return &market{
		pointCards:     [MARKET_WIDTH]string{},
		vegetableCards: [MARKET_WIDTH][MARKET_HEIGHT]string{},
	}
}

func (m *market) fill(deck []string) []string {
	for m.isFillable(len(deck)) {

		// fill first vegetable card
		for i := 0; i < MARKET_WIDTH; i++ {
			isFilled := false
			for j := 0; j < MARKET_HEIGHT; j++ {
				if m.vegetableCards[i][j] == "" && m.pointCards[i] != "" {
					m.vegetableCards[i][j] = m.pointCards[i]
					m.pointCards[i] = ""
					isFilled = true
					break
				}
			}
			if isFilled {
				break
			}
		}

		// fill first point card
		if len(deck) > 0 {
			for i := 0; i < MARKET_WIDTH; i++ {
				if m.pointCards[i] == "" {
					m.pointCards[i] = deck[0]
					deck = deck[1:]
					break
				}
			}
		}

	}
	return deck
}

func (m *market) isFillable(cardsInDeck int) bool {
	if m.countEmptyPlaces() == 0 {
		return false
	}

	if cardsInDeck >= m.countEmptyPlaces() {
		return true
	}

	if cardsInDeck > 0 {
		for i := 0; i < MARKET_WIDTH; i++ {
			if m.pointCards[i] == "" {
				return true
			}

			for j := 0; j < MARKET_HEIGHT; j++ {
				if m.vegetableCards[i][j] == "" {
					return true
				}
			}
		}
		// unreachable
		return false
	}

	for i := 0; i < MARKET_WIDTH; i++ {
		if m.pointCards[i] == "" {
			continue
		}

		isFull := true
		for j := 0; j < MARKET_HEIGHT; j++ {
			if m.vegetableCards[i][j] == "" {
				isFull = false
				break
			}
		}
		if !isFull {
			return true
		}
	}

	return false
}

func (m *market) countEmptyPlaces() int {
	count := 0
	for i := 0; i < MARKET_WIDTH; i++ {
		if m.pointCards[i] == "" {
			count++
		}

		for j := 0; j < MARKET_HEIGHT; j++ {
			if m.vegetableCards[i][j] == "" {
				count++
			}
		}
	}
	return count
}

func (m *market) veggieCount() int {
	count := 0
	for i := 0; i < MARKET_WIDTH; i++ {
		for j := 0; j < MARKET_HEIGHT; j++ {
			if m.vegetableCards[i][j] != "" {
				count++
			}
		}
	}
	return count
}

func (m *market) cardInPointCards(cardID string) bool {
	for i := 0; i < MARKET_WIDTH; i++ {
		if m.pointCards[i] == cardID {
			return true
		}
	}
	return false
}

func (m *market) cardInVegetableCards(cardID string) bool {
	for i := 0; i < MARKET_WIDTH; i++ {
		for j := 0; j < MARKET_HEIGHT; j++ {
			if m.vegetableCards[i][j] == cardID {
				return true
			}
		}
	}
	return false
}
