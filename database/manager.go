package database

type ComponentDatabase = map[string]ComponentData

type ComponentData struct {
	Name string
	Api  string
	Data map[string]float64
}

// Thermodynamics

func (data ComponentData) Cp(TK float64) float64 {
	switch data.Api {
	case "steam_table":
		return steamTableCpCalculator(data)
	default:
		return defaultCpCalculator(data)
	}
}

func defaultCpCalculator(data ComponentData) float64 {
	return data.Data["CpA"]
}

func steamTableCpCalculator(data ComponentData) float64 {
	return data.Data["CpA"]
}
