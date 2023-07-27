package game

import (
	"fmt"
	"server_logic/csvs"
)

type Cook struct {
	CookId int
}

type ModCook struct {
	CookInfo map[int]*Cook
}

func (mc *ModCook) AddItem(itemId int) {
	_, ok := mc.CookInfo[itemId]
	if ok {
		fmt.Println("已习得,", csvs.GetItemConfig(itemId).ItemName)
		return
	}
	config := csvs.GetCookConfig(itemId)
	if config == nil {
		fmt.Println("非法烹饪技能,", csvs.GetItemConfig(itemId).ItemName)
		return
	}
	mc.CookInfo[itemId] = &Cook{CookId: itemId}
	fmt.Println("学会烹饪技能,", itemId)

	return
}
