package logic

import (
	"fmt"
	"log/slog"
	"time"

	"slices"

	"github.com/google/uuid"
	"github.com/marcelbednarczyk/hackathon-jurata-2024/proto"
)

type gameLogic struct {
	currentDeck       []cardID
	market            *market
	playersHands      map[string]*playerHand
	playerToMove      int
	isFlipNext        bool
	playerIDs         []string
	finalScores       map[string]int32
	invalidMoves      map[string]int
	winner            string
	isGameOver        bool
	startingPlayerIdx int
}

type MatchLog struct {
	MatchUUID        string        `json:"matchUUID"`
	RoomUUID         uuid.UUID     `json:"roomUUID"`
	RoomID           string        `json:"roomID"`
	RoomCreationDate time.Time     `json:"roomCreationDate"`
	Winner           string        `json:"winner"`
	Players          []MatchPlayer `json:"players"`
	Games            []MatchGame   `json:"games,omitempty"`
}

type MatchPlayer struct {
	PlayerUUID  string `json:"playerUUID"`
	PlayerName  string `json:"playerName"`
	PlayerScore int    `json:"playerScore"`
}

type GamePlayer struct {
	PlayerUUID string `json:"playerUUID"`
	PlayerName string `json:"playerName"`
	FinalScore int32  `json:"finalScore"`
	RoundsWon  int    `json:"roundsWon"`
}

type MatchGame struct {
	Index    int          `json:"index"`
	GameUUID string       `json:"gameUUID"`
	Winner   string       `json:"winner"`
	Players  []GamePlayer `json:"players"`
}

type GameLog struct {
	Index    int          `json:"index"`
	GameUUID string       `json:"gameUUID"`
	Winner   string       `json:"winner"`
	Players  []GamePlayer `json:"players"`
	States   []State      `json:"states"`
}

type State struct {
	PointCards     [MARKET_WIDTH]string                `json:"pointCards"`
	VegetableCards [MARKET_WIDTH][MARKET_HEIGHT]string `json:"vegetableCards"`
	PlayerTurn     string                              `json:"playerTurn"`
	IsFlip         bool                                `json:"isFlip"`
	Hands          map[string]playerHand               `json:"hands"`
}

// type Hand struct {
// 	PointCards []string `json:""`
// 	Vegetables map[string]int `json:""`
// }

type playerHand struct {
	CurrentScore int                   `json:"currentScore"`
	PointCards   []string              `json:"pointCards"`
	Vegetables   map[vegetableType]int `json:"vegetables"`
}

type cardID = string

const MAX_INVALID_MOVES = 10

func newGame(playerIDs []string, gameIndex int) *gameLogic {
	l := &gameLogic{}
	l.market = createEmptyMarket()
	l.playerToMove = gameIndex % len(playerIDs)
	l.startingPlayerIdx = l.playerToMove
	l.playerIDs = playerIDs
	l.isGameOver = false
	l.finalScores = map[string]int32{}
	l.winner = ""
	l.currentDeck = l.createPlayingDeck(len(playerIDs))
	l.playersHands = map[string]*playerHand{}
	l.invalidMoves = map[string]int{}

	for _, playerID := range playerIDs {
		hand := &playerHand{
			PointCards: []string{},
			Vegetables: map[vegetableType]int{},
		}

		for _, vegType := range AllVegetables {
			hand.Vegetables[vegType] = 0
		}
		l.invalidMoves[playerID] = 0

		l.playersHands[playerID] = hand
	}

	l.currentDeck = l.market.fill(l.currentDeck)

	return l
}

