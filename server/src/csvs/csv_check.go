package csvs

import (
	"fmt"
	"math/rand"
)

var (
	ConfigDropGroupMap        map[int]*DropGroup
	ConfigDropItemGroupMap    map[int]*DropItemGroup
	ConfigStatueMap           map[int]map[int]*ConfigStatue
	ConfigRelicsEntryGroupMap map[int]map[int]*ConfigRelicsEntry
	ConfigRelicsLevelMap      map[int]map[int]*ConfigRelicsLevel
	ConfigRelicsSuitMap       map[int][]*ConfigRelicsSuit
	ConfigWeaponLevelMap      map[int]map[int]*ConfigWeaponLevel
	ConfigWeaponStarMap       map[int]map[int]*ConfigWeaponStar
)

type DropGroup struct {
	DropId      int
	WeightAll   int
	DropConfigs []*ConfigDrop
}

type DropItemGroup struct {
	DropId      int
	DropConfigs []*ConfigDropItem
}

func CheckLoadCsv() {
	// 二次处理
	MakeDropGroupMap()
	MakeDropItemGroupMap()
	MakeConfigStatueMap()
	MakeConfigRelicsEntryGroupMap()
	MakeConfigRelicsLevelMap()
	MakeConfigRelicsSuitMap()
	MakeConfigWeaponLevelMap()
	MakeConfigWeaponStarMap()
	fmt.Println("csv初始化完成")
	RandDropItemTest()
}

func MakeDropGroupMap() {
	ConfigDropGroupMap = make(map[int]*DropGroup)
	for _, v := range ConfigDropSlice {
		dropGroup, ok := ConfigDropGroupMap[v.DropId]
		if !ok {
			dropGroup = new(DropGroup)
			dropGroup.DropId = v.DropId
			ConfigDropGroupMap[v.DropId] = dropGroup
		}
		dropGroup.WeightAll += v.Weight
		dropGroup.DropConfigs = append(dropGroup.DropConfigs, v)
	}
	return
}

func MakeDropItemGroupMap() {
	ConfigDropItemGroupMap = make(map[int]*DropItemGroup)
	for _, v := range ConfigDropItemSlice {
		dropGroup, ok := ConfigDropItemGroupMap[v.DropId]
		if !ok {
			dropGroup = new(DropItemGroup)
			dropGroup.DropId = v.DropId
			ConfigDropItemGroupMap[v.DropId] = dropGroup
		}
		dropGroup.DropConfigs = append(dropGroup.DropConfigs, v)
	}
	return
}

func MakeConfigStatueMap() {
	ConfigStatueMap = make(map[int]map[int]*ConfigStatue)
	for _, v := range ConfigStatueSlice {
		statueMap, ok := ConfigStatueMap[v.StatueId]
		if !ok {
			statueMap = make(map[int]*ConfigStatue)
			ConfigStatueMap[v.StatueId] = statueMap
		}
		statueMap[v.Level] = v
	}
	return
}

func RandDropItemTest() {
	dropGroup := ConfigDropItemGroupMap[1]
	if dropGroup == nil {
		return
	}
	for _, v := range dropGroup.DropConfigs {
		randNum := rand.Intn(PERCENT_ALL)
		if randNum < v.Weight {
			println(v.ItemId)
		}
	}
	return
}

func GetDropItemGroup(dropId int) *DropItemGroup {
	return ConfigDropItemGroupMap[dropId]
}

func GetDropItemGroupNew(dropId int) []*ConfigDropItem {
	rel := make([]*ConfigDropItem, 0)
	config := GetDropItemGroup(dropId)
	configAll := make([]*ConfigDropItem, 0)
	for _, v := range config.DropConfigs {
		if v.DropType == DROP_ITEM_TYPE_ITEM {
			rel = append(rel, v)
		} else if v.DropType == DROP_ITEM_TYPE_GROUP {
			randNum := rand.Intn(PERCENT_ALL)
			if randNum < v.Weight {
				config := GetDropItemGroupNew(v.ItemId)
				rel = append(rel, config...)
			}

		} else if v.DropType == DROP_ITEM_TYPE_WEIGHT {
			configAll = append(configAll, v)
		}
	}
	if len(configAll) > 0 {
		allRate := 0
		for _, v := range configAll {
			allRate += v.Weight
		}
		randNum := rand.Intn(allRate)
		nowRate := 0
		for _, v := range configAll {
			nowRate += v.Weight
			if nowRate > randNum {
				newConfig := new(ConfigDropItem)
				newConfig.Weight = PERCENT_ALL
				newConfig.DropType = v.DropType
				newConfig.ItemId = v.ItemId
				newConfig.DropId = v.DropId
				newConfig.ItemNumMin = v.ItemNumMin
				newConfig.ItemNumMax = v.ItemNumMax
				newConfig.WorldAdd = v.WorldAdd
				rel = append(rel, newConfig)
				break
			}
		}
	}
	return rel
}

