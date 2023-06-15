package vk

type Config struct {
	Token  string `json:"token"`
	ChatID int    `json:"chat_id"`
	Delay  int    `json:"delay"`
}

func NewConfig() Config {
	return Config{
		Token:  "*Long pool токен группы вк*",
		ChatID: 2000000001,
		Delay:  20,
	}
}
