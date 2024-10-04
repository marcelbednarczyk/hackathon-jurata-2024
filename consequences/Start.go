package consequences

import "github.com/marcelbednarczyk/hackathon-jurata-2024/proto"

type Consequences struct {
	PointCards     map[string]int
	VegetableCards map[string]map[string]int
}

func Start(gameState *proto.GameState) Consequences {
	return Consequences{}
}