func (l *gameLogic) validateMove(playerID string, isFlip bool, cardIDs ...cardID) error {
	if l.playerIDs[l.playerToMove] != playerID {
		return fmt.Errorf("not your turn")
	}

	for _, cardID := range cardIDs {
		if _, ok := Deck[cardID]; !ok {
			return fmt.Errorf("card not found %v", cardID)
		}
	}

	if isFlip && !l.isFlipNext {
		return fmt.Errorf("cannot flip now")
	}

	if !isFlip && l.isFlipNext {
		return fmt.Errorf("must flip now")
	}

	// validate flip
	if isFlip {
		if len(cardIDs) == 0 {
			return nil
		}

		if len(cardIDs) != 1 {
			return fmt.Errorf("invalid number of cards")
		}

		playerHand := l.playersHands[playerID]
		for _, cardID := range playerHand.PointCards {
			if cardID == cardIDs[0] {
				return nil
			}
		}
		return fmt.Errorf("card not found in player hand")
	}

	if len(cardIDs) < 1 || len(cardIDs) > 2 {
		return fmt.Errorf("invalid number of cards")
	}

	if len(cardIDs) == 1 && l.market.cardInPointCards(cardIDs[0]) {
		return nil
	}

	if len(cardIDs) < 2 && l.market.veggieCount() > 1 {
		return fmt.Errorf("must take 2 cards")
	}

	if len(slices.Compact(cardIDs)) != len(cardIDs) {
		return fmt.Errorf("duplicate cards")
	}

	for _, cardID := range cardIDs {
		if !l.market.cardInVegetableCards(cardID) {
			return fmt.Errorf("card not found in market")
		}
	}

	return nil
}

func (l *gameLogic) makeMove(playerID string, isFlip bool, cardIDs ...cardID) error {
	err := l.validateMove(playerID, isFlip, cardIDs...)
	if err != nil {
		l.invalidMoves[playerID]++
		if len(l.playerIDs) == 2 && l.invalidMoves[playerID] > MAX_INVALID_MOVES {
			l.forfait(playerID)
			return nil
		}
		return err
	}

	if isFlip {
		if len(cardIDs) == 1 {
			l.handleFlip(playerID, cardIDs[0])
		}
		l.nextPlayer()

		l.isFlipNext = false

		if l.market.countEmptyPlaces() == (MARKET_HEIGHT*MARKET_WIDTH)+MARKET_WIDTH {
			l.gameOver()
		}

		return nil
	}

	l.isFlipNext = true

	// do poprawy - moze zostac tylko jedna karta
	if l.market.cardInPointCards(cardIDs[0]) {
		l.handleTakePointCard(playerID, cardIDs[0])
	} else {
		l.handleTakeVegetables(playerID, cardIDs...)
	}

	l.currentDeck = l.market.fill(l.currentDeck)

	return nil
}

func (l *gameLogic) handleFlip(playerID string, cardID cardID) {
	playerHand := l.playersHands[playerID]
	for i := 0; i < len(playerHand.PointCards); i++ {
		if playerHand.PointCards[i] == cardID {
			card := Deck[cardID]
			playerHand.Vegetables[card.getVegetable()]++
			playerHand.PointCards = append(playerHand.PointCards[:i], playerHand.PointCards[i+1:]...)
			break
		}
	}
}

func (l *gameLogic) handleTakePointCard(playerID string, cardID cardID) {
	playerHand := l.playersHands[playerID]
	for i := 0; i < MARKET_WIDTH; i++ {
		if l.market.pointCards[i] == cardID {
			playerHand.PointCards = append(playerHand.PointCards, cardID)
			l.market.pointCards[i] = ""
			l.playersHands[playerID] = playerHand
			break
		}
	}
}

func (l *gameLogic) handleTakeVegetables(playerID string, cardIDs ...cardID) {
	playerHand := l.playersHands[playerID]
	for _, cardID := range cardIDs {
		for i := 0; i < MARKET_WIDTH; i++ {
			for j := 0; j < MARKET_HEIGHT; j++ {
				if l.market.vegetableCards[i][j] == cardID {
					playerHand.Vegetables[Deck[cardID].getVegetable()]++
					l.market.vegetableCards[i][j] = ""
				}
			}
		}
	}
}

func (l *gameLogic) nextPlayer() {
	l.playerToMove++
	if l.playerToMove == len(l.playerIDs) {
		l.playerToMove = 0
	}
}

func (l *gameLogic) gameOver() {
	l.isGameOver = true
	l.finalScores = map[string]int32{}
	maxPoints := -999
	i := l.startingPlayerIdx
	for {
		points := l.countPoints(l.playerIDs[i])
		l.finalScores[l.playerIDs[i]] = int32(points)

		if points >= maxPoints {
			maxPoints = points
			l.winner = l.playerIDs[i]
		}

		i++
		if i == len(l.playerIDs) {
			i = 0
		}

		if i == l.startingPlayerIdx {
			break
		}
	}
}

