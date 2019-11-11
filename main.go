package main

import (
	"fmt"
	"github.com/MrDienns/BakfietsCalculator/calculator"
)

func main() {
	data, err := calculator.ReadData("data/storedata.json")
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}
