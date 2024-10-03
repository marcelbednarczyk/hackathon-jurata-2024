package bot

import "gitlab.com/hackathon-rainbow-2024/go-client/proto"

type bot interface {
	MakeTakeCardsMove(state *proto.GameState) []string
	MakeFlipMove(state *proto.GameState) []string
}
