package bot

import "github.com/marcelbednarczyk/hackathon-jurata-2024/proto"

type Bot interface {
	MakeTakeCardsMove(state *proto.GameState, i int) []string
	MakeFlipMove(state *proto.GameState) []string
}
