package bot

import "github.com/marcelbednarczyk/hackathon-jurata-2024/proto"

type bot interface {
	MakeTakeCardsMove(state *proto.GameState) []string
	MakeFlipMove(state *proto.GameState) []string
}
