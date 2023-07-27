package game

import (
	"fmt"
	"math/rand"
	"server_logic/csvs"
	"time"
)

type Map struct {
	MapId     int
	EventInfo map[int]*Event
}

type Event struct {
	EventId       int
	State         int
	NextResetTime int64
}

type ModMap struct {
	MapInfo map[int]*Map
	Statue  map[int]*StatueInfo
}

type StatueInfo struct {
	StatueId int
	Level    int
	ItemInfo map[int]*ItemInfo
}

func (mm *ModMap) InitData() {
	mm.MapInfo = make(map[int]*Map)
	mm.Statue = make(map[int]*StatueInfo)
	for _, v := range csvs.ConfigMapMap {
		_, ok := mm.MapInfo[v.MapId]
		if !ok {
			mm.MapInfo[v.MapId] = mm.NewMapInfo(v.MapId)
		}
	}

	for _, v := range csvs.ConfigMapEventMap {
		_, ok := mm.MapInfo[v.MapId]
		if !ok {
			continue
		}
		_, ok = mm.MapInfo[v.MapId].EventInfo[v.EventId]
		if !ok {
			mm.MapInfo[v.MapId].EventInfo[v.EventId] = new(Event)
			mm.MapInfo[v.MapId].EventInfo[v.EventId].EventId = v.EventId
			mm.MapInfo[v.MapId].EventInfo[v.EventId].State = csvs.EVENT_START
		}
	}
	for _, v := range mm.MapInfo {
		fmt.Println("地图:", csvs.ConfigMapMap[v.MapId].MapName)
		for _, v := range mm.MapInfo[v.MapId].EventInfo {
			fmt.Println("事件:", csvs.ConfigMapEventMap[v.EventId].Name)
		}
	}
}

func (mm *ModMap) NewMapInfo(mapId int) *Map {
	mapInfo := new(Map)
	mapInfo.MapId = mapId
	mapInfo.EventInfo = make(map[int]*Event)
	return mapInfo
}

func (mm *ModMap) SetEventState(mapId, eventId, state int, player *Player) {
	_, ok := mm.MapInfo[mapId]
	if !ok {
		fmt.Println("地图不存在")
		return
	}
	_, ok = mm.MapInfo[mapId].EventInfo[eventId]
	if !ok {
		fmt.Println("地图事件不存在", csvs.ConfigMapMap[mapId].MapName, csvs.ConfigMapEventMap[eventId])
		return
	}
	if mm.MapInfo[mapId].EventInfo[eventId].State >= state {
		fmt.Println("状态异常")
		return
	}

	config := csvs.ConfigMapMap[mapId]
	if config == nil {
		return
	}
	if !player.ModBag.HasEnoughItem(csvs.ConfigMapEventMap[mapId].CostItem, csvs.ConfigMapEventMap[mapId].CostNum) {
		fmt.Println(csvs.GetItemConfig(csvs.ConfigMapEventMap[mapId].CostItem).ItemName, "不足")
		return
	}
	if config.MapType == csvs.REFRESH_PLAYER && csvs.ConfigMapEventMap[mm.MapInfo[mapId].EventInfo[eventId].EventId].EventType == csvs.EVENT_TYPE_REWARD {
		for _, v := range mm.MapInfo[mapId].EventInfo {
			eventConfig := csvs.ConfigMapEventMap[v.EventId]
			if eventConfig == nil {
				return
			}
			if eventConfig.EventType == csvs.EVENT_TYPE_NOMAL {
				continue
			}
			if v.EventId == eventId {
				continue
			}
			if v.State != csvs.EVENT_END {
				fmt.Println("任务未完成!")
				return
			}
		}
	}
	mm.MapInfo[mapId].EventInfo[eventId].State = state
	if state == csvs.EVENT_FINISH {
		fmt.Println("状态完成")
	}
	eventConfig := csvs.ConfigMapEventMap[mm.MapInfo[mapId].EventInfo[eventId].EventId]
	if eventConfig == nil {
		return
	}
	if state == csvs.EVENT_END {
		for i := 0; i < eventConfig.EventDropTimes; i++ {
			config := csvs.GetDropItemGroupNew(eventConfig.EventDrop)
			for _, v := range config {
				randNum := rand.Intn(csvs.PERCENT_ALL)
				if randNum < v.Weight {
					randAll := v.ItemNumMax - v.ItemNumMin + 1
					itemNum := rand.Intn(randAll) + v.ItemNumMin
					worldLevel := player.ModPlayer.WorldLevelNow
					if worldLevel > 0 {
						itemNum = itemNum * (csvs.PERCENT_ALL + worldLevel*v.WorldAdd) / csvs.PERCENT_ALL
					}
					player.ModBag.AddItem(v.ItemId, int64(itemNum), player)
				}
			}
		}
		fmt.Println("事件领取")
	}
	if state > 0 {
		switch eventConfig.RefreshType {
		case csvs.MAP_REFRESH_SELF:
			mm.MapInfo[mapId].EventInfo[eventId].NextResetTime = time.Now().Unix() + csvs.MAP_REFRESH_SELF_TIME
		}
	}
}

