package game

import (
	"fmt"
	"math/rand"
	"server_logic/csvs"
)

type Relics struct {
	RelicsId   int
	KeyId      int
	MainEntry  int
	Level      int
	Exp        int
	OtherEntry []int
	RoleId     int
}

type ModRelics struct {
	RelicsInfo map[int]*Relics
	MasKey     int
}

func (mr *ModRelics) AddItem(itemId int, num int64) {
	config := csvs.GetRelicsConfig(itemId)
	if config == nil {
		fmt.Println("该圣遗物不存在")
		return
	}
	if len(mr.RelicsInfo)+int(num) > csvs.RELICS_MAX_COUNT {
		fmt.Println("超过最大值")
		return
	}
	for i := int64(0); i < num; i++ {
		relics := mr.NewRelics(itemId)
		mr.RelicsInfo[relics.KeyId] = relics
		fmt.Println("获得圣遗物")
		relics.ShowInfo()
	}
	return
}

func (mr *ModRelics) NewRelics(itemId int) *Relics {
	relicsRel := new(Relics)
	relicsRel.RelicsId = itemId
	mr.MasKey++
	relicsRel.KeyId = mr.MasKey
	config := csvs.ConfigRelicsMap[itemId]
	if config == nil {
		return nil
	}
	relicsRel.MainEntry = mr.MakeMainEntry(config.MainGroup)
	for i := 0; i < config.OtherGroupNum; i++ {
		if i == config.OtherGroupNum-1 {
			randNum := rand.Intn(csvs.PERCENT_ALL)
			if randNum < csvs.ALL_ENTRY_RATE {
				relicsRel.OtherEntry = append(relicsRel.OtherEntry, mr.MakeOtherEntry(relicsRel, config.OtherGroup))
			}
		} else {
			relicsRel.OtherEntry = append(relicsRel.OtherEntry, mr.MakeOtherEntry(relicsRel, config.OtherGroup))
		}
	}
	return relicsRel
}

func (mr *ModRelics) MakeMainEntry(mainGroup int) int {
	configs, ok := csvs.ConfigRelicsEntryGroupMap[mainGroup]
	if !ok {
		return 0
	}
	allRate := 0
	for _, v := range configs {
		allRate += v.Weight
	}
	randNum := rand.Intn(allRate)
	nowNum := 0
	for _, v := range configs {
		nowNum += v.Weight
		if nowNum < randNum {
			return v.Id
		}
	}
	return 0
}

func (r *Relics) ShowInfo() {
	fmt.Println(fmt.Sprintf("key:%d,Id:%d,", r.KeyId, r.RelicsId))
	mainEntryConfig := csvs.GetRelicsLevelConfig(r.MainEntry, r.Level)
	fmt.Println(fmt.Sprintf("当前等级:%d,当前经验:%d", r.Level, r.Exp))
	if mainEntryConfig != nil {
		fmt.Println(fmt.Sprintf("主词条属性:%s,值:%d", mainEntryConfig.AttrName, mainEntryConfig.AttrValue))
	}
	for _, v := range r.OtherEntry {
		otherEntryConfig := csvs.ConfigRelicsEntryMap[v]
		if otherEntryConfig != nil {
			fmt.Println(fmt.Sprintf("副词条属性:%s,值:%d", otherEntryConfig.AttrName, otherEntryConfig.AttrValue))
		}
	}
	return
}

func (mr *ModRelics) MakeOtherEntry(relics *Relics, otherGroup int) int {
	configs, ok := csvs.ConfigRelicsEntryGroupMap[otherGroup]
	if !ok {
		return 0
	}
	configNow := csvs.GetRelicsConfig(relics.RelicsId)
	if configNow == nil {
		return 0
	}
	if len(relics.OtherEntry) >= configNow.OtherGroupNum {
		allEntry := make(map[int]int)
		for _, v := range relics.OtherEntry {
			otherConfig, _ := csvs.ConfigRelicsEntryMap[v]
			if otherConfig != nil {
				allEntry[otherConfig.AttrType] = csvs.LOGIC_TRUE
			}
		}
		allRate := 0
		for _, v := range configs {
			_, ok := allEntry[v.AttrType]
			if !ok {
				continue
			}
			allRate += v.Weight
		}
		randNum := rand.Intn(allRate)
		nowNum := 0
		for _, v := range configs {
			_, ok := allEntry[v.AttrType]
			if !ok {
				continue
			}
			nowNum += v.Weight
			if nowNum < randNum {
				return v.Id
			}
		}
	} else {
		allEntry := make(map[int]int)
		mainConfig, _ := csvs.ConfigRelicsEntryMap[relics.MainEntry]
		if mainConfig != nil {
			allEntry[mainConfig.AttrType] = csvs.LOGIC_TRUE
		}
		for _, v := range relics.OtherEntry {
			otherConfig, _ := csvs.ConfigRelicsEntryMap[v]
			if otherConfig != nil {
				allEntry[otherConfig.AttrType] = csvs.LOGIC_TRUE
			}
		}
		allRate := 0
		for _, v := range configs {
			_, ok := allEntry[v.AttrType]
			if ok {
				continue
			}
			allRate += v.Weight
		}
		randNum := rand.Intn(allRate)
		nowNum := 0
		for _, v := range configs {
			_, ok := allEntry[v.AttrType]
			if ok {
				continue
			}
			nowNum += v.Weight
			if nowNum < randNum {
				return v.Id
			}
		}
	}
	return 0

}

func (mr *ModRelics) RelicsUp(player *Player) {
	relics := mr.RelicsInfo[1]
	if relics == nil {
		fmt.Println("找不到对应圣遗物")
		return
	}
	relics.Exp += 100000
	for {
		nextLevelConfig := csvs.GetRelicsLevelConfig(relics.MainEntry, relics.Level+1)
		if nextLevelConfig == nil {
			fmt.Println("升到满级了")
			break
		}
		if relics.Exp < nextLevelConfig.NeedExp {
			break
		}
		relics.Level++
		relics.Exp -= nextLevelConfig.NeedExp
		if relics.Level%4 == 0 {
			relicsConfig := csvs.ConfigRelicsMap[relics.RelicsId]
			if relicsConfig != nil {
				relics.OtherEntry = append(relics.OtherEntry, mr.MakeOtherEntry(relics, relicsConfig.OtherGroup))
			}
		}
	}
	relics.ShowInfo()
}
