package logic

type pointType string

const (
	PointsPerVegetableOne   pointType = "PointsPerVegetableOne"
	PointsPerVegetableTwo   pointType = "PointsPerVegetableTwo"
	PointsPerVegetableThree pointType = "PointsPerVegetableThree"
	SumTwo                  pointType = "SumTwo"
	SumThree                pointType = "SumThree"
	EvenOdd                 pointType = "EvenOdd"
	Fewest                  pointType = "Fewest"
	Most                    pointType = "Most"
	MostTotal               pointType = "MostTotal"
	FewestTotal             pointType = "FewestTotal"
	CompleteSet             pointType = "CompleteSet"
	AtLeastTwo              pointType = "AtLeastTwo"
	AtLeastThree            pointType = "AtLeastThree"
	MissingVegetable        pointType = "MissingVegetable"
)

var AllPointTypes = map[pointType]bool{
	PointsPerVegetableOne:   true,
	PointsPerVegetableTwo:   true,
	PointsPerVegetableThree: true,
	SumTwo:                  true,
	SumThree:                true,
	EvenOdd:                 true,
	Fewest:                  true,
	Most:                    true,
	MostTotal:               true,
	FewestTotal:             true,
	CompleteSet:             true,
	AtLeastTwo:              true,
	AtLeastThree:            true,
	MissingVegetable:        true,
}

var pointToUpper = map[pointType]string{
	PointsPerVegetableOne:   "POINTS_PER_VEGETABLE_ONE",
	PointsPerVegetableTwo:   "POINTS_PER_VEGETABLE_TWO",
	PointsPerVegetableThree: "POINTS_PER_VEGETABLE_THREE",
	SumTwo:                  "SUM_TWO",
	SumThree:                "SUM_THREE",
	EvenOdd:                 "EVEN_ODD",
	Fewest:                  "FEWEST",
	Most:                    "MOST",
	MostTotal:               "MOST_TOTAL",
	FewestTotal:             "FEWEST_TOTAL",
	CompleteSet:             "COMPLETE_SET",
	AtLeastTwo:              "AT_LEAST_TWO",
	AtLeastThree:            "AT_LEAST_THREE",
	MissingVegetable:        "MISSING_VEGETABLE",
}

func (p pointType) UpperString() string {
	return pointToUpper[p]
}
