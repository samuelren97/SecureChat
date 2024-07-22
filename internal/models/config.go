package models

type ConfigModel struct {
	Server ServerConfigModel `json:"server"`
}

type ServerConfigModel struct {
	Address string `json:"address"`
}
