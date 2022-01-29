package common

type Config struct {
	Meets      []Meet `json:"Meets"`
	SumId      int    `json:"SumId"`
	TimeMargin int    `json:"TimeMargin"`
	IsAsk      bool   `json:"IsAsk"`
}

func NewConfig() Config {
	return Config{
		Meets:      make([]Meet, 0),
		SumId:      0,
		TimeMargin: 20,
		IsAsk:      false,
	}
}
