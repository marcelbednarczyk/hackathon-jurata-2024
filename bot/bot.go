package bot

import (
	"github.com/marcelbednarczyk/hackathon-jurata-2024/counter"
	"github.com/marcelbednarczyk/hackathon-jurata-2024/proto"
)

type Bot interface {
	MakeTakeCardsMove(state *proto.GameState, cou counter.Counter, i int) []string
	MakeFlipMove(state *proto.GameState, _ counter.Counter, i int) []string
}