func (mm *ModMap) RefreshDay() {
	for _, v := range mm.MapInfo {
		for _, v := range mm.MapInfo[v.MapId].EventInfo {
			config := csvs.ConfigMapEventMap[v.EventId]
			if config == nil {
				continue
			}
			if config.RefreshType != csvs.MAP_REFRESH_DAY {
				continue
			}
			v.State = csvs.EVENT_START
		}
	}
}

func (mm *ModMap) RefreshWeek() {
	for _, v := range mm.MapInfo {
		for _, v := range mm.MapInfo[v.MapId].EventInfo {
			config := csvs.ConfigMapEventMap[v.EventId]
			if config == nil {
				continue
			}
			if config.RefreshType != csvs.MAP_REFRESH_WEEK {
				continue
			}
			v.State = csvs.EVENT_START
		}
	}
}

func (mm *ModMap) RefreshSelf() {
	for _, v := range mm.MapInfo {
		for _, v := range mm.MapInfo[v.MapId].EventInfo {
			config := csvs.ConfigMapEventMap[v.EventId]
			if config == nil {
				continue
			}
			if config.RefreshType != csvs.MAP_REFRESH_SELF {
				continue
			}
			if time.Now().Unix() <= v.NextResetTime {
				v.State = csvs.EVENT_START
			}
		}
	}
}

func (mm *ModMap) CheckRefresh(event *Event) {
	if event.NextResetTime > time.Now().Unix() {
		return
	}
	eventConfig := csvs.ConfigMapEventMap[event.EventId]
	if eventConfig == nil {
		return
	}
	switch eventConfig.RefreshType {
	case csvs.MAP_REFRESH_DAY:
		count := time.Now().Unix() / csvs.MAP_REFRESH_DAY_TIME
		count++
		event.NextResetTime = count * csvs.MAP_REFRESH_DAY_TIME
	case csvs.MAP_REFRESH_WEEK:
		event.NextResetTime = time.Now().Unix() + csvs.MAP_REFRESH_WEEK_TIME
		count := time.Now().Unix() / csvs.MAP_REFRESH_WEEK_TIME
		count++
		event.NextResetTime = count * csvs.MAP_REFRESH_WEEK_TIME
	case csvs.MAP_REFRESH_SELF:
		event.NextResetTime = time.Now().Unix() + csvs.MAP_REFRESH_SELF_TIME
		count := time.Now().Unix() / csvs.MAP_REFRESH_SELF_TIME
		count++
		event.NextResetTime = count * csvs.MAP_REFRESH_SELF_TIME
	case csvs.MAP_REFRESH_CANT:
		return
	}
	event.State = csvs.EVENT_START
}

func (mm *ModMap) RefreshByPlayer(mapId int) {
	config := csvs.ConfigMapMap[mapId]
	if config == nil {
		return
	}
	if config.MapType != csvs.REFRESH_PLAYER {
		return
	}
	for _, v := range mm.MapInfo[config.MapId].EventInfo {
		v.State = csvs.EVENT_START
	}
}

func (mm *ModMap) NewStatue(statueId int) *StatueInfo {
	data := new(StatueInfo)
	data.Level = 0
	data.StatueId = statueId
	data.ItemInfo = make(map[int]*ItemInfo)
	return data
}

func (mm *ModMap) UpStatue(statueId int, player *Player) {
	_, ok := mm.Statue[statueId]
	if !ok {
		mm.Statue[statueId] = mm.NewStatue(statueId)
	}
	info, ok := mm.Statue[statueId]
	if !ok {
		return
	}
	nextLevel := info.Level + 1
	nextConfig := csvs.GetStatueConfig(statueId, nextLevel)
	if nextConfig == nil {
		return
	}
	_, ok = info.ItemInfo[nextConfig.CostItem]
	var nowNum int64
	if ok {
		nowNum = info.ItemInfo[nextConfig.CostItem].ItemNum
	}
	needNum := nextConfig.CostNum - nowNum

	if !player.ModBag.HasEnoughItem(nextConfig.CostItem, needNum) {
		num := player.ModBag.GetItemNum(nextConfig.CostItem)
		if num <= 0 {
			fmt.Println("无升级所需物品")
			return
		}
		_, ok = info.ItemInfo[nextConfig.CostItem]
		if !ok {
			info.ItemInfo[nextConfig.CostItem] = new(ItemInfo)
			info.ItemInfo[nextConfig.CostItem].ItemId = nextConfig.CostItem
			info.ItemInfo[nextConfig.CostItem].ItemNum = 0
		}
		item, ok := info.ItemInfo[nextConfig.CostItem]
		if !ok {
			return
		}
		item.ItemNum += num
		player.ModBag.RemoveItemToBag(nextConfig.CostItem, num, player)
		fmt.Println(info.StatueId, "神像存储成功，还需要:", nextConfig.CostNum-item.ItemNum)
	} else {
		player.ModBag.RemoveItemToBag(nextConfig.CostItem, needNum, player)
		info.Level++
		info.ItemInfo = make(map[int]*ItemInfo)
		fmt.Println(info.StatueId, "神像升级成功，当前等级:", info.Level)
	}
}
