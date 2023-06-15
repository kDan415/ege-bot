package ege

type Config struct {
	Participant string `json:"participant"`
}

func NewConfig() Config {
	return Config{
		Participant: "*Сессия с сайта ЕГЭ*",
	}
}
