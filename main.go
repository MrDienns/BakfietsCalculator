package main

import (
	"github.com/MrDienns/BakfietsCalculator/calculator"
)

func main() {
	data, err := calculator.ReadData("data/storedata.json")
	if err != nil {
		panic(err)
	}
	view := calculator.NewView(data)
	view.Open()
}
