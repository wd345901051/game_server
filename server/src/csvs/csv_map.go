package csvs

import "server_logic/utils"

type ConfigMap struct {
	MapId   int    `json:"MapId"`
	MapName string `json:"MapName"`
	MapType int    `json:"MapType"`
}

type ConfigMapEvent struct {
	EventId        int    `json:"EventId"`
	EventType      int    `json:"EventType"`
	RefreshType    int    `json:"RefreshType"`
	Name           string `json:"Name"`
	EventDrop      int    `json:"EventDrop"`
	EventDropTimes int    `json:"EventDropTimes"`
	MapId          int    `json:"MapId"`
	CostItem       int    `json:"CostItem"`
	CostNum        int64  `json:"CostNum"`
}

var (
	ConfigMapMap      map[int]*ConfigMap
	ConfigMapEventMap map[int]*ConfigMapEvent
)

func init() {
	ConfigMapMap = make(map[int]*ConfigMap)
	utils.GetCsvUtilMgr().LoadCsv("Map", &ConfigMapMap)

	ConfigMapEventMap = make(map[int]*ConfigMapEvent)
	utils.GetCsvUtilMgr().LoadCsv("MapEvent", &ConfigMapEventMap)
	return
}
