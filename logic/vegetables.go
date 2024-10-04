package logic

type vegetableType string

const (
	Tomato  vegetableType = "Tomato"
	Carrot  vegetableType = "Carrot"
	Lettuce vegetableType = "Lettuce"
	Cabbage vegetableType = "Cabbage"
	Pepper  vegetableType = "Pepper"
	Onion   vegetableType = "Onion"
)

var AllVegetables = []vegetableType{Tomato, Carrot, Lettuce, Cabbage, Pepper, Onion}

var AllVegetablesMap = map[vegetableType]bool{
	Tomato:  true,
	Carrot:  true,
	Lettuce: true,
	Cabbage: true,
	Pepper:  true,
	Onion:   true,
}

var vegToUpper = map[vegetableType]string{
	Tomato:  "TOMATO",
	Carrot:  "CARROT",
	Lettuce: "LETTUCE",
	Cabbage: "CABBAGE",
	Pepper:  "PEPPER",
	Onion:   "ONION",
}

func (v vegetableType) UpperString() string {
	return vegToUpper[v]
}
