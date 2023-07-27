package csvs

import "server_logic/utils"

type ConfigCook struct {
	CookId int `json:"CookId"`
	Star   int `json:"Star"`
}

var (
	ConfigCookMap map[int]*ConfigCook
)

func init() {
	ConfigCookMap = make(map[int]*ConfigCook)
	utils.GetCsvUtilMgr().LoadCsv("Cook", &ConfigCookMap)
	return
}

func GetCookConfig(cookId int) *ConfigCook {
	return ConfigCookMap[cookId]
}
