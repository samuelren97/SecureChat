package models

type ConfigModel struct {
	Server ServerConfigModel `json:"server"`
}

type ServerConfigModel struct {
	Port string `json:"port"`
}
