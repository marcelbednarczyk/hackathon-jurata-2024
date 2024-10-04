package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/marcelbednarczyk/hackathon-jurata-2024/bot"
	"github.com/marcelbednarczyk/hackathon-jurata-2024/counter"
	"github.com/marcelbednarczyk/hackathon-jurata-2024/logic"
	"github.com/marcelbednarczyk/hackathon-jurata-2024/proto"
	"github.com/phsym/console-slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const IP_ADDR = "localhost:50051"

const (
	TIMEOUT_NEW_ROOM  = time.Second * 5
	TIMEOUT_JOIN_ROOM = time.Second * 100
	TIMEOUT_MAKE_MOVE = time.Second * 5
)

var (
	addr       = flag.String("addr", IP_ADDR, "server address")
	playerName = flag.String("name", "Dummy name", "player name")
	new        = flag.Bool("new", false, "creates a new room")
	roomID     = flag.String("room", "", "join a room with this id")
	gameCount  = flag.Int("count", 1, "number of games to play")
	debug      = flag.Bool("debug", false, "enable debug logs")
	botType    = flag.String("bot", "random", "bot type")
)

func init() {
	flag.Parse()
	setupLogger(*debug)
}

func setupLogger(debug bool) {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	logger := slog.New(
		console.NewHandler(os.Stderr, &console.HandlerOptions{Level: level}),
	)
	slog.SetDefault(logger)
}

func main() {
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Failed to connect to server", slog.Any("error", err))
		os.Exit(1)
	}
	defer conn.Close()

	client := proto.NewGameClient(conn)
	conn.GetState()

	if *new {
		ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_NEW_ROOM)
		defer cancel()

		newGameConfig, err := client.NewRoom(ctx, &proto.Config{RoomID: *roomID, NumberOfGames: int32(*gameCount)})
		if err != nil {
			slog.Error("Failed to create a new room", slog.Any("error", err))
			os.Exit(1)
		}

		slog.Info("Created a new room", slog.String("roomID", *roomID))
		*roomID = newGameConfig.RoomID
	}

	if *roomID == "" {
		slog.Error("Room ID is required")
		flag.Usage()
		os.Exit(1)
	}

	roomState := mustJoinRoom(client, *roomID, *playerName)
	slog.Info("Joined room", slog.Any("roomState", roomState))
	for roomState.StartNextGame {
		var b bot.Bot
		switch *botType {
		case "greedy":
			b = bot.NewGreedyBot()
		case "solucjaTymczasowa":
			b = bot.NewSolucjaTymczasowaBot()
		default:
			b = bot.NewRandomBot()
		}
		gameState, err := getCurrentGameState(client, roomState.RoomID, roomState.PlayerID)
		if err != nil {
			slog.Error("Failed to get current game state", slog.Any("error", err))
			os.Exit(1)
		}

		i := 0
		cardsCounter := counter.Start(gameState)
		slog.Info("Counter started", slog.String("counter", cardsCounter.String()))
		for !gameState.IsGameOver {
			for {
				logicState := logic.GetStateLog(gameState)
				slog.Info("State log",
					slog.Int("myCrrentScore", logicState.Hands["me"].CurrentScore),
					slog.Int("opponentCurrentScore", logicState.Hands["opponent"].CurrentScore),
				)
				move := &proto.MoveRequest{
					RoomID:   roomState.RoomID,
					PlayerID: roomState.PlayerID,
					MoveType: gameState.MoveToMake,
				}

				if gameState.MoveToMake == proto.MoveType_TAKE_CARDS {
					move.Cards = b.MakeTakeCardsMove(gameState, i)
				} else {
					move.Cards = b.MakeFlipMove(gameState)
				}

				newState, err := makeMove(client, move)
				if err != nil && status.Code(err) == codes.InvalidArgument {
					slog.Warn("Invalid move", slog.Any("move", move), slog.Any("error", err))
					continue
				}

				if err != nil {
					slog.Error("Failed to make move", slog.Any("error", err))
					os.Exit(1)
				}
				cardsCounter.Update(newState)
				slog.Info("Counter updated", slog.String("counter", cardsCounter.String()))

				gameState = newState
				break
			}
			i++
		}
		roomState, err = getRoomState(client, roomState.RoomID, roomState.PlayerID)
		if err != nil {
			slog.Error("Failed to get room state", slog.Any("error", err))
			os.Exit(1)
		}
	}

	err = saveGameLog(roomState.LastGameLog)
	if err != nil {
		slog.Error("Failed to save game log", slog.Any("error", err))
	}

	roomState.LastGameLog = nil

	slog.Info("Game finished", slog.Any("roomState", roomState))
}

func saveGameLog(gameLog []byte) error {
	err := os.MkdirAll("./game-logs", os.ModePerm)
	if err != nil {
		return err
	}

	now := strings.ReplaceAll(time.Now().Format(time.DateTime), " ", "_")
	now = strings.ReplaceAll(now, ":", "-")

	return os.WriteFile("./game-logs/"+now+".json", gameLog, os.ModePerm)
}

func mustJoinRoom(client proto.GameClient, roomID, playerName string) *proto.RoomState {

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_JOIN_ROOM)
	defer cancel()

	roomState, err := client.JoinRoom(ctx, &proto.JoinRoomRequest{RoomID: roomID, PlayerName: playerName})
	if err != nil {
		slog.Error("Failed to join room", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Debug("Joined room", slog.Any("roomState", roomState))
	return roomState
}

func makeMove(client proto.GameClient, move *proto.MoveRequest) (*proto.GameState, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_MAKE_MOVE)
	defer cancel()

	gameState, err := client.MakeMove(ctx, move)
	if err != nil {
		slog.Error("Failed to make move", slog.Any("error", err))
		return nil, err
	}

	slog.Debug("Made a move", slog.Any("move", move))
	return gameState, nil
}

func getCurrentGameState(client proto.GameClient, roomID, playerID string) (*proto.GameState, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_JOIN_ROOM)
	defer cancel()

	gameState, err := client.GetCurrentGameState(ctx, &proto.GetRoomRequest{RoomID: roomID, PlayerID: playerID})
	if err != nil {
		slog.Error("Failed to get current game state", slog.Any("error", err))
		return nil, err
	}

	slog.Debug("Got current game state", slog.Any("gameState", gameState))
	return gameState, nil
}

func getRoomState(client proto.GameClient, roomID, playerID string) (*proto.RoomState, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_JOIN_ROOM)
	defer cancel()

	roomState, err := client.GetRoomState(ctx, &proto.GetRoomRequest{RoomID: roomID, PlayerID: playerID})
	if err != nil {
		slog.Error("Failed to get room state", slog.Any("error", err))
		return nil, err
	}

	slog.Debug("Got room state", slog.Any("roomState", roomState))
	return roomState, nil
}