func GetStateLog(gameState *proto.GameState) *State {
	pointCards := [MARKET_WIDTH]string{}
	for i, card := range gameState.Market.PointCards {
		pointCards[i] = card.CardID
	}

	vegetableCards := [MARKET_WIDTH][MARKET_HEIGHT]string{}
	for i, card := range gameState.Market.VegetableCards {
		width, height := i/MARKET_WIDTH, i%MARKET_HEIGHT
		vegetableCards[width][height] = card.CardID
	}
	s := &State{
		IsFlip:         gameState.MoveToMake == proto.MoveType_FLIP_CARD,
		PlayerTurn:     gameState.PlayerToMove,
		PointCards:     pointCards,
		VegetableCards: vegetableCards,
	}

	myPointCards := []string{}
	for _, card := range gameState.YourHand.PointCards {
		myPointCards = append(myPointCards, card.CardID)
	}
	myVegetables := map[vegetableType]int{}
	for _, card := range gameState.YourHand.Vegetables {
		myVegetables[getVegetableType(card.VegetableType)] = int(card.Count)
	}
	playersHands := map[string]*playerHand{}
	playersHands["me"] = &playerHand{
		CurrentScore: 0,
		PointCards:   myPointCards,
		Vegetables:   myVegetables,
	}

	opponentsPointCards := []string{}
	for _, card := range gameState.OpponentsHands[0].PointCards {
		opponentsPointCards = append(opponentsPointCards, card.CardID)
	}
	opponentsVegetables := map[vegetableType]int{}
	for _, card := range gameState.OpponentsHands[0].Vegetables {
		opponentsVegetables[getVegetableType(card.VegetableType)] = int(card.Count)
	}
	playersHands["opponent"] = &playerHand{
		CurrentScore: 0,
		PointCards:   opponentsPointCards,
		Vegetables:   opponentsVegetables,
	}

	return s.GetStateLog(playersHands, countPoints(playersHands))
}

func getVegetableType(vegType proto.VegetableType) vegetableType {
	switch vegType {
	case proto.VegetableType_TOMATO:
		return Tomato
	case proto.VegetableType_CARROT:
		return Carrot
	case proto.VegetableType_LETTUCE:
		return Lettuce
	case proto.VegetableType_CABBAGE:
		return Cabbage
	case proto.VegetableType_PEPPER:
		return Pepper
	case proto.VegetableType_ONION:
		return Onion
	}
	slog.Warn("Unknown vegetable type", slog.Any("vegType", vegType))
	return ""
}

func countPoints(playersHands map[string]*playerHand) func(playerID string) int {
	return func(playerID string) int {
		playerHand := playersHands[playerID]
		points := 0
		for _, cardID := range playerHand.PointCards {
			card := Deck[cardID]
			if card == nil {
				continue
			}
			points += scoreCard(playerID, playersHands, card)
		}
		return points
	}
}

func (l *gameLogic) GetStateLog() *State {
	s := &State{
		IsFlip:         l.isFlipNext,
		PlayerTurn:     l.GetPlayerToMove(),
		PointCards:     l.market.pointCards,
		VegetableCards: l.market.vegetableCards,
	}
	return s.GetStateLog(l.playersHands, l.countPoints)
}

func (s *State) GetStateLog(playersHands map[string]*playerHand, countPoints func(string) int) *State {
	handMap := make(map[string]playerHand)
	defer func() {
		handMap = nil
	}()

	for id, h := range playersHands {
		// tutaj akrobacje bo mapy się psuły i nadpisywały
		nowaMapa := map[vegetableType]int{}
		for k, v := range h.Vegetables {
			nowaMapa[k] = v
		}

		handMap[id] = playerHand{
			CurrentScore: countPoints(id),
			PointCards:   append([]string{}, h.PointCards...),
			Vegetables:   nowaMapa,
		}
	}
	s.Hands = handMap

	return s
}

func (l *gameLogic) forfait(playerID string) {
	l.isGameOver = true
	l.finalScores = map[string]int32{}
	for _, pID := range l.playerIDs {
		if pID == playerID {
			l.finalScores[pID] = -999
		} else {
			l.finalScores[pID] = 0
		}

		if playerID != pID {
			l.winner = pID
		}
	}
}
