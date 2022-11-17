package algorithms

import (
	"math"

	"voima.com/gophase/database"
	"voima.com/gophase/numeric"
)

type Composition struct {
	Component string
	Moles     float64
}

type Phase struct {
	Type       string
	Components []Composition
}

type PhaseDistribution struct {
	TotalPhases int
	Phases      []Phase
}

// K Calculations

type KCalculator func(componentName string, TK float64, P float64) float64

type FlashCalculationMethod func(mixture []Composition, TK float64, P float64) PhaseDistribution
type BubblePCalculationMethod func(mixture []Composition, TK float64) PhaseDistribution
type BubbleTCalculationMethod func(mixture []Composition, P float64) PhaseDistribution

func GetRachfordRiceMethod(db *database.ComponentDatabase, kCalculator KCalculator) FlashCalculationMethod {
	return func(mixture []Composition, TK float64, P float64) PhaseDistribution {
		return RachfordRice(db, mixture, TK, P, kCalculator)
	}
}

func GetSimpleBubblePCalculationMethod(db *database.ComponentDatabase, kCalculator KCalculator) BubblePCalculationMethod {
	return func(mixture []Composition, TK float64) PhaseDistribution {
		return SimpleBubbleP(db, mixture, TK, kCalculator)
	}
}

func GetWilsonKCalculator(db *database.ComponentDatabase) KCalculator {
	return func(componentName string, TK float64, P float64) float64 {
		return WilsonK(db, componentName, TK, P)
	}
}

func WilsonK(db *database.ComponentDatabase, componentName string, TK float64, P float64) float64 {
	Pc := (*db)[componentName].Data["Pc"]
	Tc := (*db)[componentName].Data["Tc"]
	wacc := (*db)[componentName].Data["wacc"]
	Pr := P / Pc
	Tr := TK / Tc
	Ki := 1 / Pr * math.Exp(5.37*(1+wacc)*(1-1/Tr))
	return Ki
}

// Flash

func calculateTotalMoles(mixture []Composition) float64 {
	sum := 0.0
	for _, composition := range mixture {
		sum += composition.Moles
	}
	return sum
}

func distributePhases(mixture []Composition, vaporFraction float64, TK float64, P float64, kCalculator KCalculator) []Phase {

	liquidPhase := Phase{
		Type:       "liquid",
		Components: make([]Composition, 0),
	}
	vaporPhase := Phase{
		Type:       "vapor",
		Components: make([]Composition, 0),
	}

	totalMoles := calculateTotalMoles(mixture)
	vaporMoles := vaporFraction * totalMoles
	liquidMoles := (1 - vaporFraction) * totalMoles

	for _, composition := range mixture {
		component := composition.Component
		Ki := kCalculator(component, TK, P)

		zi := composition.Moles / totalMoles
		yi := zi * Ki / (1 + vaporFraction*(Ki-1))
		xi := yi / Ki

		liquidPhase.Components = append(liquidPhase.Components, Composition{
			Component: component,
			Moles:     liquidMoles * xi,
		})

		vaporPhase.Components = append(vaporPhase.Components, Composition{
			Component: component,
			Moles:     vaporMoles * yi,
		})

	}

	return []Phase{liquidPhase, vaporPhase}
}

// Dispatchers

func FlashTP(mixture []Composition, TK float64, P float64, method FlashCalculationMethod) PhaseDistribution {
	return method(mixture, TK, P)
}

func BubbleP(mixture []Composition, TK float64, method BubblePCalculationMethod) PhaseDistribution {
	return method(mixture, TK)
}

func BubbleT(mixture []Composition, P float64, method BubbleTCalculationMethod) PhaseDistribution {
	return method(mixture, P)
}

// Concrete implementations

func SimpleBubbleP(db *database.ComponentDatabase, mixture []Composition, TK float64, kCalculator KCalculator) PhaseDistribution {
	vaporPhase := Phase{
		Type:       "vapor",
		Components: make([]Composition, 0),
	}
	return PhaseDistribution{
		TotalPhases: 2,
		Phases:      []Phase{vaporPhase},
	}

}
func SimpleBubbleT(db *database.ComponentDatabase, mixture []Composition, P float64, kCalculator KCalculator) PhaseDistribution {
	vaporPhase := Phase{
		Type:       "vapor",
		Components: make([]Composition, 0),
	}
	return PhaseDistribution{
		TotalPhases: 2,
		Phases:      []Phase{vaporPhase},
	}
}

func RachfordRice(db *database.ComponentDatabase, mixture []Composition, TK float64, P float64, kCalculator KCalculator) PhaseDistribution {

	rachfordRiceDerivative := func(fun numeric.Function1D, V float64) float64 {
		sum := 0.0
		for _, composition := range mixture {
			zi := composition.Moles
			Ki := WilsonK(db, composition.Component, TK, P)
			sum += -1 * math.Pow(Ki-1, 2) * zi / math.Pow(1+V*(Ki-1), 2)
		}
		return sum
	}

	rachfordRice := func(V float64) float64 {
		sum := 0.0
		for _, composition := range mixture {
			zi := composition.Moles
			Ki := kCalculator(composition.Component, TK, P)
			sum += (Ki - 1) * zi / (1 + V*(Ki-1))
		}
		return sum
	}

	res := numeric.NewtonRaphson1D(rachfordRice, 0.5, 1e-12, 50, rachfordRiceDerivative)
	V := res.X

	phases := distributePhases(mixture, V, TK, P, kCalculator)

	return PhaseDistribution{
		TotalPhases: 2,
		Phases:      phases,
	}

}
