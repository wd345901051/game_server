package game

import (
	"fmt"
	"server_logic/csvs"
)

type Icon struct {
	IconId int
}

type ModIcon struct {
	IconInfo map[int]*Icon
}

func (mi *ModIcon) IsHasIcon(iconId int) bool {
	_, ok := mi.IconInfo[iconId]
	return ok
}

func (mi *ModIcon) AddItem(itemId int, player *Player) {
	_, ok := mi.IconInfo[itemId]
	if ok {
		fmt.Println("已存在该头像,", itemId)
		return
	}
	config := csvs.GetIconConfig(itemId)
	if config == nil {
		fmt.Println("非法头像,", itemId)
		return
	}
	mi.IconInfo[itemId] = &Icon{IconId: itemId}
	player.ModBag.AddItemToBag(itemId, 1)
	fmt.Println("获得头像,", itemId)

	return
}

func (mi *ModIcon) CheckGetIcon(roleId int, player *Player) {
	config := csvs.GetIconConfigByRoleId(roleId)
	if config == nil {
		return
	}
	mi.AddItem(config.IconId, player)
}
