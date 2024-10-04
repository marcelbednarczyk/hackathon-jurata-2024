package logic

import (
	"fmt"
	"log/slog"

	"github.com/marcelbednarczyk/hackathon-jurata-2024/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewGame(playerIDs []string, gameIndex int) *gameLogic {
	return newGame(playerIDs, gameIndex)
}

func (l *gameLogic) GetGameState(playerID string) (*proto.GameState, error) {
	otherPlayers := make([]*proto.Hand, 0)
	for pID := range l.playersHands {
		if pID != playerID {
			hand, err := playerHandToProto(l.playersHands[pID])
			if err != nil {
				return nil, err
			}
			otherPlayers = append(otherPlayers, hand)
		}
	}

	moveToMake := proto.MoveType_TAKE_CARDS
	if l.isFlipNext {
		moveToMake = proto.MoveType_FLIP_CARD
	}

	handProto, err := playerHandToProto(l.playersHands[playerID])
	if err != nil {
		return nil, err
	}

	marketProto, err := marketToProto(l.market)
	if err != nil {
		return nil, err
	}

	return &proto.GameState{
		PlayerToMove:   l.playerIDs[l.playerToMove],
		Market:         marketProto,
		IsGameOver:     l.isGameOver,
		Winner:         l.winner,
		YourHand:       handProto,
		OpponentsHands: otherPlayers,
		MoveToMake:     moveToMake,
		FinalScores:    l.finalScores,
	}, nil
}

func (l *gameLogic) MakeMove(move *proto.MoveRequest) error {
	if move == nil {
		connErr := fmt.Errorf("connection seems to be broken")
		slog.Error(connErr.Error())
		return status.Errorf(codes.Unknown, connErr.Error())
	}
	err := l.makeMove(move.PlayerID, move.MoveType == proto.MoveType_FLIP_CARD, move.Cards...)
	if err != nil {
		return err
	}

	return nil
}

func (l *gameLogic) GetPlayerToMove() string {
	return l.playerIDs[l.playerToMove]
}

func playerHandToProto(playerHand *playerHand) (*proto.Hand, error) {
	pHand := &proto.Hand{
		PointCards: make([]*proto.Card, len(playerHand.PointCards)),
		Vegetables: make([]*proto.VegtableHeld, 0),
	}

	for i, cardID := range playerHand.PointCards {
		c, err := cardToProto(Deck[cardID])
		if err != nil {
			return nil, err
		}
		pHand.PointCards[i] = c
	}

	for vegType, count := range playerHand.Vegetables {
		pHand.Vegetables = append(pHand.Vegetables, &proto.VegtableHeld{
			VegetableType: vegType.protoType(),
			Count:         int32(count),
		})
	}

	return pHand, nil
}

func (t vegetableType) protoType() proto.VegetableType {
	return proto.VegetableType(proto.VegetableType_value[t.UpperString()])
}

func (t pointType) protoType() proto.PointType {
	return proto.PointType(proto.PointType_value[string(t.UpperString())])
}

func cardToProto(c cardInterface) (*proto.Card, error) {
	card := &proto.Card{
		CardID:    c.getID(),
		Vegetable: c.getVegetable().protoType(),
		PointType: c.getPointType().protoType(),
	}

	if card.PointType == proto.PointType_BAD_POINT_TYPE {
		return nil, fmt.Errorf("wrong point type, card: %v", c)
	}

	if card.Vegetable == proto.VegetableType_BAD_VEGGIE {
		return nil, fmt.Errorf("wrong vegetable type, card: %v", c)
	}

	switch c := c.(type) {
	case *pointsPerVegetable:
		card.PointsPerVegetable = &proto.PointsPerVegetable{
			Points: make([]*proto.VegetablePoints, len(c.Details)),
		}

		for i, detail := range c.Details {
			card.PointsPerVegetable.Points[i] = &proto.VegetablePoints{
				Vegetable: detail.Vegetable.protoType(),
				Points:    int32(detail.Points),
			}
		}
	case *sum:
		card.Sum = &proto.Sum{
			Points:     int32(c.Details.Points),
			Vegetables: make([]proto.VegetableType, len(c.Details.Vegetables)),
		}

		for i, veg := range c.Details.Vegetables {
			card.Sum.Vegetables[i] = veg.protoType()
		}
	case *evenOdd:
		card.EvenOdd = &proto.EvenOdd{
			Odd:       int32(c.Details.Odd),
			Even:      int32(c.Details.Even),
			Vegetable: c.Details.Vegetable.protoType(),
		}
	case *fewestMost:
		card.FewestMost = &proto.FewestMost{
			Points:    int32(c.Details.Points),
			Vegetable: c.Details.Vegetable.protoType(),
		}
	case *other:
		card.Other = &proto.Other{
			Points: int32(c.Details.Points),
		}
	}

	return card, nil
}

func marketToProto(m *market) (*proto.Market, error) {
	pMarket := &proto.Market{
		PointCards:     make([]*proto.Card, MARKET_WIDTH),
		VegetableCards: make([]*proto.Card, MARKET_HEIGHT*MARKET_WIDTH),
	}

	for i := 0; i < MARKET_WIDTH; i++ {
		if m.pointCards[i] != "" {
			c, err := cardToProto(Deck[m.pointCards[i]])
			if err != nil {
				return nil, err
			}
			pMarket.PointCards[i] = c
		}

		for j := 0; j < MARKET_HEIGHT; j++ {
			if m.vegetableCards[i][j] != "" {
				c, err := cardToProto(Deck[m.vegetableCards[i][j]])
				if err != nil {
					return nil, err
				}
				pMarket.VegetableCards[i*MARKET_HEIGHT+j] = c
			}
		}
	}

	return pMarket, nil
}