func MakeConfigRelicsLevelMap() {
	ConfigRelicsLevelMap = make(map[int]map[int]*ConfigRelicsLevel)
	for _, v := range ConfigRelicsLevelSlice {
		levelMap, ok := ConfigRelicsLevelMap[v.EntryId]
		if !ok {
			levelMap = make(map[int]*ConfigRelicsLevel)
			ConfigRelicsLevelMap[v.EntryId] = levelMap
		}
		levelMap[v.Level] = v
	}
	return
}

func GetStatueConfig(statueId int, level int) *ConfigStatue {
	_, ok := ConfigStatueMap[statueId]
	if !ok {
		return nil
	}
	_, ok = ConfigStatueMap[statueId][level]
	if !ok {
		return nil
	}
	return ConfigStatueMap[statueId][level]
}

func MakeConfigRelicsEntryGroupMap() {
	ConfigRelicsEntryGroupMap = make(map[int]map[int]*ConfigRelicsEntry)
	for _, v := range ConfigRelicsEntryMap {
		groupMap, ok := ConfigRelicsEntryGroupMap[v.Group]
		if !ok {
			groupMap = make(map[int]*ConfigRelicsEntry)
			ConfigRelicsEntryGroupMap[v.Group] = groupMap
		}
		groupMap[v.Id] = v
	}
	return
}

func MakeConfigRelicsSuitMap() {
	ConfigRelicsSuitMap = make(map[int][]*ConfigRelicsSuit)
	for _, v := range ConfigRelicsSuitSlice {
		ConfigRelicsSuitMap[v.Type] = append(ConfigRelicsSuitMap[v.Type], v)
	}
	return
}

func MakeConfigWeaponLevelMap() {
	ConfigWeaponLevelMap = make(map[int]map[int]*ConfigWeaponLevel)
	for _, v := range ConfigWeaponLevelSlice {
		levelMap, ok := ConfigWeaponLevelMap[v.WeaponStar]
		if !ok {
			levelMap = make(map[int]*ConfigWeaponLevel)
			ConfigWeaponLevelMap[v.WeaponStar] = levelMap
		}
		levelMap[v.Level] = v
	}
	return
}

func GetRelicsLevelConfig(mainEntry, level int) *ConfigRelicsLevel {
	_, ok := ConfigRelicsLevelMap[mainEntry]
	if !ok {
		return nil
	}
	v, ok := ConfigRelicsLevelMap[mainEntry][level]
	if !ok {
		return nil
	}
	return v
}

func GetWeaponLevelConfig(weaponStart int, level int) *ConfigWeaponLevel {
	_, ok := ConfigWeaponLevelMap[weaponStart]
	if !ok {
		return nil
	}
	v, ok := ConfigWeaponLevelMap[weaponStart][level]
	if !ok {
		return nil
	}
	return v
}

func GetWeaponStarConfig(weaponStart int, startLevel int) *ConfigWeaponStar {
	_, ok := ConfigWeaponStarMap[weaponStart]
	if !ok {
		return nil
	}
	v, ok := ConfigWeaponStarMap[weaponStart][startLevel]
	if !ok {
		return nil
	}
	return v
}

func MakeConfigWeaponStarMap() {
	ConfigWeaponStarMap = make(map[int]map[int]*ConfigWeaponStar)
	for _, v := range ConfigWeaponStarSlice {
		starMap, ok := ConfigWeaponStarMap[v.WeaponStar]
		if !ok {
			starMap = make(map[int]*ConfigWeaponStar)
			ConfigWeaponStarMap[v.WeaponStar] = starMap
		}
		starMap[v.StarLevel] = v
	}
	return
}
