package calculator

import (
	"encoding/json"
	"io/ioutil"
)

type Data struct {
	Brands map[string]*Brand `json:"brands"`
}

type Brand struct {
	Name   string              `json:"name"`
	Models map[string]*Model `json:"models"`
}

type Model struct {
	Options    map[string]*Option `json:"options"`
	Name       string               `json:"name"`
	DailyPrice int                  `json:"dailyPrice"`
}

type Option struct {
	Name       string `json:"name"`
	DailyPrice int    `json:"dailyPrice"`
}

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