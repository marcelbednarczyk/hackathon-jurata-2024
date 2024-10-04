package consequences

import (
	"github.com/marcelbednarczyk/hackathon-jurata-2024/counter"
	"github.com/marcelbednarczyk/hackathon-jurata-2024/proto"
)

type DynamicCardPriority struct {
	cardPriority     map[proto.VegetableType]float64
	pointPriority    map[proto.VegetableType]float64
	vegetableCounter *counter.Counter
}

func NewDynamicCardPriority(vegetableCounter *counter.Counter) *DynamicCardPriority {
	return &DynamicCardPriority{
		cardPriority:     make(map[proto.VegetableType]float64),
		pointPriority:    make(map[proto.VegetableType]float64),
		vegetableCounter: vegetableCounter,
	}
}
func (d *DynamicCardPriority) UpdatePriorities(yourHand *proto.Hand, opponentHands []*proto.Hand, market *proto.Market) {
	vegetableCount := make(map[proto.VegetableType]int)
	pointNeed := make(map[proto.VegetableType]int)

	// Zlicz warzywa w Twojej ręce oraz przeciwników
	for _, vegetable := range yourHand.Vegetables {
		vegetableCount[vegetable.VegetableType] += int(vegetable.Count)
	}
	for _, hand := range opponentHands {
		for _, vegetable := range hand.Vegetables {
			vegetableCount[vegetable.VegetableType] += int(vegetable.Count)
		}
	}

	for _, pointCard := range yourHand.PointCards {
		if pointCard.PointsPerVegetable != nil {
			for _, vegPoints := range pointCard.PointsPerVegetable.Points {
				needed := vegPoints.Points - int32(vegetableCount[vegPoints.Vegetable])
				if needed > 0 {
					pointNeed[vegPoints.Vegetable] += int(needed)
				}
			}
		}
	}

	for vegetableType, remaining := range vegetableCount {
		if remaining == 0 && pointNeed[vegetableType] > 0 {
			d.cardPriority[vegetableType] = float64(10 * pointNeed[vegetableType])
		} else {
			d.cardPriority[vegetableType] = float64(1) / (float64(remaining) + float64(pointNeed[vegetableType]))
		}
	}

	for _, card := range market.PointCards {
		if card.PointsPerVegetable == nil {
			continue
		}
		for _, vegPoints := range card.PointsPerVegetable.Points {
			d.adjustPriorityForVegetable(vegPoints.Vegetable)
		}
	}
}

func (d *DynamicCardPriority) adjustPriorityForVegetable(vegetable proto.VegetableType) {
	switch vegetable {
	case proto.VegetableType_TOMATO:
		if d.vegetableCounter.TomatosRemaining == 0 {
			d.cardPriority[proto.VegetableType_TOMATO] += 5.0
		}
	case proto.VegetableType_CARROT:
		if d.vegetableCounter.CarrotsRemaining == 0 {
			d.cardPriority[proto.VegetableType_CARROT] += 5.0
		}
	case proto.VegetableType_PEPPER:
		if d.vegetableCounter.PeppersRemaining == 0 {
			d.cardPriority[proto.VegetableType_PEPPER] += 5.0
		}
	case proto.VegetableType_ONION:
		if d.vegetableCounter.OnionsRemaining == 0 {
			d.cardPriority[proto.VegetableType_ONION] += 5.0
		}
	case proto.VegetableType_CABBAGE:
		if d.vegetableCounter.CabbagesRemaining == 0 {
			d.cardPriority[proto.VegetableType_CABBAGE] += 5.0
		}
	case proto.VegetableType_LETTUCE:
		if d.vegetableCounter.LettucesRemaining == 0 {
			d.cardPriority[proto.VegetableType_LETTUCE] += 5.0
		}
	}
}

func (d *DynamicCardPriority) GetVegetablePriority(vegetableType proto.VegetableType) float64 {
	if priority, exists := d.cardPriority[vegetableType]; exists {
		return priority
	}
	return -1
}

func (d *DynamicCardPriority) GetPointPriority(vegetableType proto.VegetableType) float64 {
	if priority, exists := d.pointPriority[vegetableType]; exists {
		return priority
	}
	return -1
}
