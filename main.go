package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"voima.com/gophase/algorithms"
	"voima.com/gophase/database"
)

type Calculation interface {
	calculate()
}

type Evaporator struct {
	Name string
}

func (e Evaporator) calculate() {
	fmt.Println("Calculating evaporator " + e.Name)
}

func runCalculations(entities []Calculation) {
	for _, calc := range entities {
		calc.calculate()
	}
}

func parseLine(dbText string, db *database.ComponentDatabase) error {
	data := strings.Split(dbText, "|")

	componentData := database.ComponentData{
		Data: make(map[string]float64),
	}

	for _, datum := range data {
		chunks := strings.Split(datum, "=")
		name := strings.TrimSpace(chunks[0])

		if name == "name" {
			value := strings.TrimSpace(chunks[1])
			componentData.Name = value
			continue
		}

		if name == "api" {
			value := strings.TrimSpace(chunks[1])
			componentData.Api = value
			continue
		}

		value, err := strconv.ParseFloat(strings.TrimSpace(chunks[1]), 64)
		if err != nil {
			return err
		}
		componentData.Data[name] = value
	}
	(*db)[componentData.Name] = componentData

	return nil
}

func parseComponentDB(dbText string) (database.ComponentDatabase, error) {

	db := make(database.ComponentDatabase)

	lines := strings.Split(dbText, "\n")
	for _, rawline := range lines {
		line := strings.TrimSpace(rawline)
		if len(line) == 0 {
			continue
		}
		if !strings.HasPrefix(line, "//") {
			err := parseLine(line, &db)
			if err != nil {
				return nil, err
			}
		}
	}

	return db, nil
}

// Main

func quad(x float64) float64 {
	return math.Pow(x-5, 2)
}

func main() {

	dbFilePath := "./components.txt"
	data, err := os.ReadFile(dbFilePath)

	if err != nil {
		panic(err)
	}

	db, err := parseComponentDB(string(data))

	fmt.Printf("Database loaded. %d components registered.\n", len(db))

	if err != nil {
		panic(err)
	}

	mix := []algorithms.Composition{
		{
			Component: "CH4",
			Moles:     1,
		},
		{
			Component: "C2H6",
			Moles:     2,
		},
		{
			Component: "C3H8",
			Moles:     7,
		},
	}

	wilson := algorithms.GetWilsonKCalculator(&db)
	rachfordRice := algorithms.GetRachfordRiceMethod(&db, wilson)
	simpleBubP := algorithms.GetSimpleBubblePCalculationMethod(&db, wilson)
	phases := algorithms.FlashTP(mix, 253.15, 14e5, rachfordRice)
	bubblePhases := algorithms.BubbleP(mix, 14e5, simpleBubP)

	fmt.Println(phases)
	fmt.Println(bubblePhases)

}
