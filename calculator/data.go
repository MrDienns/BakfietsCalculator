package calculator

import (
	"encoding/json"
	"io/ioutil"
)

// Data is a general root level struct used to represent the top level JSON structure. In this struct are sub-elements
// that build up the JSON further.
type Data struct {
	Brands map[string]*Brand `json:"brands"`
}

// Brand struct represents a Bike brand. It contains the name, and a map which then contains the models.
type Brand struct {
	Name   string            `json:"name"`
	Models map[string]*Model `json:"models"`
}

// Model struct represents a bike model. It contains a name, a daily price and a map for the options that this bike
// has to offer.
type Model struct {
	Options    map[string]*Option `json:"options"`
	Name       string             `json:"name"`
	DailyPrice int                `json:"dailyPrice"`
}

// Option struct represents the option for a bike model. It has the name and a daily price.
type Option struct {
	Name       string `json:"name"`
	DailyPrice int    `json:"dailyPrice"`
}

// ReadData function accepts a path as argument and reads the data from the disk. If the function failed to load the
// data, an error is returned. If the function succeeded, a *Data instance is returned.
func ReadData(path string) (*Data, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	data := &Data{}
	err = json.Unmarshal(contents, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
