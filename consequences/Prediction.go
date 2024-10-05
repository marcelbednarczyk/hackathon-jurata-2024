package consequences

import (
	"os"
	"strconv"

	"github.com/marcelbednarczyk/hackathon-jurata-2024/counter"
	"github.com/marcelbednarczyk/hackathon-jurata-2024/proto"
)

var cardChance = 0.73 //37% 57/43, 50% 71/29, 75% 66/34

func init() {
	if os.Getenv("CARD_CHANCE") != "" {
		var err error
		cardChance, err = strconv.ParseFloat(os.Getenv("CARD_CHANCE"), 64)
		if err != nil {
			panic(err)
		}
	}
}

func Prediction(gameState *proto.GameState, remainingCards counter.Counter) {
	predictedCards := []*proto.Card{}
	predictedCards = append(predictedCards, getWholePrediction(remainingCards, cardChance)...)
	for _, card := range predictedCards {
		veggieHand(gameState.YourHand, card)
	}
}

func getWholePrediction(remainingCards counter.Counter, chance float64) []*proto.Card {
	predictedCards := []*proto.Card{}
	// Cabbage
	count := int(float64(remainingCards.CabbagesRemaining) * chance)
	predictedCards = append(predictedCards, convertToVeggieCards(proto.VegetableType_CABBAGE, count)...)

	// Carrot
	count = int(float64(remainingCards.CarrotsRemaining) * chance)
	predictedCards = append(predictedCards, convertToVeggieCards(proto.VegetableType_CARROT, count)...)

	// Onion
	count = int(float64(remainingCards.OnionsRemaining) * chance)
	predictedCards = append(predictedCards, convertToVeggieCards(proto.VegetableType_ONION, count)...)

	// Pepper
	count = int(float64(remainingCards.PeppersRemaining) * chance)
	predictedCards = append(predictedCards, convertToVeggieCards(proto.VegetableType_PEPPER, count)...)

	// Tomato
	count = int(float64(remainingCards.TomatosRemaining) * chance)
	predictedCards = append(predictedCards, convertToVeggieCards(proto.VegetableType_TOMATO, count)...)

	// Lettuce
	count = int(float64(remainingCards.LettucesRemaining) * chance)
	predictedCards = append(predictedCards, convertToVeggieCards(proto.VegetableType_LETTUCE, count)...)

	return predictedCards
}

func convertToVeggieCards(veg proto.VegetableType, count int) []*proto.Card {
	cards := []*proto.Card{}
	for i := 0; i < count; i++ {
		cards = append(cards, &proto.Card{
			Vegetable: veg,
		})
	}
	return cards
}
